package models

import "time"

type User struct {
	ID          int    `gorm:"primarykey;AUTO_INCREMENT"`
	Fullname    string `gorm:"type:varchar(100);not null"`
	PhoneNumber string `gorm:"type:varchar(15);unique;not null"`
	Gender      string `gorm:"type:enum('male', 'female');not null"`
	Address     string `gorm:"type:longtext;not null"`
	Point       int    `gorm:"type:bigint;not null;default:0"`
	Longitude 	string `gorm:"type:longtext"`
	Latitude  	string `gorm:"type:longtext"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}

type RegisterUser struct {
	Fullname    string `json:"fullname" form:"fullname"`
	Email     	string `json:"email" form:"email"`
	Username  	string `json:"username" form:"username"`
	Password  	string `json:"password" form:"password"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Gender      string `json:"gender" form:"gender"`
	Address     string `json:"address" form:"address"`
}

type ResponseGetUser struct {
	ID			int		`json:"id"`
	Fullname    string 	`json:"fullname"`
	Email     	string 	`json:"email"`
	Username  	string 	`json:"username"`
	PhoneNumber string 	`json:"phone_number"`
	Gender      string 	`json:"gender"`
	Address     string 	`json:"address"`
	Role		string	`json:"role"`
}