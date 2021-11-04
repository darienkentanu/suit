package models

import "time"

type Category struct {
	ID        int    `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	Name      string `gorm:"type:varchar(255);unique;not null" json:"name" form:"name"`
	Point     int    `gorm:"type:bigint;not null" json:"point" form:"point"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category_Response struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Point int    `json:"point"`
}

type ResponseGetCategory struct {
	ID        		int    `json:"category_id"`
	Name      		string `json:"category_name"`
	Point     		int    `json:"category_point"`
	Weight    		int    `json:"weight"`
	ReceivedPoints 	int    `json:"received_points"`
}
