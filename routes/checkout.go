package routes

import (
	"database/sql"

	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func checkout(e *echo.Echo, db *gorm.DB, dbSQL *sql.DB) {
	cdb := database.NewCheckoutDB(db)
	crdb := database.NewCartDB(db)
	ctdb := database.NewCategoryDB(db)
	dpdb := database.NewDropPointsDB(db)
	udb := database.NewUserDB(db, dbSQL)
	tdb := database.NewTransactionDB(db, dbSQL)

	cc := controllers.NewCheckoutController(cdb, crdb, ctdb, dpdb, udb, tdb)
	e.POST("/checkoutbypickup", cc.CreateCheckoutPickup, middlewares.IsLoggedIn, middlewares.IsUser)
	e.POST("/checkoutbydropoff", cc.CreateCheckoutDropOff, middlewares.IsLoggedIn, middlewares.IsUser)
	// e.PUT("/verification/:id", cc.verification, middlewares.IsLoggedIn, middlewares.IsStaff)
}
