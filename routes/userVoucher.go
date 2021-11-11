package routes

import (
	// _ "github.com/go-sql-driver/mysql"

	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func userVoucherRoute(e *echo.Echo, db *gorm.DB) {
	udb := database.NewUserDB(db)
	vdb := database.NewVoucherDB(db)
	uvdb := database.NewUserVoucherDB(db)
	uvc := controllers.NewUserVoucherController(uvdb, udb, vdb)

	e.POST("/claim/:id", uvc.ClaimVoucher, middlewares.IsLoggedIn, middlewares.IsUser)
	e.PUT("/redeem/:id", uvc.RedeemVoucher, middlewares.IsLoggedIn, middlewares.IsUser)
	e.GET("/uservouchers", uvc.GetUserVoucher, middlewares.IsLoggedIn, middlewares.IsUser)
}
