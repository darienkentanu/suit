package main

import (
	"fmt"

	"github.com/darienkentanu/suit/config"
	"github.com/darienkentanu/suit/routes"
)

func main() {
	db1 := config.InitDB()
	db2 := config.InitDBSQL()
	e := routes.New(db1, db2)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.GetInt("port"))))
}
