package models

import "time"

type Admin struct {
	ID           int        `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	Fullname     string     `gorm:"type:varchar(100);not null" json:"fullname" form:"fullname"`
	PhoneNumber  string     `gorm:"type:varchar(15);unique;not null" json:"phone_number" form:"phone_number"`
	Drop_PointID int        `gorm:"type:bigint;not null" json:"drop_point_id" form:"drop_point_id"`
	Drop_Point   Drop_Point `gorm:"foreignkey:Drop_PointID;" json:"-"`
	LoginID      int        `gorm:"type:bigint;not null" json:"login_id" form:"login_id"`
	Login        Login      `gorm:"foreignkey:LoginID;" json:"-"`
	CreatedAt    time.Time
}
