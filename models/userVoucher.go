package models

import "time"

type User_Voucher struct {
	ID        int     `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	VoucherID int     `gorm:"type:bigint;not null" json:"voucher_id" form:"voucher_id"`
	Voucher   Voucher `gorm:"foreignkey:VoucherID;" json:"-"`
	UserID    int     `gorm:"type:bigint;not null" json:"user_id" form:"user_id"`
	User      User    `gorm:"foreignkey:UserID;" json:"-"`
	Used      int     `gorm:"type:tinyint;not null;default:0" json:"used" form:"used"`
	UpdatedAt time.Time
	CreatedAt time.Time
}
