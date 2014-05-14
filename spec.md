# Objective
To provide a single source of communication that allows tracking of status, todo, and conversations.

# Features
* Forum style messaging
    * Topics
    * Threads
    * Permissions ?

* Ability to track and update item status
    * i.e. status of client interaction, etc.
* CLI for interaction
* Web UI for interaction
* Slack integration

# Design
Dialogue is a client/server model.  This project contains the daemon, simplistic web ui, and a command line tool.

# API
HTTP JSON API

## Endpoints

* `/auth`
    * `POST`: authenticates to the system ; returns an auth token as JSON
* `/topics`
    * `GET`: returns all topics as JSON
    * `POST`: creates a new topic
* `/topics/<id>`
    * `GET`: returns a single topic as JSON
    * `PUT`: updates the topic
    * `DELETE`: deletes the topic
* `/posts`
    * `GET`: returns all posts as JSON
    * `POST`: creates a new post
* `/posts/<id>`
    * `GET`: returns a single post as JSON
    * `PUT`: updates the post
    * `DELETE`: deletes the post
* `/search`
    * `POST`: queries the datastore and returns results as JSON

# Command Line Interface

`dialogue topics` : returns all open topics (paginated?)
    * `--all` : returns all topics (paginated?)

`dialogue topics create` : create a new topic
    * `--title "<title>"` : title of topic

`dialogue posts <topic-id>` : returns all posts for topic

`dialogue posts create --topicId <topic-id>` : create a new post in a topic
    * `--content "<content>"` : content of post

# Web UI
Realtime view of topics and posts.

TODO
