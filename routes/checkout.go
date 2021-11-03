package routes

import (
	_ "github.com/go-sql-driver/mysql"
)

// func checkout(e *echo.Echo, db *gorm.DB, dbSQL *sql.DB) {
// 	cdb := database.newCheckoutDB(db)
// 	cc := controllers.newCheckoutController(cdb)
// 	e.POST("/checkoutbydropoff", cc.checkoutByDropOff, middlewares.IsLoggedIn, middlewares.IsUser)
// 	e.POST("/checkoutbypickup", cc.checkoutByPickUp, middlewares.IsLoggedIn, middlewares.IsUser)
// 	e.PUT("/verification/:id", cc.verification, middlewares.IsLoggedIn, middlewares.IsStaff)
// }
