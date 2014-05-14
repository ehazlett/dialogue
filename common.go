package dialogue

type (
	Authorization struct {
		Id       string `json:"id" gorethink:"id,omitempty"`
		Token    string `json:"token" gorethink:"token"`
		Username string `json:"username" gorethink:"username"`
	}
	User struct {
		Id       string `json:"id" gorethink:"id,omitempty"`
		Username string `json:"username" gorethink:"username"`
		Password string `json:"password" gorethink:"password"`
	}
	Topic struct {
		Id     string `json:"id" gorethink:"id,omitempty"`
		Title  string `json:"title" gorethink:"title"`
		Closed bool   `json:"closed" gorethink:"closed"`
	}
	Post struct {
		Id      string `json:"id" gorethink:"id,omitempty"`
		TopicId string `json:"topicId" gorethink:"topicId"`
		Author  string `json:"author" gorethink:"author"`
		Content string `json:"content" gorethink:"content"`
	}
)
