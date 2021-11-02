package controllers

import (
	"net/http"

	"github.com/darienkentanu/suit/gmaps"
	"github.com/darienkentanu/suit/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserDB interface {
	GetEmail(string) (int)
	GetPhoneNumber(string) (int)
	GetUsername(string) (int)
	CreateUser(user models.User) (models.User, error)
	GetAllUsers() ([]models.ResponseGetUser, error)
	CreateCart(cart models.Cart) (error)
}

type UserController struct {
	db UserDB
	ldb LoginDB
}

func NewUserController(db UserDB) UserController {
	return UserController{db: db}
}

func (uc *UserController) RegisterUsers(c echo.Context) error {
	var register models.RegisterUser

	if err := c.Bind(&register); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	row := uc.db.GetEmail(register.Email)
	if row != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is already registered")
	}

	row = uc.db.GetPhoneNumber(register.PhoneNumber)
	if row != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Phone number is already registered")
	}

	row = uc.db.GetUsername(register.Username)
	if row != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Username is already registered")
	}

	hashPassword, err := GenerateHashPassword(register.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error in password hash")
	}

	lat, lng := gmaps.Geocoding(register.Address)
	
	var user models.User
	user.Fullname = register.Fullname
	user.PhoneNumber = register.PhoneNumber
	user.Gender = register.Gender
	user.Address = register.Address
	user.Latitude = lat
	user.Longitude = lng
	
	user, err = uc.db.CreateUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot create user")
	}

	var login models.Login
	login.Email = register.Email
	login.Username = register.Username
	login.Password = hashPassword
	login.Role = "user"
	login.UserID = user.ID

	login, err = uc.db.CreateLogin(login)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot create user")
	}

	var cart models.Cart
	cart.UserID = user.ID
	err = uc.db.CreateCart(cart)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot create cart")
	}

	var response models.ResponseGetUser
	response.Fullname 	= user.Fullname
	response.Email 		= login.Email
	response.Username 	= login.Username
	response.PhoneNumber= user.PhoneNumber
	response.Gender		= user.Gender
	response.Address 	= user.Address
	response.Role		= login.Role

	return c.JSON(http.StatusCreated, M{
		"status": "sucess",
		"data": response,
	})
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (uc *UserController) GetAllUsers(c echo.Context) error {
	users, err := uc.db.GetAllUsers()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   users,
	})
}

