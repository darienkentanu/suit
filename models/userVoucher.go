package models

import "time"

type User_Voucher struct {
	ID        int     `gorm:"primarykey;AUTO_INCREMENT"`
	VoucherID int     `gorm:"type:bigint;not null"`
	Voucher   Voucher `gorm:"foreignkey:VoucherID;"`
	UserID    int     `gorm:"type:bigint;not null"`
	User      User    `gorm:"foreignkey:UserID;"`
	Status    string  `gorm:"type:enum('used','available')"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RequestUserVoucher struct {
	VoucherID	int		`json:"voucher_id" form:"voucher_id"`
}

type ResponseGetUserVoucher struct {
	ID			int		`json:"id"`
	UserID		int		`json:"user_id"`
	VoucherID	int		`json:"voucher_id"`
	VoucherName	string	`json:"voucher_name"`
	Point		int		`json:"point"`
	Status		string	`json:"status"`
}