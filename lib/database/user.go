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

type UserModel interface {
	GetPhoneNumber(string) (int)
	CreateUser(user models.User) (models.User, error)
	GetAllUsers() ([]models.ResponseGetUser, error)
	UpdateUser(id int, newUser models.User) (models.User, error)
	GetUserProfile(id int) (models.ResponseGetUser, error)
	GetUserByID(userID int) (models.User, error)
	UpdatePoint(userID int, newPoint int) (models.User, error)
}

func (m *UserDB) GetPhoneNumber(number string) int {
	var user models.User
	row := m.db.Where("phone_number = ?", number).Find(&user).RowsAffected
	return int(row)
}

func (m *UserDB) CreateUser(user models.User) (models.User, error) {
	if err := m.db.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (m *UserDB) GetAllUsers() ([]models.ResponseGetUser, error) {
	var users []models.ResponseGetUser

	rows, err := m.dbSQL.Query("SELECT u.id, u.fullname, l.email, l.username, u.point, u.phone_number, u.gender, u.address, l.role FROM users u JOIN logins l ON u.id = l.user_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var user models.ResponseGetUser
		if err := rows.Scan(&user.ID, &user.Fullname, &user.Email, &user.Username, &user.Point, &user.PhoneNumber, &user.Gender, &user.Address, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (m *UserDB) UpdateUser(id int, newUser models.User) (models.User, error) {
	var user models.User

	if err := m.db.First(&user, id).Error; err != nil {
		return user, err
	}

	user.Fullname = newUser.Fullname
	user.PhoneNumber = newUser.PhoneNumber
	user.Gender = newUser.Gender
	user.Address = newUser.Address
	if err := m.db.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (m *UserDB) GetUserProfile(id int) (models.ResponseGetUser, error) {
	var account models.ResponseGetUser
	rows, err := m.dbSQL.Query("SELECT u.id, u.fullname, l.email, l.username, u.point, u.phone_number, u.gender, u.address, l.role FROM users u JOIN logins l ON u.id = l.user_id WHERE user_id = ?", id)
	if err != nil {
		return account, err
	}
	defer rows.Close()
	
	for rows.Next() {
		if err := rows.Scan(&account.ID, &account.Fullname, &account.Email, &account.Username, &account.Point, &account.PhoneNumber, &account.Gender, &account.Address, &account.Role); err != nil {
			return account, err
		}		
	}

	return account, nil
}

func (m *UserDB) UpdatePoint(userID int, newPoint int) (models.User, error) {
	var user models.User
	if err := m.db.First(&user, userID).Error; err != nil {
		return user, err
	}

	if err := m.db.Model(&user).Update("point", newPoint).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (m *UserDB) GetUserByID(userID int) (models.User, error) {
	var user models.User
	if err := m.db.First(&user, userID).Error; err != nil {
		return user, err
	}

	return user, nil
}