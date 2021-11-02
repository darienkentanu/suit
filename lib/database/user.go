package database

import (
	"database/sql"

	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type UserDB struct {
	db *gorm.DB
	dbSQL *sql.DB
}

func NewUserDB(db *gorm.DB, dbSQL *sql.DB) *UserDB {
	return &UserDB{db: db, dbSQL: dbSQL}
}

func (udb *UserDB) GetPhoneNumber(number string) int {
	var user models.User
	row := udb.db.Where("phone_number = ?", number).Find(&user).RowsAffected
	return int(row)
}

func (udb *UserDB) CreateUser(user models.User) (models.User, error) {
	if err := udb.db.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (udb *UserDB) GetAllUsers() ([]models.ResponseGetUser, error) {
	var users []models.ResponseGetUser

	rows, err := udb.dbSQL.Query("SELECT u.id, u.fullname, l.email, l.username, u.phone_number, u.gender, u.address, l.role FROM users u JOIN logins l ON u.id = l.user_id")
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

func (udb *UserDB) UpdateUser(id int, newUser models.User) (models.User, error) {
	var user models.User

	if err := udb.db.First(&user, id).Error; err != nil {
		return user, err
	}

	user.Fullname = newUser.Fullname
	user.PhoneNumber = newUser.PhoneNumber
	user.Gender = newUser.Gender
	user.Address = newUser.Address
	if err := udb.db.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (udb *UserDB) GetUserProfile(id int) (interface{}, error) {
	rows, err := udb.dbSQL.Query("SELECT u.id, u.fullname, l.email, l.username, u.phone_number, u.gender, u.address, l.role FROM users u JOIN logins l ON u.id = l.user_id WHERE user_id = ?", id)
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

func (udb *UserDB) CreateCart(cart models.Cart) error {
	if err := udb.db.Save(&cart).Error; err != nil {
		return err
	}
	return nil
}