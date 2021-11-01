package models

import "time"

type User_Voucher struct {
	ID        int     `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	VoucherID int     `gorm:"type:bigint;not null" json:"voucher_id" form:"voucher_id"`
	Voucher   Voucher `gorm:"foreignkey:VoucherID;" json:"-"`
	UserID    int     `gorm:"type:bigint;not null" json:"user_id" form:"user_id"`
	User      User    `gorm:"foreignkey:UserID;" json:"-"`
	Status    int     `gorm:"type:enum('used','available')" json:"status" form:"status"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
