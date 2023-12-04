package api

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	httpController "semi_systems/attendance/adopter/controller/http"
	attendanceMysqlRepository "semi_systems/attendance/adopter/gateway/mysql"
	"semi_systems/attendance/adopter/presenter"
	"semi_systems/attendance/usecase"
	"semi_systems/chat"
	"semi_systems/config"
	"semi_systems/driver"
	userHttpController "semi_systems/keijiban/adopter/controller/http"
	keijibanMysqlRepository "semi_systems/keijiban/adopter/gateway/mysql"
	userPresenter "semi_systems/keijiban/adopter/presenter"
	userUsecase "semi_systems/keijiban/usecase"
	"semi_systems/packages/http/middleware"
	"semi_systems/packages/http/router"
	"semi_systems/packages/log"
	"syscall"
	"time"
)

func Execute(hub *chat.Hub) {
	logger := log.Logger()
	defer logger.Sync()

	engine := gin.New()

	engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	engine.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})

	// cors
	engine.Use(middleware.Cors(nil))

	store := cookie.NewStore([]byte("secret"))
	engine.Use(sessions.Sessions("user", store))

	r := router.New(engine, driver.GetRDB)

	//mysql
	attendanceRepository := attendanceMysqlRepository.NewAttendanceRepository()
	userRepository := keijibanMysqlRepository.NewUserRepository()
	articleRepository := keijibanMysqlRepository.NewArticleRepository()

	//usecase
	attendanceInputFactory := usecase.NewAttendanceInputFactory(attendanceRepository)
	attendanceOutputFactory := presenter.NewAttendanceOutputFactory()
	userInputFactory := userUsecase.NewUserInputFactory(userRepository)
	userOutputFactory := userPresenter.NewUserOutputFactory()
	articleInputFactory := userUsecase.NewArticleInputFactory(articleRepository)
	articleOutputFactory := userPresenter.NewArticleOutputFactory()

	//controller
	httpController.NewAttendance(r, attendanceInputFactory, attendanceOutputFactory)
	userHttpController.NewUser(r, userInputFactory, userOutputFactory)
	userHttpController.NewArticle(r, articleInputFactory, articleOutputFactory)

	//serve
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Env.Port),
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
			// logger.Error(fmt.Sprintf("Server forced to shutdown: %+v", err))
		}
	}()

	logger.Info("Succeeded in listen and serve.")

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Server forced to shutdown: %+v", err))
	}

	logger.Info("Server existing")
}

func serveWs(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// ここでクライアントのオリジンを検証します。
			// 安全でないがテスト目的であれば、すべてのオリジンを許可することができます。
			// 本番環境では、特定のオリジンのみを許可するように厳密に設定してください。
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to set websocket upgrade: %+v", err))
		return
	}
	client := chat.NewClient(conn)
	hub.RegisterCh <- client

	go client.WriteLoop()
	go client.ReadLoop(hub)
}
