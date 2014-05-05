package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/ehazlett/dialogue"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type (
	dialogueApi struct {
		m   *martini.ClassicMartini
		rdb *rethinkdb
	}
	ApiError struct {
		Error string `json:"error"`
	}
)

func NewApi(host string, port int, rdb *rethinkdb) (*dialogueApi, error) {
	m := martini.Classic()

	a := &dialogueApi{
		m:   m,
		rdb: rdb,
	}
	// middleware
	m.Use(render.Renderer())
	// routes
	m.Get("/topics", a.GetTopics)
	m.Post("/topics", a.PostTopics)
	m.Get("/topics/:id", a.GetTopic)
	m.Delete("/topics/:id", a.DeleteTopic)

	return a, nil
}

func (api *dialogueApi) Run() {
	api.m.Run()
}

func (api *dialogueApi) unmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// route handlers
func (api *dialogueApi) GetTopic(params martini.Params, r render.Render) {
	id := params["id"]
	res, err := api.rdb.GetTopic(id)
	if err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error getting topic: %s", err),
		}
		r.JSON(500, e)
		return
	}
	// check for nil
	if res == nil {
		e := ApiError{
			Error: "Topic not found",
		}
		r.JSON(404, e)
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
	topic := &dialogue.Topic{
		Id:    nil,
		Title: title,
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

func (api *dialogueApi) DeleteTopic(w http.ResponseWriter, params martini.Params, rndr render.Render) {
	id := params["id"]
	if err := api.rdb.DeleteTopic(id); err != nil {
		e := ApiError{
			Error: fmt.Sprintf("Error deleting topic: %s", id),
		}
		rndr.JSON(500, e)
		return
	}
	w.WriteHeader(204)
}
