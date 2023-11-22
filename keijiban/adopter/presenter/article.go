package presenter

import (
	"github.com/gin-gonic/gin"
	"semi_systems/keijiban/domain"
	"semi_systems/keijiban/usecase"
)

type article struct {
	c *gin.Context
}

type ArticleOutputFactory func(c *gin.Context) usecase.ArticleOutputPort

func NewArticleOutputFactory() ArticleOutputFactory {
	return func(c *gin.Context) usecase.ArticleOutputPort {
		return &article{c: c}
	}
}

func (a article) Create(id uint) error {
	a.c.JSON(200, gin.H{"id": id})
	return nil
}

func (a article) GetAll(res []*domain.Article) error {
	a.c.JSON(200, res)
	return nil
}

func (a article) Update() error {
	a.c.JSON(200, gin.H{})
	return nil
}

func (a article) Delete() error {
	a.c.JSON(200, gin.H{})
	return nil
}
