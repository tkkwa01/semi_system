package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"semi_systems/attendance/repository"
	"semi_systems/attendance/usecase"
	"semi_systems/config"
	"semi_systems/driver"
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
	attendanceRepository := repository.NewAttendanceRepository()

	//usecase
	attendanceUsecase := usecase.NewAttendanceUsecase(attendanceRepository)

	//controller
	attendanceController := handler.NewAttendanceHandler(attendanceUsecase)

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
