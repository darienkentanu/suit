package routes

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func New(db *gorm.DB, dbSQL *sql.DB) *echo.Echo {
	e := echo.New()

	categoryRoute(e, db, dbSQL)
	userRoute(e, db, dbSQL)
	dropPointsRoute(e, db, dbSQL)

	return e
}
