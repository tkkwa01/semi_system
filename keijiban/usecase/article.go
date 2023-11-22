package usecase

import (
	"semi_systems/keijiban/domain"
	"semi_systems/keijiban/resource/request"
	"semi_systems/packages/context"
)

type ArticleInputPort interface {
	Create(ctx context.Context, req *request.ArticleCreate) error
	GetAll(ctx context.Context) error
	Update(ctx context.Context, req *request.ArticleUpdate) error
	Delete(ctx context.Context, id uint) error
}
type ArticleOutputPort interface {
	Create(id uint) error
	GetAll(res []*domain.Article) error
	Update() error
	Delete() error
}

type ArticleRepository interface {
	Create(ctx context.Context, article *domain.Article) (uint, error)
	GetAll(ctx context.Context) ([]*domain.Article, error)
	Update(ctx context.Context, article *domain.Article) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*domain.Article, error)
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
	newArticle, err := domain.NewArticle(ctx, req)
	if err != nil {
		return err
	}

	if ctx.IsInValid() {
		return ctx.ValidationError()
	}

	id, err := a.articleRepo.Create(ctx, newArticle)
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
	article, err := a.articleRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
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
	return a.outputPort.Update()
}

func (a Article) Delete(ctx context.Context, id uint) error {
	err := a.articleRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return a.outputPort.Delete()
}
