package request

type ArticleCreate struct {
	AuthorID uint   `json:"author_id"`
	Author   string `json:"author"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}

type ArticleUpdate struct {
	ID       uint   `json:"id"`
	AuthorID uint   `json:"author_id"`
	Author   string `json:"author"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}
