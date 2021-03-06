package controllers

import (
	"net/http"
	"strconv"

	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"

	"github.com/labstack/echo/v4"
)

type CategoryController struct {
	db database.CategoryModel
}

func NewCategoryController(db database.CategoryModel) *CategoryController {
	return &CategoryController{db: db}
}

func (cc *CategoryController) GetCategories(c echo.Context) error {
	categories, err := cc.db.GetCategories()

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}
	var res models.Category_Response
	var resSlc []models.Category_Response
	for _, v := range categories {
		res.ID = v.ID
		res.Name = v.Name
		res.Point = v.Point
		resSlc = append(resSlc, res)
	}
	return c.JSON(http.StatusOK, M{
		"status": "Success",
		"data":   resSlc,
	})
}

func (cc *CategoryController) AddCategories(c echo.Context) error {
	var category models.Category
	if err := c.Bind(&category); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	category, err := cc.db.AddCategories(category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusCreated, M{
		"status": "success",
		"data": M{
			"id":    category.ID,
			"name":  category.Name,
			"point": category.Point,
		},
	})
}

func (cc *CategoryController) EditCategories(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}
	var newCategory models.Category
	c.Bind(&newCategory)
	newCategory, err = cc.db.EditCategoriesById(id, newCategory)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}
	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data": M{
			"id":    newCategory.ID,
			"name":  newCategory.Name,
			"point": newCategory.Point,
		},
	})
}

func (cc *CategoryController) DeleteCategories(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	err = cc.db.DeleteCategoriesById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}

	return c.JSON(http.StatusOK, M{
		"message": "category succesfully deleted",
	})
}
