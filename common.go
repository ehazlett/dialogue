package dialogue

type (
	Authorization struct {
		Id       *string `json:"id" gorethink:"id"`
		Token    string  `json:"token" gorethink:"token"`
		Username string  `json:"username" gorethink:"username"`
	}
	User struct {
		Id       *string `json:"id" gorethink:"id"`
		Username string  `json:"username" gorethink:"username"`
		Password string  `json:"password" gorethink:"password"`
	}
	Topic struct {
		Id     *string `json:"id" gorethink:"id"`
		Title  string  `json:"title" gorethink:"title"`
		Closed bool    `json:"closed" gorethink:"closed"`
	}
	Post struct {
		Id      *string `json:"id" gorethink:"id"`
		TopicId string  `json:"topicId" gorethink:"topicId"`
		Author  string  `json:"author" gorethink:"author"`
		Content string  `json:"content" gorethink:"content"`
	}
)
