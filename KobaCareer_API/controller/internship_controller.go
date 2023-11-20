package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"semi_systems/KobaCareer_API/domain"
	"semi_systems/KobaCareer_API/usecase"
	"strconv"
)

type IInternshipController interface {
	GetAllInternships(c echo.Context) error
	GetInternshipById(c echo.Context) error
	CreateInternship(c echo.Context) error
	UpdateInternship(c echo.Context) error
	DeleteInternship(c echo.Context) error
}

type internshipController struct {
	iu usecase.IInternshipUsecase
}

func NewInternshipController(iu usecase.IInternshipUsecase) IInternshipController {
	return &internshipController{iu}
}

func (i internshipController) GetAllInternships(c echo.Context) error {
	internshipsRes, err := i.iu.GetAllInternships()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, internshipsRes)
}

func (i internshipController) GetInternshipById(c echo.Context) error {
	id := c.Param("internshipId")
	// パスパラメータは文字列で与えられるためinternshipIdをuint型に変換
	internshipId, _ := strconv.Atoi(id)
	internshipRes, err := i.iu.GetInternshipById(uint(internshipId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, internshipRes)
}

func (i internshipController) CreateInternship(c echo.Context) error {
	internship := domain.Internships{}
	if err := c.Bind(&internship); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	internshipRes, err := i.iu.CreateInternship(internship)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, internshipRes)
}

func (i internshipController) UpdateInternship(c echo.Context) error {
	internship := domain.Internships{}
	if err := c.Bind(&internship); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	id := c.Param("internshipId")
	internshipId, _ := strconv.Atoi(id)
	internshipRes, err := i.iu.UpdateInternship(internship, uint(internshipId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, internshipRes)
}

func (i internshipController) DeleteInternship(c echo.Context) error {
	id := c.Param("internshipId")
	internshipId, _ := strconv.Atoi(id)

	err := i.iu.DeleteInternship(uint(internshipId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Internship deleted successfully")
}
