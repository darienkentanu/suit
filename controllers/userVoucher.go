package controllers

import (
	"net/http"
	"strconv"

	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/labstack/echo/v4"
)

type UserVoucherController struct {
	userVoucherModel	database.UserVoucherModel
	userModel			database.UserModel
	voucherModel		database.VoucherModel
}

func NewUserVoucherController(userVoucherModel database.UserVoucherModel, userModel database.UserModel, voucherModel database.VoucherModel) *UserVoucherController {
	return &UserVoucherController{
		userVoucherModel: userVoucherModel,
		userModel: userModel,
		voucherModel: voucherModel,
	}
}

func (controllers *UserVoucherController) ClaimVoucher(c echo.Context) error {
	voucherID, err := strconv.Atoi(c.Param("id"))
	if err != nil || voucherID < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	voucher, err := controllers.voucherModel.GetVoucherByID(voucherID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	userID := middlewares.CurrentLoginUser(c)

	user, err := controllers.userModel.GetUserByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	if voucher.Point > user.Point {
		return echo.NewHTTPError(http.StatusBadRequest, "Not enough points")
	}

	userVoucher, err := controllers.userVoucherModel.AddUserVoucher(userID, voucherID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	newPoint := user.Point - voucher.Point

	user, err = controllers.userModel.UpdatePoint(userID, newPoint)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	resUserVoucher, err := controllers.userVoucherModel.GetUserVoucherByID(userVoucher.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusCreated, M{
		"status": "success",
		"data": resUserVoucher,
	})
}

func (controllers *UserVoucherController) RedeemVoucher(c echo.Context) error {
	voucherID, err := strconv.Atoi(c.Param("id"))
	if err != nil || voucherID < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	userID := middlewares.CurrentLoginUser(c)

	row := controllers.userVoucherModel.GetVoucherByUserAndVoucherID(userID, voucherID)
	if row == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Not available")
	}

	userVoucher, err := controllers.userVoucherModel.UpdateStatusUserVoucher(userID, voucherID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	resUserVoucher, err := controllers.userVoucherModel.GetUserVoucherByID(userVoucher.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data": resUserVoucher,
	})
}

func (controllers *UserVoucherController) GetUserVoucher(c echo.Context) error {
	userID := middlewares.CurrentLoginUser(c)

	resUserVoucher, err := controllers.userVoucherModel.GetAllVoucher(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data": resUserVoucher,
	})
}