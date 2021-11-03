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

func transactionRoute(e *echo.Echo, db *gorm.DB, dbSQL *sql.DB) {
	tdb := database.NewTransactionDB(db, dbSQL)
	cdb := database.NewCategoryDB(db)
	ccdb := database.NewCartDB(db)
	tvc := controllers.NewTransactionController(tdb, cdb, ccdb)

	e.GET("/transactions", tvc.GetTransactions, middlewares.IsLoggedIn, middlewares.IsStaff)
}
