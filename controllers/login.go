package controllers

import (
	"net/http"

	"github.com/darienkentanu/suit/gmaps"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/darienkentanu/suit/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	userModel  		database.UserModel
	loginModel 		database.LoginModel
	staffModel 		database.StaffModel
	dropPointsModel	database.DropPointsModel
}

func NewLoginController(userModel database.UserModel, loginModel database.LoginModel, staffModel database.StaffModel, dropPointsModel database.DropPointsModel) *LoginController {
	return &LoginController{
		userModel:  userModel,
		loginModel: loginModel,
		staffModel: staffModel,
		dropPointsModel: dropPointsModel,
	}
}

func (controllers *LoginController) Login(c echo.Context) error {
	var requestLogin models.RequestLogin

	if err := c.Bind(&requestLogin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	account, err := controllers.loginModel.GetAccountByEmailOrUsername(requestLogin)
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

	loginID := account.ID
	role := account.Role

	var newToken string
	newToken, err = middlewares.CreateToken(id, loginID, role)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot login")
	}

	account.Token = newToken
	account, err = controllers.loginModel.UpdateToken(int(account.ID), newToken)
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

func (controllers *LoginController) GetProfile(c echo.Context) error {
	role := middlewares.CurrentRoleLoginUser(c)
	id := middlewares.CurrentLoginUser(c)

	if role == "user" {
		user, err := controllers.userModel.GetUserProfile(id)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		return c.JSON(http.StatusOK, M{
			"status": "success",
			"data":   user,
		})
	} else if role == "staff" {
		staff, err := controllers.staffModel.GetStaffByID(id)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		return c.JSON(http.StatusOK, M{
			"status": "success",
			"data":   staff,
		})
	}
	
	return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
}

func (controllers *LoginController) UpdateProfile(c echo.Context) error {
	role := middlewares.CurrentRoleLoginUser(c)
	id := middlewares.CurrentLoginUser(c)
	loginID := middlewares.CurrentLoginID(c)

	if role == "user" {
		var newUserProfile models.RegisterUser

		if err := c.Bind(&newUserProfile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
		}

		row := controllers.loginModel.GetEmail(newUserProfile.Email)
		if row > 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "Email is already registered")
		}

		row = controllers.userModel.GetPhoneNumber(newUserProfile.PhoneNumber)
		if row > 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "Phone number is already registered")
		}

		row = controllers.loginModel.GetUsername(newUserProfile.Username)
		if row > 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "Username is already registered")
		}

		hashPassword, err := GenerateHashPassword(newUserProfile.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error in password hash")
		}

		lat, lng := gmaps.Geocoding(newUserProfile.Address)

		var user models.User
		user.Fullname = newUserProfile.Fullname
		user.PhoneNumber = newUserProfile.PhoneNumber
		user.Gender = newUserProfile.Gender
		user.Address = newUserProfile.Address
		user.Latitude = lat
		user.Longitude = lng

		user, err = controllers.userModel.UpdateUser(id, user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		var login models.Login
		login.Email = newUserProfile.Email
		login.Username = newUserProfile.Username
		login.Password = hashPassword

		login, err = controllers.loginModel.UpdateLogin(loginID, login)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		var responseUser models.ResponseGetUser
		responseUser.ID = id
		responseUser.Fullname = user.Fullname
		responseUser.Email = login.Email
		responseUser.Username = login.Username
		responseUser.PhoneNumber = user.PhoneNumber
		responseUser.Gender = user.Gender
		responseUser.Address = user.Address
		responseUser.Role = login.Role

		return c.JSON(http.StatusOK, M{
			"status": "success",
			"data":   responseUser,
		})
	} else if role == "staff" {
		var newStaffProfile models.RegisterStaff

		if err := c.Bind(&newStaffProfile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
		}

		row := controllers.loginModel.GetEmail(newStaffProfile.Email)
		if row > 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "Email is already registered")
		}

		row = controllers.loginModel.GetUsername(newStaffProfile.Username)
		if row > 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "Username is already registered")
		}

		row = controllers.staffModel.GetPhoneNumberStaff(newStaffProfile.PhoneNumber)
		if row > 1 {
			return echo.NewHTTPError(http.StatusBadRequest, "Phone number is already registered")
		}

		hashPassword, err := GenerateHashPassword(newStaffProfile.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error in password hash")
		}

		var staff models.Staff
		staff.Fullname = newStaffProfile.Fullname
		staff.PhoneNumber = newStaffProfile.PhoneNumber
		staff.Drop_PointID = newStaffProfile.DropPointID

		staff, err = controllers.staffModel.UpdateStaff(id, staff)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		var login models.Login
		login.Email = newStaffProfile.Email
		login.Username = newStaffProfile.Username
		login.Password = hashPassword

		login, err = controllers.loginModel.UpdateLogin(loginID, login)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		dropPoint, err := controllers.dropPointsModel.GetDropPointsByID(staff.Drop_PointID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		var responseStaff models.ResponseGetStaff
		responseStaff.ID = id
		responseStaff.Fullname = staff.Fullname
		responseStaff.Email = login.Email
		responseStaff.Username = login.Username
		responseStaff.PhoneNumber = staff.PhoneNumber
		responseStaff.Role = login.Role
		responseStaff.DropPointID = staff.Drop_PointID
		responseStaff.DropPointAddress = dropPoint.Address

		return c.JSON(http.StatusOK, M{
			"status": "success",
			"data":   responseStaff,
		})
	}
	
	return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
}
