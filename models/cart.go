package models

import "time"

type Cart struct {
	UserID    int  `gorm:"primarykey" json:"users_id" form:"users_id"`
	User      User `gorm:"primarykey;foreignkey:UserID;" json:"-"`
	CreatedAt time.Time
}
