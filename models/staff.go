package models

import "time"

type Staff struct {
	ID           	int        `gorm:"primarykey;AUTO_INCREMENT"`
	Fullname     	string     `gorm:"type:varchar(100);not null"`
	PhoneNumber  	string     `gorm:"type:varchar(15);unique;not null"`
	Drop_PointID 	int        `gorm:"type:bigint;not null"`
	Drop_Point   	Drop_Point `gorm:"foreignkey:Drop_PointID;"`
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
}

type RegisterStaff struct {
	Fullname    string	`json:"fullname" form:"fullname"`
	Email     	string	`json:"email" form:"email"`
	Username  	string	`json:"username" form:"username"`
	Password  	string	`json:"password" form:"password"`
	PhoneNumber string	`json:"phone_number" form:"phone_number"`
	DropPointID	int		`json:"drop_point_id" form:"drop_point_id"`
}

type ResponseGetStaff struct {
	ID					int		`json:"id"`
	Fullname    		string 	`json:"fullname"`
	Email     			string 	`json:"email"`
	Username  			string 	`json:"username"`
	PhoneNumber 		string 	`json:"phone_number"`
	Role				string	`json:"role"`
	DropPointID			int		`json:"drop_point_id"`
	DropPointAddress	string	`json:"drop_point_address"`
}