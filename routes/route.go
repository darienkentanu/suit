package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func New(*gorm.DB, *sql.DB) *echo.Echo {
	e := echo.New()
	return e
}
