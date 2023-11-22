package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	httpController "semi_systems/attendance/adopter/controller/http"
	attendanceMysqlRepository "semi_systems/attendance/adopter/gateway/mysql"
	"semi_systems/attendance/adopter/presenter"
	"semi_systems/attendance/usecase"
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

func Execute() {
	logger := log.Logger()
	defer logger.Sync()

	engine := gin.New()

	engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// cors
	engine.Use(middleware.Cors(nil))

	r := router.New(engine, driver.GetRDB)

	//mysql
	attendanceRepository := attendanceMysqlRepository.NewAttendanceRepository()
	userRepository := keijibanMysqlRepository.NewUserRepository()

	//usecase
	attendanceInputFactory := usecase.NewAttendanceInputFactory(attendanceRepository)
	attendanceOutputFactory := presenter.NewAttendanceOutputFactory()
	userInputFactory := userUsecase.NewUserInputFactory(userRepository)
	userOutputFactory := userPresenter.NewUserOutputFactory()

	//controller
	httpController.NewAttendance(r, attendanceInputFactory, attendanceOutputFactory)
	userHttpController.NewUser(r, userInputFactory, userOutputFactory)

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
