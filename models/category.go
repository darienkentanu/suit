package models

import "time"

type Category struct {
	ID        int    `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	Name      string `gorm:"type:varchar(255);unique;not null" json:"name" form:"name"`
	Point     int    `gorm:"type:bigint;not null" json:"point" form:"point"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
