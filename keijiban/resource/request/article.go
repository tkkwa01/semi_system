package request

type ArticleCreate struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
}

type ArticleUpdate struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
}
