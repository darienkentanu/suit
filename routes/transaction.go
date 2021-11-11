package routes

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func transactionRoute(e *echo.Echo, db *gorm.DB) {
	tdb := database.NewTransactionDB(db)
	cdb := database.NewCategoryDB(db)
	ccdb := database.NewCartDB(db)
	dpdb := database.NewDropPointsDB(db)
	tvc := controllers.NewTransactionController(tdb, cdb, ccdb, dpdb)

	e.GET("/transactions", tvc.GetTransactions, middlewares.IsLoggedIn)
	e.GET("/transactionsbydroppoint/:id", tvc.GetTransactionsDropPoint, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.GET("/transactionreport/:range", tvc.GetTransactionsWithRangeDate, middlewares.IsLoggedIn)	// daily, weekly, monthly

	e.GET("/totaltransaction", tvc.GetTransactionTotal, middlewares.IsLoggedIn)
	e.GET("/totaltransaction/:range", tvc.GetTransactionTotalWithRangeDate, middlewares.IsLoggedIn)	// daily, weekly, monthly
}
