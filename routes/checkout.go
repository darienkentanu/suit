package routes

import (
	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func checkout(e *echo.Echo, db *gorm.DB) {
	cdb := database.NewCheckoutDB(db)
	crdb := database.NewCartDB(db)
	ctdb := database.NewCategoryDB(db)
	dpdb := database.NewDropPointsDB(db)
	udb := database.NewUserDB(db)
	tdb := database.NewTransactionDB(db)
	ldb := database.NewLoginDB(db)

	cc := controllers.NewCheckoutController(cdb, crdb, ctdb, dpdb, udb, tdb, ldb)
	e.POST("/checkoutbypickup", cc.CreateCheckoutPickup, middlewares.IsLoggedIn, middlewares.IsUser)
	e.POST("/checkoutbydropoff", cc.CreateCheckoutDropOff, middlewares.IsLoggedIn, middlewares.IsUser)
	e.PUT("/verification/:id", cc.Verification, middlewares.IsLoggedIn, middlewares.IsStaff)
}
