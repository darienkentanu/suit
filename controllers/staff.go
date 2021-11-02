package controllers

import (
	"net/http"

	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"
	"github.com/labstack/echo/v4"
)

type StaffController struct {
	staffModel	database.StaffModel
	loginModel	database.LoginModel
}

func NewStaffController(staffModel database.StaffModel, loginModel database.LoginModel) *StaffController {
	return &StaffController{
		staffModel: staffModel,
		loginModel: loginModel,
	}
}

func (controllers *StaffController) AddStaff(c echo.Context) error {
	var register models.RegisterStaff

	if err := c.Bind(&register); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	row := controllers.loginModel.GetEmail(register.Email)
	if row != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Email is already registered")
	}

	row = controllers.staffModel.GetPhoneNumberStaff(register.PhoneNumber)
	if row != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Phone number is already registered")
	}

	row = controllers.loginModel.GetUsername(register.Username)
	if row != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Username is already registered")
	}

	hashPassword, err := GenerateHashPassword(register.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error in password hash")
	}
	
	var staff models.Staff
	staff.Fullname = register.Fullname
	staff.PhoneNumber = register.PhoneNumber
	staff.Drop_PointID = register.DropPointID
	
	staff, err = controllers.staffModel.CreateStaff(staff)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot create staff")
	}

	var login models.Login
	login.Email = register.Email
	login.Username = register.Username
	login.Password = hashPassword
	login.Role = "staff"
	login.StaffID = staff.ID

	login, err = controllers.loginModel.CreateLoginStaff(login)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot create staff")
	}

	var response models.ResponseGetStaff
	response.ID			= staff.ID
	response.Fullname 	= staff.Fullname
	response.Email 		= login.Email
	response.Username 	= login.Username
	response.PhoneNumber= staff.PhoneNumber
	response.Role		= login.Role
	response.DropPointID= staff.Drop_PointID

	return c.JSON(http.StatusCreated, M{
		"status": "success",
		"data": response,
	})
}

func (controllers *StaffController) GetAllStaff(c echo.Context) error {
	staff, err := controllers.staffModel.GetAllStaff()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   staff,
	})
}