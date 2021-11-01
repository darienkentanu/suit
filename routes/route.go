package routes

import (
	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	e.POST("/register", controllers.RegisterUsersController)
	e.POST("/login", controllers.LoginController)

	e.GET("/users", controllers.GetAllUsersController, middlewares.IsLoggedIn)
	e.GET("/profile", controllers.GetProfileController, middlewares.IsLoggedIn)
	e.PUT("/profile", controllers.UpdateProfileController, middlewares.IsLoggedIn)

	e.GET("/categories", controllers.GetCategories)
	e.POST("/categories", controllers.AddCategories, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.PUT("/categories/:id", controllers.EditCategories, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.DELETE("/categories/:id", controllers.DeleteCategories, middlewares.IsLoggedIn, middlewares.IsStaff)

	return e
}
