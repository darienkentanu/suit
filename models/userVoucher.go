package models

import "time"

type UserVoucher struct {
	ID        int     `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	VoucherID int     `gorm:"type:int;not null" json:"voucher_id" form:"voucher_id"`
	Voucher   Voucher `gorm:"foreignkey:VoucherID;" json:"-"`
	UserID    int     `gorm:"type:int;not null" json:"users_id" form:"users_id"`
	User      User    `gorm:"foreignkey:UserID;" json:"-"`
	Used      int     `gorm:"type:int;not null;default:null" json:"used" form:"used"`
	UpdatedAt time.Time
	CreatedAt time.Time
}
