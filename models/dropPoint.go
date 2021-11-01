package models

import "time"

type Drop_Point struct {
	ID        int    `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	Address   string `gorm:"type:longtext;not null" json:"address" form:"address"`
	Longitude string `gorm:"type:longtext" json:"longitute" form:"longitude"`
	Latitude  string `gorm:"type:longtext" json:"latitude" form:"latitude"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
