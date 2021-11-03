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
	c.Bind(&input)
	userID := middlewares.CurrentLoginUser(c)
	newItem, err := cc.db.AddToCart(userID, input)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err)
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
		echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   cartItems,
	})
}

func (cc *CartController) EditCartItem(c echo.Context) error {
	var input models.CartItem_Input
	c.Bind(&input)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cart item id")
	}
	newCartItem, err := cc.db.EditCartItem(id, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, M{
		"message": "cart item succesfully deleted",
	})
}
