package routes

import (
	// _ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *echo.Echo {
	e := echo.New()

	categoryRoute(e, db)
	userRoute(e, db)
	dropPointsRoute(e, db)
	staffRoute(e, db)
	voucherRoute(e, db)
	cartRoute(e, db)
	userVoucherRoute(e, db)
	transactionRoute(e, db)
	checkout(e, db)

	return e
}
