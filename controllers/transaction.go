package controllers

import (
	"net/http"

	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	transactionModel	database.TransactionModel
	categoryModel		database.CategoryModel
	cartModel			database.CartModel
}

func NewTransactionController(transactionModel database.TransactionModel, categoryModel database.CategoryModel, cartModel database.CartModel) *TransactionController {
	return &TransactionController{
		transactionModel: transactionModel,
		categoryModel: categoryModel,
		cartModel: cartModel,
	}
}

func (controllers *TransactionController) GetTransactions(c echo.Context) error {
	transactions, err := controllers.transactionModel.GetAllTransaction()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	var resAllTransactions []models.ResponseGetTransactions

	for _, transaction := range transactions {
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

			categories = append(categories, resCategory)
		}

		var resTransaction models.ResponseGetTransactions
		resTransaction.ID = transaction.ID
		resTransaction.UserID = transaction.UserID
		resTransaction.Method = transaction.Method
		resTransaction.DropPointID = transaction.Drop_PointID
		resTransaction.Point = transaction.Point
		resTransaction.Categories = categories

		resAllTransactions = append(resAllTransactions, resTransaction)
	}

	return c.JSON(http.StatusCreated, M{
		"status": "success",
		"data": resAllTransactions,
	})
}