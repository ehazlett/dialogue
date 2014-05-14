package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/ehazlett/dialogue"
	"bitbucket.org/ehazlett/dialogue/auth"
	"bitbucket.org/ehazlett/dialogue/db"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

type (
	dialogueApi struct {
		m       *martini.ClassicMartini
		rdb     *db.Rethinkdb
		auth    auth.Authenticator
		address string
	}
	AuthToken struct {
		Token string `json:"token"`
	}
	ApiError struct {
		Error string `json:"error"`
	}
	ApiResponse struct {
		Response string `json:"response"`
	}
)

func NewApi(address string, rdb *db.Rethinkdb, auth auth.Authenticator, sessionKey string) (*dialogueApi, error) {
	m := martini.Classic()
	// sessions
	store := sessions.NewCookieStore([]byte(sessionKey))
	m.Use(sessions.Sessions("dialogue", store))

	a := &dialogueApi{
		m:       m,
		rdb:     rdb,
		auth:    auth,
		address: address,
	}
	// middleware
	m.Use(render.Renderer())
	// routes
	// content
	m.Get("/topics", a.apiAuthorize, a.GetTopics)
	m.Post("/topics", a.apiAuthorize, a.PostTopics)
	m.Post("/topics/:topicId", a.apiAuthorize, a.PostTopicsPosts)
	m.Get("/topics/:topicId", a.apiAuthorize, a.GetTopic)
	m.Delete("/topics/:topicId", a.apiAuthorize, a.DeleteTopic)
	m.Delete("/posts/:postId", a.apiAuthorize, a.DeletePost)

	// authentication
	m.Post("/auth", a.Authenticate)
	m.Post("/users", a.apiAuthorize, a.PostUsers)
	m.Put("/users/:username", a.apiAuthorize, a.PutUser)
	// setup
	m.Get("/setup", a.Setup)

	return a, nil
}

func (api *dialogueApi) Run() {
	log.Info("Listening on " + api.address)
	log.Fatal(http.ListenAndServe(api.address, api.m))
}

func (api *dialogueApi) unmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func (api *dialogueApi) Setup(r *http.Request, rndr render.Render) {
	user, err := api.rdb.GetUser("admin")
	if err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error checking for admin user: %s", err),
		}
		rndr.JSON(500, e)
		return
	}
	// create admin user if missing
	if user == nil {
		log.Info("Creating admin user: username: admin password: dialogue")
		pw, err := api.auth.HashPassword("dialogue")
		if err != nil {
			e := ApiError{
				Error: fmt.Sprintf("Error generating password for admin user: %s", err),
			}
			rndr.JSON(500, e)
			return
		}
		user = &dialogue.User{
			Username: "admin",
			Password: pw,
		}
		if err := api.rdb.SaveUser(user); err != nil {
			e := ApiError{
				Error: fmt.Sprintf("Error creating user: %s", err),
			}
			rndr.JSON(500, e)
			return
		}
		r := ApiResponse{
			Response: "admin user created: username: admin password: dialogue",
		}
		rndr.JSON(200, r)
		return
	}
	e := ApiError{
		Error: "admin user already present",
	}
	rndr.JSON(400, e)
	return
}

func (api *dialogueApi) apiAuthorize(r *http.Request, session sessions.Session, rndr render.Render) {
	// check authorization headers
	username := r.Header.Get("X-Auth-User")
	token := r.Header.Get("X-Auth-Token")
	if username == "" || token == "" {
		e := ApiError{
			Error: "username and token must be present",
		}
		rndr.JSON(401, e)
		return
	}
	// verify token
	auth, err := api.rdb.GetAuthorization(username)
	if err != nil {
		e := ApiError{
			Error: fmt.Sprintf("error verifying token: %s", err),
		}
		rndr.JSON(401, e)
		return
	}
	if auth == nil || auth.Token != token {
		e := ApiError{
			Error: "invalid username/token",
		}
		rndr.JSON(401, e)
		return
	}
	// all is well, set session
	session.Set("username", username)
}

// route handlers
func (api *dialogueApi) GetTopic(params martini.Params, r render.Render) {
	topicId := params["topicId"]
	res, err := api.rdb.GetPosts(topicId)
	if err != nil {
		e := ApiError{
			Error: "Error getting posts",
		}
		r.JSON(500, e)
		return
	}
	r.JSON(200, res)
}

func (api *dialogueApi) GetTopics(params martini.Params, r render.Render) {
	res, err := api.rdb.GetTopics()
	if err != nil {
		e := ApiError{
			Error: "Error getting topics",
		}
		r.JSON(500, e)
		return
	}
	r.JSON(200, res)
}

