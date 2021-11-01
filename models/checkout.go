package models

import "time"

type Checkout struct {
	ID        int `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	CreatedAt time.Time
}
