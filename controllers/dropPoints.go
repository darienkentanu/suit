package controllers

import (
	"net/http"
	"strconv"

	"github.com/darienkentanu/suit/gmaps"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"

	"github.com/labstack/echo/v4"
)

type DropPointsController struct {
	db database.DropPointsModel
}

func NewDropPointsController(db database.DropPointsModel) *DropPointsController {
	return &DropPointsController{db: db}
}

func (dpc *DropPointsController) GetDropPoints(c echo.Context) error {
	dropPoints, err := dpc.db.GetDropPoints()

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}
	var res models.Drop_Points_Response
	var resSlc []models.Drop_Points_Response
	for _, v := range dropPoints {
		res.ID = v.ID
		res.Address = v.Address
		resSlc = append(resSlc, res)
	}
	return c.JSON(http.StatusOK, M{
		"status": "Success",
		"data":   resSlc,
	})
}

func (dpc *DropPointsController) AddDropPoints(c echo.Context) error {
	var dropPoints models.Drop_Point
	if err := c.Bind(&dropPoints); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	dropPoints.Latitude, dropPoints.Longitude = gmaps.Geocoding(dropPoints.Address)

	dropPoints, err := dpc.db.AddDropPoints(dropPoints)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusCreated, M{
		"status": "success",
		"data": M{
			"id":      dropPoints.ID,
			"address": dropPoints.Address,
		},
	})
}

func (dpc *DropPointsController) EditDropPoints(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}
	var newDropPoints models.Drop_Point
	if err := c.Bind(&newDropPoints); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	newDropPoints.Latitude, newDropPoints.Longitude = gmaps.Geocoding(newDropPoints.Address)
	newDropPoints, err = dpc.db.EditDropPointsById(id, newDropPoints)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}
	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data": M{
			"id":      newDropPoints.ID,
			"address": newDropPoints.Address,
		},
	})
}

func (dpc *DropPointsController) DeleteDropPoints(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	err = dpc.db.DeleteDropPointsById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	return c.JSON(http.StatusOK, M{
		"message": "drop point succesfully deleted",
	})
}
