package main

import (
	"fmt"

	"github.com/darienkentanu/suit/config"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/darienkentanu/suit/routes"
)

// func init() {
// 	config.InitDB()
// 	config.InitDBSQL()
// }

func main() {
	e := routes.New()

	middlewares.LogMiddlewares(e)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.GetInt("port"))))
}
