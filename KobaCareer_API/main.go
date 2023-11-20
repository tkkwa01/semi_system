package main

import (
	"github.com/labstack/echo/v4/middleware"
	"semi_systems/KobaCareer_API/controller"
	"semi_systems/KobaCareer_API/db"
	"semi_systems/KobaCareer_API/repository"
	"semi_systems/KobaCareer_API/router"
	"semi_systems/KobaCareer_API/usecase"
	"semi_systems/KobaCareer_API/validator"
)

func main() {
	db := db.NewDB()
	internValidate := validator.NewInternshipValidator()
	internRepository := repository.NewInternshipRepository(db)
	internUsecase := usecase.NewInternshipUsecase(internRepository, internValidate)
	internController := controller.NewInternshipController(internUsecase)
	e := router.NewRouter(internController)
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":8080"))
}
