package request

type ArticleCreate struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type ArticleUpdate struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}
