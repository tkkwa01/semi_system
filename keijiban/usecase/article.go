package usecase

import (
	"errors"
	"log"
	"semi_systems/keijiban/domain"
	"semi_systems/keijiban/resource/request"
	"semi_systems/packages/context"
)

type ArticleInputPort interface {
	Create(ctx context.Context, req *request.ArticleCreate) error
	GetAll(ctx context.Context) error
	Update(ctx context.Context, req *request.ArticleUpdate) error
	Delete(ctx context.Context, id uint) error
	GetMy(ctx context.Context) error
}
type ArticleOutputPort interface {
	Create(id uint) error
	GetAll(res []*domain.Article) error
	Update(res *domain.Article) error
	Delete() error
	GetMy(res *domain.Article) error
}

type ArticleRepository interface {
	Create(ctx context.Context, article *domain.Article) (uint, error)
	GetAll(ctx context.Context) ([]*domain.Article, error)
	Update(ctx context.Context, article *domain.Article) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*domain.Article, error)
	GetMy(ctx context.Context, id uint) error
	GetByUserID(ctx context.Context, userID uint) ([]*domain.Article, error)
}

type Article struct {
	articleRepo ArticleRepository
	outputPort  ArticleOutputPort
}

type ArticleInputFactory func(outputPort ArticleOutputPort) ArticleInputPort

func NewArticleInputFactory(ar ArticleRepository) ArticleInputFactory {
	return func(outputPort ArticleOutputPort) ArticleInputPort {
		return &Article{
			articleRepo: ar,
			outputPort:  outputPort,
		}
	}
}

func (a Article) Create(ctx context.Context, req *request.ArticleCreate) error {
	userID := ctx.UID()
	userName := ctx.UserName()

	log.Printf("Creating article with userID: %d, userName: %s\n", userID, userName)

	newArticle := domain.Article{
		AuthorID: userID,
		Author:   userName,
		Title:    req.Title,
		Text:     req.Text,
	}

	if ctx.IsInValid() {
		return ctx.ValidationError()
	}

	id, err := a.articleRepo.Create(ctx, &newArticle)
	if err != nil {
		return err
	}

	return a.outputPort.Create(id)
}

func (a Article) GetAll(ctx context.Context) error {
	article, err := a.articleRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	return a.outputPort.GetAll(article)
}

func (a Article) Update(ctx context.Context, req *request.ArticleUpdate) error {
	currentUserID := ctx.UID()

	article, err := a.articleRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if article.AuthorID != currentUserID {
		return errors.New("unauthorized: you are not the author of this article")
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Text != "" {
		article.Text = req.Text
	}

	err = a.articleRepo.Update(ctx, article)
	if err != nil {
		return err
	}
	return a.outputPort.Update(article)
}

func (a Article) Delete(ctx context.Context, id uint) error {
	currentUserID := ctx.UID()

	article, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if article.AuthorID != currentUserID {
		return errors.New("unauthorized: you are not the author of this article")
	}

	err = a.articleRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return a.outputPort.Delete()
}

func (a Article) GetMy(ctx context.Context) error {
	currentUserID := ctx.UID()

	article, err := a.articleRepo.GetByUserID(ctx, currentUserID)
	if err != nil {
		return err
	}

	return a.outputPort.GetAll(article)
}
