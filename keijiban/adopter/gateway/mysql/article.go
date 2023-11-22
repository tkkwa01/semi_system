package mysql

import (
	"semi_systems/keijiban/domain"
	"semi_systems/keijiban/usecase"
	"semi_systems/packages/context"
)

type article struct{}

func NewArticleRepository() usecase.ArticleRepository {
	return &article{}
}

func (a article) Create(ctx context.Context, article *domain.Article) (uint, error) {
	db := ctx.DB()

	if err := db.Create(article).Error; err != nil {
		return 0, dbError(err)
	}
	return article.ID, nil
}

func (a article) GetAll(ctx context.Context) ([]*domain.Article, error) {
	db := ctx.DB()

	var articles []*domain.Article
	if err := db.Find(&articles).Error; err != nil {
		return nil, dbError(err)
	}
	return articles, nil
}

func (a article) Update(ctx context.Context, article *domain.Article) error {
	db := ctx.DB()

	if err := db.Save(article).Error; err != nil {
		return dbError(err)
	}
	return nil
}

func (a article) Delete(ctx context.Context, id uint) error {
	db := ctx.DB()

	var article domain.Article
	res := db.Where("id = ?", id).Delete(&article)
	if res.Error != nil {
		return dbError(res.Error)
	}
	return nil
}

func (a article) GetByID(ctx context.Context, id uint) (*domain.Article, error) {
	db := ctx.DB()

	var article domain.Article
	err := db.Where(&domain.Article{ID: id}).First(&article).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &article, nil
}
