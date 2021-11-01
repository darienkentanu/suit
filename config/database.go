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

func InitDBTest() *gorm.DB {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		GetString("database_test.username"),
		GetString("database_test.password"),
		GetString("database_test.address"),
		GetInt("database_test.port"),
		GetString("database_test.name"),
	)

	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	initMigrationTest(db)
	return db
}

func InitDBSQLTest() *sql.DB {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		GetString("database_test.username"),
		GetString("database_test.password"),
		GetString("database_test.address"),
		GetInt("database_test.port"),
		GetString("database_test.name"),
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

func initMigrationTest(db *gorm.DB) {
	db.Migrator().DropTable(&models.User{})
	db.AutoMigrate(&models.User{})
	db.Migrator().DropTable(&models.Voucher{})
	db.AutoMigrate(&models.Voucher{})
	db.Migrator().DropTable(&models.User_Voucher{})
	db.AutoMigrate(&models.User_Voucher{})
	db.Migrator().DropTable(&models.Drop_Point{})
	db.AutoMigrate(&models.Drop_Point{})
	db.Migrator().DropTable(&models.Staff{})
	db.AutoMigrate(&models.Staff{})
	db.Migrator().DropTable(&models.Login{})
	db.AutoMigrate(&models.Login{})
	db.Migrator().DropTable(&models.Checkout{})
	db.AutoMigrate(&models.Checkout{})
	db.Migrator().DropTable(&models.Transaction{})
	db.AutoMigrate(&models.Transaction{})
	db.Migrator().DropTable(&models.Cart{})
	db.AutoMigrate(&models.Cart{})
	db.Migrator().DropTable(&models.Category{})
	db.AutoMigrate(&models.Category{})
	db.Migrator().DropTable(&models.CartItem{})
	db.AutoMigrate(&models.CartItem{})
}
