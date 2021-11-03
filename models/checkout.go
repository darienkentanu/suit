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
	TransactionID    int                   `json:"transaction_id"`
	Method           string                `json:"method"`
	DropPointID      int                   `json:"drop_point_id"`
	DropPointAddress string                `json:"drop_point_address"`
	Distance         float64               `json:"distance"`
	Duration         int                   `json:"duration"`
	TotalPoint       int                   `json:"total_points"`
	Categories       []ResponseGetCategory `json:"categories"`
}

type Checkout_Response_DropOff struct {
	TransactionID    int                   `json:"transaction_id"`
	Method           string                `json:"method"`
	DropPointID      int                   `json:"drop_point_id"`
	DropPointAddress string                `json:"drop_point_address"`
	TotalPoint       int                   `json:"total_points"`
	Categories       []ResponseGetCategory `json:"categories"`
}
