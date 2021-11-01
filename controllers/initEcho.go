package controllers

import (
	"database/sql"

	"github.com/darienkentanu/suit/config"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	db_test    = config.InitDBTest()
	dbSQL_test = config.InitDBSQL()
)

func InitEcho() (*echo.Echo, *gorm.DB, *sql.DB) {
	e := echo.New()
	return e, db_test, dbSQL_test
}
