package models

import "time"

type Voucher struct {
	ID          int    `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name" form:"name"`
	Point       int    `gorm:"type:int;not null" json:"point" form:"point"`
	UserVoucher []UserVoucher
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
