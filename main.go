package main

import (
	"fmt"

	"github.com/darienkentanu/suit/config"
	"github.com/darienkentanu/suit/routes"
)

func init() {
	config.InitDB()
	config.InitDBSQL()
}

func main() {

	e := routes.New()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.GetInt("port"))))
}
