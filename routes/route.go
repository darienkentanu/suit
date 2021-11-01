package routes

import (
	"github.com/darienkentanu/suit/controllers"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()
	e.GET("/categories", controllers.GetCategories)
	e.POST("/categories", controllers.AddCategories)
	e.PUT("/categories/:id", controllers.EditCategories)
	e.DELETE("/categories/:id", controllers.DeleteCategories)
	return e
}
