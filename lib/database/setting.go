package database

import "github.com/darienkentanu/suit/config"

type M map[string]interface{}

var (
	db    = config.InitDB()
	dbSQL = config.InitDBSQL()
)
