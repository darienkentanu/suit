package models

type Login struct {
	ID       int    `gorm:"primarikey;AUTO_INCREMENT" json:"id" form:"id"`
	Email    string `gorm:"type:varchar(55);not null" json:"email" form:"email"`
	Username string `gorm:"type:varchar(55);not null" json:"username" form:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"password" form:"password"`
	Role     string `gorm:"type:enum('admin', 'user');not null" json:"role" form:"role"`
	Token    string `gorm:"type:longtext" json:"token" form:"token"`
}
