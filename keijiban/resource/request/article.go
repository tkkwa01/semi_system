package request

type ArticleCreate struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

type ArticleUpdate struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}
