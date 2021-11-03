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

func voucherRoute(e *echo.Echo, db *gorm.DB, dbSQL *sql.DB) {
	vdb := database.NewVoucherDB(db)
	vc := controllers.NewVoucherController(vdb)
	e.GET("/vouchers", vc.GetVouchers)
	e.POST("/vouchers", vc.AddVouchers, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.PUT("/vouchers/:id", vc.EditVouchers, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.DELETE("/vouchers/:id", vc.DeleteVouchers, middlewares.IsLoggedIn, middlewares.IsStaff)
}
