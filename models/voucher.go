package models

import "time"

type Voucher struct {
	ID        int    `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	Name      string `gorm:"type:varchar(100);unique;not null" json:"name" form:"name"`
	Point     int    `gorm:"type:bigint;not null" json:"point" form:"point"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Voucher_Response struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Point int    `json:"point"`
}
