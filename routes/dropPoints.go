package routes

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func dropPointsRoute(e *echo.Echo, db *gorm.DB) {
	dpdb := database.NewDropPointsDB(db)
	dpc := controllers.NewDropPointsController(dpdb)
	e.GET("/droppoints", dpc.GetDropPoints)
	e.POST("/droppoints", dpc.AddDropPoints, middlewares.IsLoggedIn, middlewares.IsStaff) // untuk insert pertama kali middleware is staff perlu dihapus terlebih dahulu.
	e.PUT("/droppoints/:id", dpc.EditDropPoints, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.DELETE("/droppoints/:id", dpc.DeleteDropPoints, middlewares.IsLoggedIn, middlewares.IsStaff)
}
