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

func userRoute(e *echo.Echo, db *gorm.DB, dbSQL *sql.DB) {
	udb := database.NewUserDB(db, dbSQL)
	ldb := database.NewLoginDB(db)
	cdb := database.NewCartDB(db)
	uc := controllers.NewUserController(udb, ldb, cdb)
	lc := controllers.NewLoginController(udb, ldb)
	e.POST("/register", uc.RegisterUsers)
	e.POST("/login", lc.Login)

	e.GET("/users", uc.GetAllUsers, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.GET("/profile", lc.GetProfile, middlewares.IsLoggedIn)
	e.PUT("/profile", lc.UpdateProfile, middlewares.IsLoggedIn)
}
