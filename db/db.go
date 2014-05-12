package db

import (
	"errors"
	"time"

	"bitbucket.org/ehazlett/dialogue"
	"github.com/Sirupsen/logrus"
	rdb "github.com/dancannon/gorethink"
)

type (
	Db interface {
		SaveTopic(*dialogue.Topic) error
		DeleteTopic(string) error
		GetTopic(string) (*dialogue.Topic, error)
		SavePost(*dialogue.Post) error
		DeletePost(string) error
		GetPost(string) (*dialogue.Post, error)
		SaveUser(*dialogue.User) error
		GetUser(string) (*dialogue.User, error)
		DeleteUser(string) error
		GetAuthorization(string) (*dialogue.Authorization, error)
		SaveAuthorization(*dialogue.Authorization) error
	}
	Rethinkdb struct {
		session *rdb.Session
	}
)

var (
	ErrTopicNotFound = errors.New("topic not found")
	ErrPostNotFound  = errors.New("post not found")
	ErrUserExists    = errors.New("user exists")
	ErrTopicExists   = errors.New("topic exists")
	log              = logrus.New()
)

const (
	AUTH_TABLE  = "auth"
	POST_TABLE  = "post"
	TOPIC_TABLE = "topic"
	USER_TABLE  = "user"
)

func NewRethinkdbSession(address string, database string) (*Rethinkdb, error) {
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

	r := &Rethinkdb{
		session: session,
	}
	// initialize database
	rdb.Db(database).TableCreate(AUTH_TABLE).Run(session)
	rdb.Db(database).TableCreate(TOPIC_TABLE).Run(session)
	rdb.Db(database).TableCreate(POST_TABLE).Run(session)
	rdb.Db(database).TableCreate(USER_TABLE).Run(session)
	return r, nil
}

func (s *Rethinkdb) topicExists(title string) bool {
	row, err := rdb.Table(TOPIC_TABLE).Filter(map[string]string{"title": title}).RunRow(s.session)
	if err != nil {
		log.Errorf("Error checking for topic: %s", err)
		return true
	}
	return !row.IsNil()
}

func (s *Rethinkdb) SaveTopic(topic *dialogue.Topic) error {
	if !s.topicExists(topic.Title) {
		if _, err := rdb.Table(TOPIC_TABLE).Insert(topic).Run(s.session); err != nil {
			return err
		}
	} else {
		return ErrTopicExists
	}
	return nil
}

func (s *Rethinkdb) UpdateTopic(topic *dialogue.Topic) error {
	if _, err := rdb.Table(TOPIC_TABLE).Update(topic).Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *Rethinkdb) DeleteTopic(id string) error {
	tbl := rdb.Table(TOPIC_TABLE)
	row, err := tbl.Get(id).RunRow(s.session)
	if err != nil {
		return err
	}
	if row.IsNil() {
		return ErrTopicNotFound
	}
	// delete
	if _, err := tbl.Get(id).Delete().Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *Rethinkdb) GetTopic(id string) (*dialogue.Topic, error) {
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

func (s *Rethinkdb) GetTopics() ([]*dialogue.Topic, error) {
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

func (s *Rethinkdb) SavePost(post *dialogue.Post) error {
	if _, err := rdb.Table(POST_TABLE).Insert(post).Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *Rethinkdb) UpdatePost(post *dialogue.Post) error {
	if _, err := rdb.Table(POST_TABLE).Update(post).Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *Rethinkdb) DeletePost(id string) error {
	tbl := rdb.Table(POST_TABLE)
	row, err := tbl.Get(id).RunRow(s.session)
	if err != nil {
		return err
	}
	if row.IsNil() {
		return ErrPostNotFound
	}
	// delete
	if _, err := tbl.Get(id).Delete().Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *Rethinkdb) GetPost(id string) (*dialogue.Post, error) {
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

func (s *Rethinkdb) GetPosts(topicId string) ([]*dialogue.Post, error) {
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

func (s *Rethinkdb) userExists(username string) bool {
	row, err := rdb.Table(USER_TABLE).Filter(map[string]string{"username": username}).RunRow(s.session)
	if err != nil {
		log.Errorf("Error checking for user: %s", err)
		return true
	}
	return !row.IsNil()
}

func (s *Rethinkdb) SaveUser(user *dialogue.User) error {
	if !s.userExists(user.Username) {
		if _, err := rdb.Table(USER_TABLE).Insert(user).Run(s.session); err != nil {
			return err
		}
	} else {
		return ErrUserExists
	}
	return nil
}

func (s *Rethinkdb) UpdateUser(user *dialogue.User) error {
	if _, err := rdb.Table(USER_TABLE).Update(user).Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *Rethinkdb) GetUser(username string) (*dialogue.User, error) {
	res, err := rdb.Table(USER_TABLE).Filter(map[string]string{"username": username}).RunRow(s.session)
	if err != nil {
		log.Errorf("Unable to get user from db: %s", err)
		return nil, err
	}
	var user *dialogue.User
	if !res.IsNil() {
		if err := res.Scan(&user); err != nil {
			log.Errorf("Unable to get user from db: %s", err)
			return nil, err
		}
	}
	return user, nil
}

func (s *Rethinkdb) DeleteUser(username string) error {
	if _, err := rdb.Table(USER_TABLE).Filter(map[string]string{"username": username}).Delete().Run(s.session); err != nil {
		return err
	}
	return nil
}

func (s *Rethinkdb) GetAuthorization(username string) (*dialogue.Authorization, error) {
	res, err := rdb.Table(AUTH_TABLE).Filter(map[string]string{"username": username}).RunRow(s.session)
	if err != nil {
		log.Errorf("Unable to get user authorization from db: %s", err)
		return nil, err
	}
	var auth *dialogue.Authorization
	if !res.IsNil() {
		if err := res.Scan(&auth); err != nil {
			log.Errorf("Unable to get user authorization from db: %s", err)
			return nil, err
		}
	}
	return auth, nil
}

func (s *Rethinkdb) SaveAuthorization(auth *dialogue.Authorization) error {
	// delete exiting
	rdb.Table(AUTH_TABLE).Filter(map[string]string{"username": auth.Username}).Delete().Run(s.session)
	// add auth
	if _, err := rdb.Table(AUTH_TABLE).Insert(auth).Run(s.session); err != nil {
		return err
	}
	return nil
}
