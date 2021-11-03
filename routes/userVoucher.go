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

func userVoucherRoute(e *echo.Echo, db *gorm.DB, dbSQL *sql.DB) {
	udb := database.NewUserDB(db, dbSQL)
	vdb := database.NewVoucherDB(db)
	uvdb := database.NewUserVoucherDB(db, dbSQL)
	uvc := controllers.NewUserVoucherController(uvdb, udb, vdb)

	e.POST("/claim/:id", uvc.ClaimVoucher, middlewares.IsLoggedIn, middlewares.IsUser)
	e.POST("/redeem/:id", uvc.RedeemVoucher, middlewares.IsLoggedIn, middlewares.IsUser)
	e.GET("/uservouchers", uvc.GetUserVoucher, middlewares.IsLoggedIn, middlewares.IsUser)
}
