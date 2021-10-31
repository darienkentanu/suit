package models

type Drop_Point struct {
	ID     int    `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	Alamat string `gorm:"type:longtext;not null" json:"alamat" form:"alamat"`
}
