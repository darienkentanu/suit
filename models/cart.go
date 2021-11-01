package models

import "time"

type Cart struct {
	UserID    int  `gorm:"primarykey" json:"user_id" form:"user_id"`
	User      User `gorm:"primarykey;foreignkey:UserID;" json:"-"`
	CreatedAt time.Time
}
