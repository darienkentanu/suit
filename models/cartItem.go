package models

import "time"

type CartItem struct {
	ID         int      `gorm:"primarykey;not null;AUTO_INCREMENT" json:"id" form:"id"`
	CategoryID int      `gorm:"type:bigint" json:"category_id" form:"category_id"`
	Category   Category `gorm:"foreignkey:CategoryID" json:"-"`
	Weight     int      `gorm:"type:int;not null" json:"weight" form:"weight"`
	CheckoutID int      `gorm:"type:bigint" json:"checkout_id" form:"checkout_id"`
	Checkout   Checkout `gorm:"foreignkey:CheckoutID" json:"-"`
	CartUserID int      `gorm:"type:bigint;not null" json:"cart_user_id" form:"cart_user_id"`
	Cart       Cart     `gorm:"foreignkey:CartUserID" json:"-"`
	CreatedAt  time.Time
}
