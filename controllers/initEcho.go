package controllers

import (
	"github.com/darienkentanu/suit/config"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitEcho() (*echo.Echo, *gorm.DB) {
	e := echo.New()
	db_test := config.InitDBTest()
	return e, db_test
}
