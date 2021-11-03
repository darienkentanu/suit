package controllers

import (
	"fmt"
	"net/http"

	"github.com/darienkentanu/suit/gmaps"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/darienkentanu/suit/models"
	"github.com/labstack/echo/v4"
)

type CheckoutController struct {
	checkoutModel		database.CheckoutModel
	cartModel			database.CartModel
	categoryModel		database.CategoryModel
	dropPointModel		database.DropPointsModel
	userModel			database.UserModel
	transactionModel	database.TransactionModel
}

func NewCheckoutController(checkoutModel database.CheckoutModel, cartModel database.CartModel, categoryModel database.CategoryModel, dropPointModel database.DropPointsModel, userModel database.UserModel, transactionModel database.TransactionModel) *CheckoutController {
	return &CheckoutController{
		checkoutModel: checkoutModel,
		cartModel: cartModel,
		categoryModel: categoryModel,
		dropPointModel: dropPointModel,
		userModel: userModel,
		transactionModel: transactionModel,
	}
}

func (controllers *CheckoutController) CreateCheckoutPickup(c echo.Context) error {
	var checkout models.Checkout_Input_PickUp

	if err := c.Bind(&checkout); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	userID := middlewares.CurrentLoginUser(c)

	// id category yang ingin di checkout dan terdapat pada cart
	var categoryIDSelected []int // berupa category id yg dipilih
	var categories []models.ResponseGetCategory
	for _, categoryID := range checkout.CategoryID {
		item, err := controllers.cartModel.GetItemInCart(userID, categoryID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		categoryIDSelected = append(categoryIDSelected, item.CategoryID)

		var resCategory models.ResponseGetCategory
		category, err := controllers.categoryModel.GetCategoryById(categoryID)
		resCategory.ID = category.ID
		resCategory.Name = category.Name
		resCategory.Point = category.Point
		resCategory.Weight = item.Weight

		categories = append(categories, resCategory)
	}

	if len(categoryIDSelected) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Category is not exist in cart")
	}

	checkoutID, err := controllers.checkoutModel.AddCheckout()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	for _, categoryID := range categoryIDSelected {
		_, err := controllers.cartModel.UpdateCheckoutIdInCartItem(checkoutID, userID, categoryID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}

	var totalPoint int
	for _, resCategory := range categories {
		weight := resCategory.Weight
		point := resCategory.Point

		totalPointItem := weight * point

		totalPoint += totalPointItem
	}

	user, err := controllers.userModel.GetUserByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	destination := fmt.Sprintf("%s,%s", user.Latitude, user.Longitude)

	dropPoints, err := controllers.dropPointModel.GetDropPoints()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	var dropPointDistance = make(map[int]float64)
	var dropPointDuration = make(map[int]int)
	var dropPointAddress = make(map[int]string)
	for _, dropPoint := range dropPoints {
		lat := dropPoint.Latitude
		lng := dropPoint.Longitude

		origin := fmt.Sprintf("%s,%s", lat, lng)

		distance, duration := gmaps.Distancematrix(origin, destination)
		dropPointDistance[dropPoint.ID] = distance
		dropPointDuration[dropPoint.ID] = duration
		dropPointAddress[dropPoint.ID] = dropPoint.Address
	}

	minDistance := -1.0
	var minDropPointID int
	for dropPointID, km := range dropPointDistance {
		if km < minDistance || minDistance == -1.0 {
			minDistance = km
			minDropPointID = dropPointID
		}
	}

	var transaction models.Transaction
	transaction.UserID = userID
	transaction.Point = totalPoint
	transaction.Method = "pickup"
	transaction.CheckoutID = checkoutID
	transaction.Drop_PointID = minDropPointID

	transaction, err = controllers.transactionModel.CreateTransaction(transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	var checkoutResponse models.Checkout_Response_PickUp
	checkoutResponse.TransactionID = transaction.ID
	checkoutResponse.Method = transaction.Method
	checkoutResponse.DropPointID = transaction.Drop_PointID
	checkoutResponse.DropPointAddress = dropPointAddress[minDropPointID]
	checkoutResponse.Distance = minDistance
	checkoutResponse.Duration = dropPointDuration[minDropPointID]
	checkoutResponse.TotalPoint = totalPoint
	checkoutResponse.Categories = categories

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   checkoutResponse,
	})
}