package routes

import (
	"github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func userRoute(e *echo.Echo, db *gorm.DB) {
	udb := database.NewUserDB(db)
	ldb := database.NewLoginDB(db)
	cdb := database.NewCartDB(db)
	sdb := database.NewStaffDB(db)
	dpdb := database.NewDropPointsDB(db)
	uc := controllers.NewUserController(udb, ldb, cdb)
	lc := controllers.NewLoginController(udb, ldb, sdb, dpdb)
	e.POST("/register", uc.RegisterUsers)
	e.POST("/login", lc.Login) // sekalian buat route login staff

	e.GET("/users", uc.GetAllUsers, middlewares.IsLoggedIn, middlewares.IsStaff)
	e.GET("/profile", lc.GetProfile, middlewares.IsLoggedIn)    // bisa untuk staff juga
	e.PUT("/profile", lc.UpdateProfile, middlewares.IsLoggedIn) // bisa untuk staff juga
}
