package dialogue

type (
	Authorization struct {
		Token    string `json:"token" gorethink:"token"`
		Username string `json:"username" gorethink:"username"`
	}
	User struct {
		Username string `json:"username" gorethink:"username"`
		Password string `json:"password" gorethink:"password"`
	}
	Topic struct {
		Id    *string `json:"id" gorethink:"id"`
		Title string  `json:"title" gorethink:"title"`
	}
	Post struct {
		Id      *string `json:"id" gorethink:"id"`
		TopicId string  `json:"topicId" gorethink:"topicId"`
		Author  string  `json:"author" gorethink:"author"`
		Content string  `json:"content" gorethink:"content"`
	}
)
