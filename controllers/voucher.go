package controllers

import (
	"net/http"
	"strconv"

	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"

	"github.com/labstack/echo/v4"
)

type VoucherController struct {
	db database.VoucherModel
}

func NewVoucherController(db database.VoucherModel) *VoucherController {
	return &VoucherController{db: db}
}

func (vc *VoucherController) GetVouchers(c echo.Context) error {
	vouchers, err := vc.db.GetVouchers()

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}
	var response models.Voucher_Response
	var responseSlice []models.Voucher_Response
	for _, value := range vouchers {
		response.ID = value.ID
		response.Name = value.Name
		response.Point = value.Point
		responseSlice = append(responseSlice, response)
	}
	return c.JSON(http.StatusOK, M{
		"status": "Success",
		"data":   responseSlice,
	})
}

func (vc *VoucherController) AddVouchers(c echo.Context) error {
	var vouchers models.Voucher
	if err := c.Bind(&vouchers); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	vouchers, err := vc.db.AddVouchers(vouchers)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusCreated, M{
		"status": "success",
		"data": M{
			"id":    vouchers.ID,
			"name":  vouchers.Name,
			"point": vouchers.Point,
		},
	})
}

func (vc *VoucherController) EditVouchers(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}
	var newVoucher models.Voucher
	if err := c.Bind(&newVoucher); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	newVoucher, err = vc.db.EditVouchersById(id, newVoucher)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}
	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data": M{
			"id":    newVoucher.ID,
			"name":  newVoucher.Name,
			"point": newVoucher.Point,
		},
	})
}

func (vc *VoucherController) DeleteVouchers(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	err = vc.db.DeleteVouchersById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	return c.JSON(http.StatusOK, M{
		"message": "voucher succesfully deleted",
	})
}
