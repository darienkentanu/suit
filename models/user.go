package models

import "time"

type User struct {
	ID          int    `gorm:"primarykey:AUTO_INCREMENT" json:"id" form:"id"`
	Fullname    string `gorm:"type:varchar(100);not null" json:"fullname" form:"fullname"`
	PhoneNumber string `gorm:"type:varchar(15);not null;unique" json:"phone_number" form:"phone_number"`
	Gender      string `gorm:"type:enum('male', 'female');not null" json:"gender" form:"gender"`
	Address     string `gorm:"type:longtext;not null" json:"address" form:"address"`
	Point       string `gorm:"type:int" json:"point" form:"point"`
	LoginID     int    `gorm:"primarykey" json:"login_id" form:"login_id"`
	Login       Login  `gorm:"foreignkey:LoginID;" json:"-"`
	UserVoucher []UserVoucher
	CreatedAt   time.Time
}
