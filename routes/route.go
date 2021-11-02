package routes

import (
	"database/sql"

	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
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

	dbC := database.NewCategoryDB(db)
	cc := controllers.NewCategoryController(dbC)
	e.GET("/categories", cc.GetCategories)
	e.POST("/categories", cc.AddCategories, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.PUT("/categories/:id", cc.EditCategories, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.DELETE("/categories/:id", cc.DeleteCategories, middlewares.IsLoggedIn, middlewares.IsStaff)

	return e
}
