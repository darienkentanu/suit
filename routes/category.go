package routes

import (
	// _ "github.com/go-sql-driver/mysql"

	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func categoryRoute(e *echo.Echo, db *gorm.DB) {
	cdb := database.NewCategoryDB(db)
	cc := controllers.NewCategoryController(cdb)
	e.GET("/categories", cc.GetCategories)
	e.POST("/categories", cc.AddCategories, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.PUT("/categories/:id", cc.EditCategories, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.DELETE("/categories/:id", cc.DeleteCategories, middlewares.IsLoggedIn, middlewares.IsStaff)
}
