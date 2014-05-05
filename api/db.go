package main

import (
	"time"

	"bitbucket.org/ehazlett/dialogue"
	rdb "github.com/dancannon/gorethink"
)

type (
	Db interface {
		SaveTopic(*dialogue.Topic) error
		DeleteTopic(*dialogue.Topic) error
		GetTopic(id string) (*dialogue.Topic, error)
		SavePost(*dialogue.Post) error
		DeletePost(*dialogue.Post) error
		GetPost(id string) (*dialogue.Post, error)
	}
	rethinkdb struct {
		session *rdb.Session
	}
)

const (
	TOPIC_TABLE = "topic"
	POST_TABLE  = "post"
)

func NewRethinkdbSession(address string, database string) (*rethinkdb, error) {
	var session *rdb.Session
	session, err := rdb.Connect(rdb.ConnectOpts{
		Address:     address,
		Database:    database,
		MaxIdle:     10,
		IdleTimeout: time.Second * 10,
	})

	if err != nil {
		return nil, err
	}

	r := &rethinkdb{
		session: session,
	}
	// initialize database
	rdb.Db(database).TableCreate(TOPIC_TABLE).Run(session)
	rdb.Db(database).TableCreate(POST_TABLE).Run(session)

	return r, nil
}

func (s *rethinkdb) topicExists(title string) bool {
	row, err := rdb.Table(TOPIC_TABLE).Filter(map[string]string{"title": title}).RunRow(s.session)
	if err != nil {
		log.Errorf("Error checking for topic: %s", err)
		return true
	}
	return !row.IsNil()
}

func (s *rethinkdb) SaveTopic(topic *dialogue.Topic) error {
	if !s.topicExists(topic.Title) {
		if _, err := rdb.Table(TOPIC_TABLE).Insert(topic).Run(s.session); err != nil {
			return err
		}
	} else {
		log.Warn("Topic " + topic.Title + " already exists")
	}
	return nil
}

func (s *rethinkdb) UpdateTopic(topic *dialogue.Topic) error {
	if _, err := rdb.Table(TOPIC_TABLE).Update(topic).Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *rethinkdb) DeleteTopic(id string) error {
	if _, err := rdb.Table(TOPIC_TABLE).Filter(map[string]string{"id": id}).Delete().Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *rethinkdb) GetTopic(id string) (*dialogue.Topic, error) {
	res, err := rdb.Table(TOPIC_TABLE).Get(id).RunRow(s.session)
	if err != nil {
		log.Errorf("Unable to get topic from db: %s", err)
		return nil, err
	}
	var topic *dialogue.Topic
	if !res.IsNil() {
		if err := res.Scan(&topic); err != nil {
			log.Errorf("Unable to get topic from db: %s", err)
			return nil, err
		}
	}
	return topic, nil
}

func (s *rethinkdb) GetTopics() ([]*dialogue.Topic, error) {
	var topics []*dialogue.Topic
	res, err := rdb.Table(TOPIC_TABLE).Run(s.session)
	if err != nil {
		log.Errorf("Unable to get topics from db: %s", err)
		return nil, err
	}
	for res.Next() {
		var t *dialogue.Topic
		if err := res.Scan(&t); err != nil {
			log.Errorf("Unable to deserialize topic from db: %s", err)
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

func (s *rethinkdb) SavePost(post *dialogue.Post) error {
	if _, err := rdb.Table(POST_TABLE).Insert(post).Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *rethinkdb) UpdatePost(post *dialogue.Post) error {
	if _, err := rdb.Table(POST_TABLE).Update(post).Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *rethinkdb) DeletePost(post *dialogue.Post) error {
	if _, err := rdb.Table(POST_TABLE).Get(post).Delete().Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *rethinkdb) GetPost(id string) (*dialogue.Post, error) {
	res, err := rdb.Table(POST_TABLE).Get(id).RunRow(s.session)
	if err != nil {
		log.Errorf("Unable to get post from db: %s", err)
		return nil, err
	}
	var post *dialogue.Post
	if !res.IsNil() {
		if err := res.Scan(&post); err != nil {
			log.Errorf("Unable to get post from db: %s", err)
			return nil, err
		}
	}
	return post, nil
}

func (s *rethinkdb) GetPosts(topicId string) ([]*dialogue.Post, error) {
	var posts []*dialogue.Post
	res, err := rdb.Table(POST_TABLE).Run(s.session)
	if err != nil {
		log.Errorf("Unable to get posts from db: %s", err)
		return nil, err
	}
	for res.Next() {
		var p *dialogue.Post
		if err := res.Scan(&p); err != nil {
			log.Errorf("Unable to deserialize post from db: %s", err)
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
