package controllers

import (
	"net/http"
	"strconv"

	"github.com/darienkentanu/suit/models"

	"github.com/labstack/echo/v4"
)

type CategoryDB interface {
	GetCategories() ([]models.Category, error)
	AddCategories(categories models.Category) (models.Category, error)
	EditCategoriesById(id int, newCategories models.Category) (models.Category, error)
	DeleteCategoriesById(id int) error
}

type CategoryController struct {
	categoryDB CategoryDB
}

func NewCategoryController(categoryDB CategoryDB) CategoryController {
	return CategoryController{categoryDB: categoryDB}
}

func (cc *CategoryController) GetCategories(c echo.Context) error {
	categories, err := cc.categoryDB.GetCategories()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	c.Bind(&category)

	category, err := cc.categoryDB.AddCategories(category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var newCategory models.Category
	c.Bind(&newCategory)
	newCategory, err = cc.categoryDB.EditCategoriesById(id, newCategory)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, M{
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = cc.categoryDB.DeleteCategoriesById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, M{
		"message": "category succesfully deleted",
	})
}
