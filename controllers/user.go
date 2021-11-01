package controllers

import (
	"net/http"

	"github.com/darienkentanu/suit/config"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUsersController(c echo.Context) error {
	var register models.RegisterUser

	if err := c.Bind(&register); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	row := database.GetEmail(register.Email)
	if row != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is already registered")
	}

	row = database.GetPhoneNumber(register.PhoneNumber)
	if row != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Phone number is already registered")
	}

	row = database.GetUsername(register.Username)
	if row != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Username is already registered")
	}

	hashPassword, err := GenerateHashPassword(register.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error in password hash")
	}
	
	var user models.User
	user.Fullname = register.Fullname
	user.PhoneNumber = register.PhoneNumber
	user.Gender = register.Gender
	user.Address = register.Address

	user, err = database.CreateUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot create user")
	}

	var login models.Login
	login.Email = register.Email
	login.Username = register.Username
	login.Password = hashPassword
	login.Role = "user"
	login.UserID = user.ID

	login, err = database.CreateLogin(login)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot create user")
	}

	var cart models.Cart
	cart.UserID = user.ID
	err = database.CreateCart(cart)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot create cart")
	}

	return c.JSON(http.StatusCreated, M{
		"status": "user created successfully",
	})
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GetAllUsersController(c echo.Context) error {
	users, err := database.GetAllUsers()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   users,
	})
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