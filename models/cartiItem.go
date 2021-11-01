package models

import "time"

type CartItem struct {
	ID         int      `gorm:"primarykey;not null;AUTO_INCREMENT" json:"id" form:"id"`
	CategoryID int      `gorm:"type:int" json:"category_id" form:"category_id"`
	Category   Category `gorm:"foreignkey:CategoryID" json:"-"`
	Quantity   int      `gorm:"type:int;not null" json:"quantity" form:"quantity"`
	CheckoutID int      `gorm:"type:int" json:"checkout_id" form:"checkout_id"`
	Checkout   Checkout `gorm:"foreignkey:CheckoutID;null" json:"-"`
	CartUserID int      `gorm:"type:int;not null" json:"cart_user_id" form:"cart_user_id"`
	Cart       Cart     `gorm:"foreignkey:CartUserID" json:"-"`
	CreatedAt  time.Time
}
