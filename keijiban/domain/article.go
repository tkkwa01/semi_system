package domain

import (
	"semi_systems/keijiban/resource/request"
	"semi_systems/packages/context"
	"time"
)

type Article struct {
	ID        uint      `json:"id"`
	AuthorID  uint      `json:"author_id"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func NewArticle(ctx context.Context, dto *request.ArticleCreate) (*Article, error) {
	var article = Article{
		AuthorID: dto.AuthorID,
		Author:   dto.Author,
		Title:    dto.Title,
		Text:     dto.Text,
	}

	if ctx.IsInValid() {
		return nil, ctx.ValidationError()
	}

	return &article, nil
}
