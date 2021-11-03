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

func cartRoute(e *echo.Echo, db *gorm.DB, dbSQL *sql.DB) {
	cdb := database.NewCartDB(db)
	cc := controllers.NewCartController(cdb)
	e.POST("/cart", cc.AddToCart, middlewares.IsLoggedIn)
	e.GET("/cart", cc.GetCartItem, middlewares.IsLoggedIn)
	e.PUT("/cartitems/:id", cc.EditCartItem, middlewares.IsLoggedIn)
	e.DELETE("/cartitems/:id", cc.DeleteCartItem, middlewares.IsLoggedIn)
}
