package main

import (
	"github.com/darienkentanu/suit/config"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/darienkentanu/suit/routes"
)

func main() {
	db := config.InitDB()
	dbSQL := config.InitDBSQL()
	e := routes.New(db, dbSQL)

	middlewares.LogMiddlewares(e)
	e.Logger.Fatal(e.Start(":8000"))
}
