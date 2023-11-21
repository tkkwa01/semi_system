package domain

import (
	"semi_systems/keijiban/resource/request"
	"semi_systems/packages/context"
)

type Article struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
}

func NewArticle(ctx context.Context, dto *request.ArticleCreate) (*Article, error) {
	var article = Article{
		Title:  dto.Title,
		Text:   dto.Text,
		Author: dto.Author,
	}

	if ctx.IsInValid() {
		return nil, ctx.ValidationError()
	}

	return &article, nil
}