func (api *dialogueApi) PostTopics(w http.ResponseWriter, r *http.Request, rndr render.Render) {
	title := r.FormValue("title")
	// check for title
	if title == "" {
		e := ApiError{
			Error: "title must be specified",
		}
		rndr.JSON(500, e)
		return
	}
	// new topic
	topic := &dialogue.Topic{
		Title:  title,
		Closed: false,
	}
	if err := api.rdb.SaveTopic(topic); err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error saving topic: %s", err),
		}
		rndr.JSON(500, e)
		return
	}
	w.WriteHeader(204)
}

func (api *dialogueApi) PostTopicsPosts(w http.ResponseWriter, r *http.Request, session sessions.Session, params martini.Params, rndr render.Render) {
	content := r.FormValue("content")
	topicId := params["topicId"]
	author := session.Get("username")
	// check for content
	if content == "" {
		e := ApiError{
			Error: "content must be specified",
		}
		rndr.JSON(500, e)
		return
	}
	// new post
	post := &dialogue.Post{
		Content: content,
		TopicId: topicId,
		Author:  author.(string),
	}
	if err := api.rdb.SavePost(post); err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error saving post: %s", err),
		}
		rndr.JSON(500, e)
		return
	}
	w.WriteHeader(204)
}

func (api *dialogueApi) DeleteTopic(w http.ResponseWriter, params martini.Params, rndr render.Render) {
	id := params["topicId"]
	if err := api.rdb.DeleteTopic(id); err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error deleting topic: %s", err),
		}
		rndr.JSON(500, e)
		return
	}
	w.WriteHeader(204)
}

func (api *dialogueApi) DeletePost(w http.ResponseWriter, params martini.Params, rndr render.Render) {
	id := params["postId"]
	if err := api.rdb.DeletePost(id); err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error deleting post: %s", err),
		}
		rndr.JSON(500, e)
		return
	}
	w.WriteHeader(204)
}

func (api *dialogueApi) Authenticate(r *http.Request, rndr render.Render, params martini.Params) {
	username := r.FormValue("username")
	pass := r.FormValue("password")
	user, err := api.rdb.GetUser(username)
	if err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error authenticating: %s", err),
		}
		rndr.JSON(500, e)
		return
	}
	if api.auth.Authenticate(user.Password, pass) {
		t := api.auth.GenerateToken()
		token := AuthToken{
			Token: t,
		}
		// update current auth token
		a := &dialogue.Authorization{
			Username: username,
			Token:    t,
		}
		if err := api.rdb.SaveAuthorization(a); err != nil {
			e := ApiError{
				Error: fmt.Sprintf("Error saving auth token: %s", err),
			}
			rndr.JSON(500, e)
			return
		}
		rndr.JSON(200, token)
		return
	}
	e := ApiError{
		Error: "Invalid username/password",
	}
	rndr.JSON(401, e)
	return
}

func (api *dialogueApi) PostUsers(w http.ResponseWriter, r *http.Request, rndr render.Render) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	// check for username and password
	if username == "" || password == "" {
		e := ApiError{
			Error: "username and password must be specified",
		}
		rndr.JSON(500, e)
		return
	}
	// hash password
	pw, err := api.auth.HashPassword(password)
	if err != nil {
		e := ApiError{
			Error: "error hashing password",
		}
		rndr.JSON(500, e)
		return

	}
	// new user
	user := &dialogue.User{
		Username: username,
		Password: pw,
	}
	if err := api.rdb.SaveUser(user); err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error creating user: %s", err),
		}
		rndr.JSON(500, e)
		return
	}
	w.WriteHeader(204)
}

func (api *dialogueApi) PutUser(r *http.Request, w http.ResponseWriter, params martini.Params, session sessions.Session, rndr render.Render) {
	username := session.Get("username")
	updateUsername := params["username"]
	// check for admin user
	// if not admin, verify request is updating own account or deny
	if username != "admin" && username != updateUsername {
		e := ApiError{
			Error: "you are not allowed to update this resource",
		}
		log.Warn(fmt.Sprintf("User %s attempted to update user %s", username, updateUsername))
		rndr.JSON(403, e)
		return
	}
	// update user
	user, err := api.rdb.GetUser(updateUsername)
	if err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error updating user: %s", err),
		}
		rndr.JSON(500, e)
		return
	}
	password := r.FormValue("password")
	// hash password
	pw, err := api.auth.HashPassword(password)
	if err != nil {
		e := ApiError{
			Error: "error hashing password",
		}
		rndr.JSON(500, e)
		return

	}
	user.Password = pw
	if err := api.rdb.UpdateUser(user); err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error updating user: %s", err),
		}
		rndr.JSON(500, e)
		return
	}
	log.Info(fmt.Sprintf("User %s updated user %s", username, updateUsername))
	w.WriteHeader(204)
}
