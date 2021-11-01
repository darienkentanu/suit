package models

import "time"

type Transaction struct {
	ID           int         `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	UserID       int         `gorm:"type:int;not null" json:"user_id" form:"user_id"`
	User         User        `gorm:"foreignkey:UserID" json:"-"`
	Status       int         `gorm:"type:int;not null;default:0;" json:"status" form:"status"`
	Point        int         `gorm:"type:int;not null" json:"point" form:"point"`
	Method       string      `gorm:"type:enum('courier', 'droppoint');" json:"method" form:"method"`
	Drop_PointID int         `gorm:"type:int;not null" json:"drop_point_id" form:"drop_point_id"`
	DropPoint    []DropPoint `gorm:"foreignkey:DropPointID;" json:"-"`
	CheckoutID   int         `gorm:"primarykey;not null" json:"checkout_id" form:"checkout_id"`
	Checkout     Checkout    `gorm:"foreignkey:CheckoutID;" json:"-"`
	CreatedAt    time.Time
}
