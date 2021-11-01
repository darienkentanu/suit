package models

import "time"

type Admins struct {
	ID          int       `gorm:"primarikey;AUTO_INCREMENT" json:"id" form:"id"`
	Fullname    string    `gorm:"type:varchar(100);not null" json:"fullname" form:"fullname"`
	PhoneNumber int       `gorm:"type:varchar(15);unique;not null" json:"phone_number" form:"phone_number"`
	DropPointID int       `gorm:"primarykey" json:"drop_point_id" form:"drop_point_id"`
	DropPoint   DropPoint `gorm:"foreignkey:DropPointID;" json:"-"`
	LoginID     int       `gorm:"primarykey" json:"login_id" form:"login_id"`
	Login       Login     `gorm:"foreignkey:LoginID;" json:"-"`
	CreatedAt   time.Time
}

type DropPoint struct {
	ID     int    `gorm:"primarikey;AUTO_INCREMENT" json:"id" form:"id"`
	Alamat string `gorm:"type:longtext;not null" json:"alamat" form:"alamat"`
	Admins []Admins
}
