package routes

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func dropPointsRoute(e *echo.Echo, db *gorm.DB, dbSQL *sql.DB) {
	dpdb := database.NewDropPointsDB(db)
	dpc := controllers.NewDropPointsController(dpdb)
	e.GET("/droppoints", dpc.GetDropPoints)
	e.POST("/droppoints", dpc.AddDropPoints, middlewares.IsLoggedIn)
	e.PUT("/droppoints/:id", dpc.EditDropPoints, middlewares.IsLoggedIn)
	e.DELETE("/droppoints/:id", dpc.DeleteDropPoints, middlewares.IsLoggedIn)
}
