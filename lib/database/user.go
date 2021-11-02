package database

import (
	"github.com/darienkentanu/suit/config"
	"github.com/darienkentanu/suit/models"
)

func GetPhoneNumber(number string) int {
	var user models.User
	row := config.InitDB().Where("phone_number = ?", number).Find(&user).RowsAffected
	return int(row)
}

func CreateUser(user models.User) (models.User, error) {
	if err := config.InitDB().Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetAllUsers() ([]models.ResponseGetUser, error) {
	var users []models.ResponseGetUser

	rows, err := config.InitDBSQL().Query("SELECT u.id, u.fullname, l.email, l.username, u.phone_number, u.gender, u.address, l.role FROM users u JOIN logins l ON u.id = l.user_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var user models.ResponseGetUser
		if err := rows.Scan(&user.ID, &user.Fullname, &user.Email, &user.Username, &user.PhoneNumber, &user.Gender, &user.Address, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func UpdateUser(id int, newUser models.User) (models.User, error) {
	var user models.User

	if err := config.InitDB().First(&user, id).Error; err != nil {
		return user, err
	}

	user.Fullname = newUser.Fullname
	user.PhoneNumber = newUser.PhoneNumber
	user.Gender = newUser.Gender
	user.Address = newUser.Address
	if err := config.InitDB().Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func GetUserProfile(id int, role string) (interface{}, error) {
	rows, err := config.InitDBSQL().Query("SELECT u.id, u.fullname, l.email, l.username, u.phone_number, u.gender, u.address, l.role FROM users u JOIN logins l ON u.id = l.user_id WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var account models.ResponseGetUser
		if err := rows.Scan(&account.ID, &account.Fullname, &account.Email, &account.Username, &account.PhoneNumber, &account.Gender, &account.Address, &account.Role); err != nil {
			return nil, err
		}
		
		return account, nil
	}

	return nil, nil
}