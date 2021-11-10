package routes

import (
	// "database/sql"

	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func staffRoute(e *echo.Echo, db *gorm.DB) {
	sdb := database.NewStaffDB(db)
	ldb := database.NewLoginDB(db)
	sc := controllers.NewStaffController(sdb, ldb)

	e.POST("/registerstaff", sc.AddStaff)
	e.GET("/staff", sc.GetAllStaff, middlewares.IsLoggedIn, middlewares.IsStaff)
}
