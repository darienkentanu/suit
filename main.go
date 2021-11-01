package main

import (
	"suit/config"
	"suit/routes"
)

func main() {
	config.InitDB()
	e := routes.New()
	e.Logger.Fatal(e.Start(":8000"))
}
