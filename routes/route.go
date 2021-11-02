package routes

import (
	"database/sql"

	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func New(db *gorm.DB, dbSQL *sql.DB) *echo.Echo {
	e := echo.New()

	e.POST("/register", controllers.RegisterUsersController)
	e.POST("/login", controllers.LoginController)

	e.GET("/users", controllers.GetAllUsersController, middlewares.IsLoggedIn)
	e.GET("/profile", controllers.GetProfileController, middlewares.IsLoggedIn)
	e.PUT("/profile", controllers.UpdateProfileController, middlewares.IsLoggedIn)

	categoryRoute(e, db, dbSQL)

	return e
}
