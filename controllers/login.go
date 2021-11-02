package controllers

import (
	"net/http"

	"github.com/darienkentanu/suit/gmaps"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/darienkentanu/suit/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginDB interface {
	GetEmail(string) (int)
	GetPhoneNumber(string) (int)
	GetUsername(string) (int)
	GetAccountByEmailOrUsername(requestLogin models.RequestLogin) (models.Login, error)
	UpdateToken(id int, token string) (models.Login, error)
	GetUserProfile(id int) (models.ResponseGetUser, error)
	UpdateUser(id int, user models.User) (models.User, error)
	UpdateLogin(id int, login models.Login) (models.Login, error)
}

type LoginController struct {
	db LoginDB
}

func NewLoginController(db LoginDB) LoginController {
	return LoginController{db: db}
}


func (lc *LoginController) Login(c echo.Context) error {
	var requestLogin models.RequestLogin
	
	if err := c.Bind(&requestLogin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	account, err := lc.db.GetAccountByEmailOrUsername(requestLogin)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect email or username")
	}

	check := CheckPasswordHash(requestLogin.Password, account.Password)
	if !check {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect password")
	}

	var id int
	if account.Role == "user" {
		id = account.UserID
	} else if account.Role == "staff" {
		id = account.StaffID
	}

	var newToken string
	newToken, err = middlewares.CreateToken(id, account.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot login")
	}

	account.Token = newToken
	account, err = lc.db.UpdateToken(int(account.ID), newToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot add token")
	}

	var responseLogin models.ResponseLogin
	responseLogin.Username = account.Username
	responseLogin.Email = account.Email
	responseLogin.Role = account.Role
	responseLogin.Token = account.Token

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   responseLogin,
	})
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (lc *LoginController) GetProfile(c echo.Context) error {
	// role := middlewares.CurrentRoleLoginUser(c)
	id := middlewares.CurrentLoginUser(c)
	user, err := lc.db.GetUserProfile(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   user,
	})
}

func (lc *LoginController) UpdateProfile(c echo.Context) error {
	var newProfile models.RegisterUser

	id := middlewares.CurrentLoginUser(c)

	if err := c.Bind(&newProfile); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	row := lc.db.GetEmail(newProfile.Email)
	if row > 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is already registered")
	}

	row = lc.db.GetPhoneNumber(newProfile.PhoneNumber)
	if row > 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Phone number is already registered")
	}

	row = lc.db.GetUsername(newProfile.Username)
	if row > 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Username is already registered")
	}

	hashPassword, err := GenerateHashPassword(newProfile.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error in password hash")
	}

	lat, lng := gmaps.Geocoding(newProfile.Address)

	var user models.User
	user.Fullname 		= newProfile.Fullname
	user.PhoneNumber 	= newProfile.PhoneNumber
	user.Gender			= newProfile.Gender
	user.Address 		= newProfile.Address
	user.Latitude		= lat
	user.Longitude		= lng

	user, err = lc.db.UpdateUser(id, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var login models.Login
	login.Email = newProfile.Email
	login.Username = newProfile.Username
	login.Password = hashPassword

	login, err = lc.db.UpdateLogin(id, login)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var response models.ResponseGetUser
	response.ID 			= id
	response.Fullname 		= user.Fullname
	response.Email 			= login.Email
	response.Username 		= login.Username
	response.PhoneNumber	= user.PhoneNumber
	response.Gender 		= user.Gender
	response.Address 		= user.Address
	response.Role 			= login.Role

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   response,
	})
}