package controllers

import (
	"net/http"

	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	transactionModel database.TransactionModel
	categoryModel    database.CategoryModel
	cartModel        database.CartModel
	dropPointModel   database.DropPointsModel
}

func NewTransactionController(transactionModel database.TransactionModel, categoryModel database.CategoryModel, cartModel database.CartModel, dropPointModel database.DropPointsModel) *TransactionController {
	return &TransactionController{
		transactionModel: transactionModel,
		categoryModel:    categoryModel,
		cartModel:        cartModel,
		dropPointModel:   dropPointModel,
	}
}

func (controllers *TransactionController) GetTransactions(c echo.Context) error {
	transactions, err := controllers.transactionModel.GetAllTransaction()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	var resAllTransactions []models.ResponseGetTransactions

	for _, transaction := range transactions {
		var totalPointsUsed int
		cartItems, err := controllers.cartModel.GetCartItemByCheckoutID(transaction.CheckoutID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		var categories []models.ResponseGetCategory
		for _, item := range cartItems {
			category, err := controllers.categoryModel.GetCategoryById(item.CategoryID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}

			var resCategory models.ResponseGetCategory

			resCategory.ID = category.ID
			resCategory.Name = category.Name
			resCategory.Point = category.Point
			resCategory.Weight = item.Weight
			resCategory.PointUsed = resCategory.Point * resCategory.Weight
			totalPointsUsed += resCategory.PointUsed
			categories = append(categories, resCategory)
		}

		dropPoint, err := controllers.dropPointModel.GetDropPointsByID(transaction.Drop_PointID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}

		var resTransaction models.ResponseGetTransactions
		resTransaction.ID = transaction.ID
		resTransaction.UserID = transaction.UserID
		resTransaction.Method = transaction.Method
		resTransaction.DropPointID = transaction.Drop_PointID
		resTransaction.DropPointAddress = dropPoint.Address
		// resTransaction.Point = transaction.Point
		resTransaction.Categories = categories
		if transaction.Status == 1 {
			resTransaction.Status = "transaction succeed"
		} else {
			resTransaction.Status = "transaction is being processed by staff"
		}
		resTransaction.TotalPointUsed = totalPointsUsed
		resAllTransactions = append(resAllTransactions, resTransaction)
	}

	return c.JSON(http.StatusCreated, M{
		"status": "success",
		"data":   resAllTransactions,
	})
}
