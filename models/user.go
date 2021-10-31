package models

import "time"

type User struct {
	ID          int    `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	Fullname    string `gorm:"type:varchar(100);not null" json:"fullname" form:"fullname"`
	PhoneNumber string `gorm:"type:varchar(15);unique;not null" json:"phone_number" form:"phone_number"`
	Gender      string `gorm:"type:enum('male', 'female');not null" json:"gender" form:"gender"`
	Address     string `gorm:"type:longtext;not null" json:"address" form:"address"`
	Point       int    `gorm:"type:bigint" json:"point" form:"point"`
	LoginID     int    `gorm:"type:bigint;not null" json:"login_id" form:"login_id"`
	Login       Login  `gorm:"foreignkey:LoginID;" json:"-"`
	CreatedAt   time.Time
}
