package routes

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func cartRoute(e *echo.Echo, db *gorm.DB) {
	cdb := database.NewCartDB(db)
	cc := controllers.NewCartController(cdb)
	e.POST("/cart", cc.AddToCart, middlewares.IsLoggedIn, middlewares.IsUser)
	e.GET("/cart", cc.GetCartItem, middlewares.IsLoggedIn, middlewares.IsUser)
	e.PUT("/cartitems/:id", cc.EditCartItem, middlewares.IsLoggedIn, middlewares.IsUser)
	e.DELETE("/cartitems/:id", cc.DeleteCartItem, middlewares.IsLoggedIn, middlewares.IsUser)
}
