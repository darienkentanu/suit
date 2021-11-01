package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/darienkentanu/suit/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	initConfig()
}

func InitDB() *gorm.DB {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		GetString("database.username"),
		GetString("database.password"),
		GetString("database.address"),
		GetInt("database.port"),
		GetString("database.name"),
	)

	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	initMigration(db)
	return db
}

func InitDBSQL() *sql.DB {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		GetString("database.username"),
		GetString("database.password"),
		GetString("database.address"),
		GetInt("database.port"),
		GetString("database.name"),
	)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

func initMigration(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Voucher{})
	db.AutoMigrate(&models.User_Voucher{})
	db.AutoMigrate(&models.Drop_Point{})
	db.AutoMigrate(&models.Staff{})
	db.AutoMigrate(&models.Login{})
	db.AutoMigrate(&models.Checkout{})
	db.AutoMigrate(&models.Transaction{})

	db.AutoMigrate(&models.Cart{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.CartItem{})
}
