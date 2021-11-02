package database

import (
	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type LoginDB struct {
	db *gorm.DB
}

func NewLoginDB(db *gorm.DB) *LoginDB {
	return &LoginDB{db: db}
}

func (ldb *LoginDB) GetEmail(email string) int {
	var login models.Login
	row := ldb.db.Where("email = ?", email).Find(&login).RowsAffected
	return int(row)
}

func (ldb *LoginDB) GetUsername(username string) int {
	var login models.Login
	row := ldb.db.Where("username = ?", username).Find(&login).RowsAffected
	return int(row)
}

func (ldb *LoginDB) CreateLogin(login models.Login) (models.Login, error) {
	if err := ldb.db.Select("email", "username", "password", "role", "user_id").Create(&login).Error; err != nil {
		return login, err
	}

	return login, nil
}

func (ldb *LoginDB) GetAccountByEmailOrUsername(requestLogin models.RequestLogin) (models.Login, error) {
	var login models.Login
	if err := ldb.db.Where("email = ? OR username = ?", requestLogin.Email, requestLogin.Username).First(&login).Error; err != nil {
		return login, err
	}

	return login, nil
}

func (ldb *LoginDB) UpdateToken(id int, newToken string) (models.Login, error) {
	var login models.Login
	if err := ldb.db.First(&login, id).Error; err != nil {
		return login, err
	}

	login.Token = newToken

	if err := ldb.db.Model(&login).Update("token", newToken).Error; err != nil {
		return login, err
	}

	return login, nil
}

func (ldb *LoginDB) UpdateLogin(id int, newLogin models.Login) (models.Login, error) {
	var login models.Login
	if err := ldb.db.First(&login, id).Error; err != nil {
		return login, err
	}

	login.Email 	= newLogin.Email
	login.Username	= newLogin.Username
	login.Password	= newLogin.Password

	if err := ldb.db.Model(&login).Updates(models.Login{
		Email: login.Email,
		Username: login.Username,
		Password: login.Password,
	}).Error; err != nil {
		return login, err
	}

	return login, nil
}