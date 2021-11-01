package models

import "time"

type Login struct {
	ID        int    `gorm:"primarykey;AUTO_INCREMENT"`
	Email     string `gorm:"type:varchar(55);unique"`
	Username  string `gorm:"type:varchar(55);unique"`
	Password  string `gorm:"type:varchar(255)"`
	Role      string `gorm:"type:enum('staff', 'user')"`
	UserID    int    `gorm:"type:bigint"`
	User      User   `gorm:"foreignkey:UserID;"`
	StaffID   int    `gorm:"type:bigint"`
	Staff     Staff  `gorm:"foreignkey:StaffID;"`
	Token     string `gorm:"type:longtext;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RequestLogin struct {
	Username	string	`json:"username" form:"username"`
	Email 		string 	`json:"email" form:"email"`
	Password 	string	`json:"password" form:"password"` 
}

type ResponseLogin struct {
	Username	string	`json:"username"`
	Email 		string 	`json:"email"`
	Role		string	`json:"role"`
	Token		string	`json:"token"`
}

type ResponseGetStaff struct {
	ID			int		`json:"id"`
	Fullname    string 	`json:"fullname"`
	Email     	string 	`json:"email"`
	Username  	string 	`json:"username"`
	PhoneNumber string 	`json:"phone_number"`
	Role		string	`json:"role"`
	DropPointID	int		`json:"drop_point_id"`
}