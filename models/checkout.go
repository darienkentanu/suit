package models

import "time"

type Checkout struct {
	ID        int `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	CreatedAt time.Time
}

type Checkout_Input_PickUp struct {
	CategoryID []int `json:"category_id" form:"category_id"`
}

type Checkout_Input_DropOff struct {
	CategoryID  []int `json:"category_id" form:"category_id"`
	DropPointID int   `json:"drop_point_id" form:"drop_point_id"`
}

type Checkout_Response_PickUp struct {
	Distance int
}

type Checkout_Response_DropOff struct {
}
