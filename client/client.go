package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ehazlett/dialogue"
)

var (
	ErrLoginFailed   = errors.New("invalid username/password")
	ErrCreatingTopic = errors.New("error creating topic")
)

type (
	client struct {
		baseUrl  string
		username string
		token    string
	}
	authResponse struct {
		Token string `json:"token"`
	}
	apiError struct {
		Error string `json:"error"`
	}
)

func Authenticate(baseUrl, username, password string) (string, error) {
	baseUrl = baseUrl + "/auth"
	resp, err := http.PostForm(baseUrl, url.Values{"username": {username}, "password": {password}})
	if err != nil {
		return "", err
	}
	// check for unauth
	if resp.StatusCode == 401 {
		return "", ErrLoginFailed
	}
	var r authResponse
	contents, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	cb := bytes.NewBufferString(string(contents))
	d := json.NewDecoder(cb)
	if err := d.Decode(&r); err != nil {
		return "", err
	}
	return r.Token, nil
}

func NewDialogueClient(url, username, token string) (*client, error) {
	c := &client{
		baseUrl:  url,
		username: username,
		token:    token,
	}
	return c, nil
}

func getApiErrorFromResponse(resp *http.Response) apiError {
	dec := json.NewDecoder(resp.Body)
	var apiErr apiError
	if err := dec.Decode(&apiErr); err != nil {
		apiErr = apiError{
			Error: fmt.Sprintf("Unable to parse response: %s", err),
		}
	}
	return apiErr
}

func (c *client) buildUrl(path string) string {
	return fmt.Sprintf("%s%s", c.baseUrl, path)
}

func (c *client) doRequest(method, path string) (*http.Response, error) {
	url := c.buildUrl(path)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	// add auth headers
	req.Header.Add("X-Auth-User", c.username)
	req.Header.Add("X-Auth-Token", c.token)
	resp, err := client.Do(req)
	return resp, err
}

func (c *client) postRequest(path string, data url.Values) (*http.Response, error) {
	url := c.buildUrl(path)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	// add auth headers
	req.Header.Add("X-Auth-User", c.username)
	req.Header.Add("X-Auth-Token", c.token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	return resp, err
}

func (c *client) GetTopics() ([]*dialogue.Topic, error) {
	var topics []*dialogue.Topic
	resp, err := c.doRequest("GET", "/topics")
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	cb := bytes.NewBufferString(string(contents))
	d := json.NewDecoder(cb)
	if err := d.Decode(&topics); err != nil {
		return nil, err
	}
	return topics, nil
}

func (c *client) DeleteTopic(id string) error {
	resp, err := c.doRequest("DELETE", "/topics/"+id)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		apiErr := getApiErrorFromResponse(resp)
		return errors.New(apiErr.Error)
	}
	return nil
}

func (c *client) CreateTopic(title string) error {
	vals := url.Values{
		"title": {title},
	}
	resp, err := c.postRequest("/topics", vals)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		apiErr := getApiErrorFromResponse(resp)
		return errors.New(apiErr.Error)
	}
	return nil
}

func (c *client) CreatePost(topicId string, content string) error {
	vals := url.Values{
		"content": {content},
	}
	resp, err := c.postRequest("/topics/"+topicId, vals)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		apiErr := getApiErrorFromResponse(resp)
		return errors.New(apiErr.Error)
	}
	return nil
}

func (c *client) GetPosts(topicId string) ([]*dialogue.Post, error) {
	var posts []*dialogue.Post
	resp, err := c.doRequest("GET", "/topics/"+topicId)
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	cb := bytes.NewBufferString(string(contents))
	d := json.NewDecoder(cb)
	if err := d.Decode(&posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (c *client) DeletePost(id string) error {
	resp, err := c.doRequest("DELETE", "/posts/"+id)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		apiErr := getApiErrorFromResponse(resp)
		return errors.New(apiErr.Error)
	}
	return nil
}
