package database

import (
	"github.com/darienkentanu/suit/config"
	"github.com/darienkentanu/suit/models"
)

func GetEmail(email string) int {
	var login models.Login
	row := config.InitDB().Where("email = ?", email).Find(&login).RowsAffected
	return int(row)
}

func GetUsername(username string) int {
	var login models.Login
	row := config.InitDB().Where("username = ?", username).Find(&login).RowsAffected
	return int(row)
}

func CreateLogin(login models.Login) (models.Login, error) {
	if err := config.InitDB().Select("email", "username", "password", "role", "user_id").Create(&login).Error; err != nil {
		return login, err
	}

	return login, nil
}

func GetAccountByEmailOrUsername(requestLogin models.RequestLogin) (models.Login, error) {
	var login models.Login
	if err := config.InitDB().Where("email = ? OR username = ?", requestLogin.Email, requestLogin.Username).First(&login).Error; err != nil {
		return login, err
	}

	return login, nil
}

func UpdateToken(id int, newToken string) (models.Login, error) {
	var login models.Login
	if err := config.InitDB().First(&login, id).Error; err != nil {
		return login, err
	}

	login.Token = newToken

	if err := config.InitDB().Model(&login).Update("token", newToken).Error; err != nil {
		return login, err
	}

	return login, nil
}

func UpdateLogin(id int, newLogin models.Login) (models.Login, error) {
	var login models.Login
	if err := config.InitDB().First(&login, id).Error; err != nil {
		return login, err
	}

	login.Email 	= newLogin.Email
	login.Username	= newLogin.Username
	login.Password	= newLogin.Password

	if err := config.InitDB().Model(&login).Updates(models.Login{
		Email: login.Email,
		Username: login.Username,
		Password: login.Password,
	}).Error; err != nil {
		return login, err
	}

	return login, nil
}