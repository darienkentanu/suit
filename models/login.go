package models

type Login struct {
	ID       int    `gorm:"primarykey;AUTO_INCREMENT" json:"id" form:"id"`
	Email    string `gorm:"type:varchar(55);unique" json:"email" form:"email"`
	Username string `gorm:"type:varchar(55);unique" json:"username" form:"username"`
	Password string `gorm:"type:varchar(255)" json:"password" form:"password"`
	Role     string `gorm:"type:enum('admin', 'user')" json:"role" form:"role"`
	Token    string `gorm:"type:longtext;" json:"token" form:"token"`
}
