package main

import (
	"github.com/darienkentanu/suit/middlewares"
	"github.com/darienkentanu/suit/routes"
)

func main() {
	e := routes.New()

	middlewares.LogMiddlewares(e)
	e.Logger.Fatal(e.Start(":8000"))
}
