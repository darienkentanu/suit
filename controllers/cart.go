package controllers

import (
	"net/http"
	"strconv"

	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/darienkentanu/suit/models"
	"github.com/labstack/echo/v4"
)

type CartController struct {
	db database.CartModel
}

func NewCartController(db database.CartModel) *CartController {
	return &CartController{db: db}
}

func (cc *CartController) AddToCart(c echo.Context) error {
	var input models.CartItem_Input
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	userID := middlewares.CurrentLoginUser(c)
	exist := cc.db.CheckCartByCategoryID(userID, input.CategoryID)
	if !exist {
		newItem, err := cc.db.AddToCart(userID, input)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
		return c.JSON(http.StatusCreated, M{
			"status": "success",
			"data":   newItem,
		})
	}
	
	newItem, err := cc.db.AddCartWeight(userID, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.JSON(http.StatusCreated, M{
		"status": "success",
		"data":   newItem,
	})
}

func (cc *CartController) GetCartItem(c echo.Context) error {
	userID := middlewares.CurrentLoginUser(c)
	cartItems, err := cc.db.GetCartItem(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   cartItems,
	})
}

func (cc *CartController) EditCartItem(c echo.Context) error {
	var input models.CartItem_Input
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cart item id")
	}
	newCartItem, err := cc.db.EditCartItem(id, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}
	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   newCartItem,
	})
}

func (cc *CartController) DeleteCartItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cart item id")
	}
	err = cc.db.DeleteCartItem(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	}
	return c.JSON(http.StatusOK, M{
		"message": "cart item succesfully deleted",
	})
}
