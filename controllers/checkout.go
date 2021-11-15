package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/darienkentanu/suit/gmaps"
	"github.com/darienkentanu/suit/gomails"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
	"github.com/darienkentanu/suit/models"
	"github.com/labstack/echo/v4"
)

type CheckoutController struct {
	checkoutModel    database.CheckoutModel
	cartModel        database.CartModel
	categoryModel    database.CategoryModel
	dropPointModel   database.DropPointsModel
	userModel        database.UserModel
	transactionModel database.TransactionModel
	loginModel       database.LoginModel
}

func NewCheckoutController(checkoutModel database.CheckoutModel, cartModel database.CartModel, categoryModel database.CategoryModel, dropPointModel database.DropPointsModel, userModel database.UserModel, transactionModel database.TransactionModel, loginModel database.LoginModel) *CheckoutController {
	return &CheckoutController{
		checkoutModel:    checkoutModel,
		cartModel:        cartModel,
		categoryModel:    categoryModel,
		dropPointModel:   dropPointModel,
		userModel:        userModel,
		transactionModel: transactionModel,
		loginModel:       loginModel,
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
		exist := controllers.cartModel.CheckCartByCategoryID(userID, categoryID)
		if !exist {
			continue
		}

		categoryIDSelected = append(categoryIDSelected, categoryID)

		var resCategory models.ResponseGetCategory
		category, err := controllers.categoryModel.GetCategoryById(categoryID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		item, err := controllers.cartModel.GetItemInCart(userID, categoryID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
		resCategory.ID = category.ID
		resCategory.Name = category.Name
		resCategory.Point = category.Point
		resCategory.Weight = item.Weight
		resCategory.ReceivedPoints = category.Point * item.Weight

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
	checkoutResponse.TotalReceivedPoints = totalPoint
	checkoutResponse.Categories = categories

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   checkoutResponse,
	})
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (controllers *CheckoutController) CreateCheckoutDropOff(c echo.Context) error {
	var checkout models.Checkout_Input_DropOff
	c.Bind(&checkout)

	userID := middlewares.CurrentLoginUser(c)
	// id category yang ingin di checkout dan terdapat pada cart
	var categoryIDSelected []int
	var categories []models.ResponseGetCategory
	for _, categoryID := range checkout.CategoryID {
		if exist := controllers.cartModel.CheckCartByCategoryID(userID, categoryID); !exist {
			continue
		}
		categoryIDSelected = append(categoryIDSelected, categoryID)

		var resCategory models.ResponseGetCategory
		category, err := controllers.categoryModel.GetCategoryById(categoryID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		item, err := controllers.cartModel.GetItemInCart(userID, categoryID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
		resCategory.ID = category.ID
		resCategory.Name = category.Name
		resCategory.Point = category.Point
		resCategory.Weight = item.Weight
		resCategory.ReceivedPoints = category.Point * item.Weight

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

	var transaction models.Transaction
	transaction.UserID = userID
	transaction.Point = totalPoint
	transaction.Method = "dropoff"
	transaction.CheckoutID = checkoutID
	transaction.Drop_PointID = checkout.DropPointID

	transaction, err = controllers.transactionModel.CreateTransaction(transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	dropPoint, err := controllers.dropPointModel.GetDropPointsByID(checkout.DropPointID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	var checkoutResponse models.Checkout_Response_DropOff
	checkoutResponse.TransactionID = transaction.ID
	checkoutResponse.Method = transaction.Method
	checkoutResponse.DropPointID = transaction.Drop_PointID
	checkoutResponse.DropPointAddress = dropPoint.Address
	checkoutResponse.TotalReceivedPoints = totalPoint
	checkoutResponse.Categories = categories

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   checkoutResponse,
	})
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (controllers *CheckoutController) Verification(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}
	transaction, err := controllers.transactionModel.UpdateStatusTransaction(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	user, err := controllers.userModel.GetUserByID(transaction.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	// var newPoint int
	newPoint := user.Point + transaction.Point
	_, err = controllers.userModel.UpdatePoint(transaction.UserID, newPoint)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	userLogin, err := controllers.loginModel.GetLoginByUserID(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	message := fmt.Sprintf("transaksi kamu untuk drop point %d telah di verifikasi! Kamu mendapatkan %d point. Berikut detail transaksinya:<br/>ID=%d<br>Method=%s<br/>Tanggal Transaksi=%v<br/>Tanggal Verifikasi=%v", 
				transaction.Drop_PointID, 
				transaction.Point, 			
				transaction.ID,
				transaction.Method,  
				transaction.CreatedAt, 
				transaction.UpdatedAt)

	err = gomails.SendMail(userLogin.Email, message)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   transaction,
	})
}
