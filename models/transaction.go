package models

import "time"

type Transaction struct {
	ID           int        `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	UserID       int        `gorm:"type:bigint;not null" json:"user_id" form:"user_id"`
	User         User       `gorm:"foreignkey:UserID" json:"-"`
	Status       int        `gorm:"type:tinyint;not null;default:0;" json:"status" form:"status"`
	Point        int        `gorm:"type:bigint;not null" json:"point" form:"point"`
	Method       string     `gorm:"type:enum('dropoff', 'pickup');" json:"method" form:"method"`
	Drop_PointID int        `gorm:"type:bigint;not null" json:"drop_point_id" form:"drop_point_id"`
	Drop_Point   Drop_Point `gorm:"foreignkey:Drop_PointID;" json:"-"`
	CheckoutID   int        `gorm:"not null" json:"checkout_id" form:"checkout_id"`
	Checkout     Checkout   `gorm:"foreignkey:CheckoutID;" json:"-"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ResponseGetTransactions struct {
	ID               	int    `json:"id"`
	UserID           	int    `json:"user_id"`
	Method           	string `json:"method"`
	DropPointID      	int    `json:"drop_point_id"`
	DropPointAddress 	string `json:"drop_point_address"`
	TotalReceivedPoints	int    `json:"total_received_points"`
	Categories     		[]ResponseGetCategory `json:"categories"`
	Status         		string                `json:"status" form:"status"`
	// TotalPointUsed int                   `json:"total_point_used" form:"total_point_used"`
	CreatedAt      		time.Time
	UpdatedAt      		time.Time
}
