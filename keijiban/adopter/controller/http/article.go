package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"semi_systems/keijiban/adopter/presenter"
	"semi_systems/keijiban/resource/request"
	"semi_systems/keijiban/usecase"
	"semi_systems/packages/context"
	"semi_systems/packages/http/middleware"
	"semi_systems/packages/http/router"
	"strconv"
)

type article struct {
	inputFactory  usecase.ArticleInputFactory
	outputFactory func(c *gin.Context) usecase.ArticleOutputPort
	articleRepo   usecase.ArticleRepository
}

func NewArticle(r *router.Router, inputFactory usecase.ArticleInputFactory, outputFactory presenter.ArticleOutputFactory) {
	handler := article{
		inputFactory:  inputFactory,
		outputFactory: outputFactory,
	}

	// 認証のない記事の取得
	r.Get("articles", handler.GetAll)

	// 認証を必要とする記事関連のルーティング
	r.Group("", []gin.HandlerFunc{middleware.Auth(true, true)}, func(r *router.Router) {
		r.Group("articles", nil, func(r *router.Router) {
			r.Post("", handler.Create)
			r.Put(":id", handler.Update)
			r.Delete(":id", handler.Delete)
			r.Get(":id", handler.GetMy)
		})
	})
}

// 記事の作成
func (a article) Create(ctx context.Context, c *gin.Context) error {
	var req request.ArticleCreate

	if !bind(c, &req) {
		return nil
	}

	outputPort := a.outputFactory(c)
	inputPort := a.inputFactory(outputPort)

	return inputPort.Create(ctx, &req)
}

// 全ての記事の取得
func (a article) GetAll(ctx context.Context, c *gin.Context) error {
	outputPort := a.outputFactory(c)
	inputPort := a.inputFactory(outputPort)

	return inputPort.GetAll(ctx)
}

// 記事の更新
func (a article) Update(ctx context.Context, c *gin.Context) error {
	var req request.ArticleUpdate

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return err
	}
	req.ID = uint(id)

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return err
	}

	outputPort := a.outputFactory(c)
	inputPort := a.inputFactory(outputPort)

	return inputPort.Update(ctx, &req)
}

// 記事の削除
func (a article) Delete(ctx context.Context, c *gin.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return err
	}

	outputPort := a.outputFactory(c)
	inputPort := a.inputFactory(outputPort)

	return inputPort.Delete(ctx, uint(id))
}

func (a article) GetMy(ctx context.Context, c *gin.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return err
	}

	outputPort := a.outputFactory(c)
	inputPort := a.inputFactory(outputPort)

	return inputPort.GetMy(ctx, uint(id))
}
