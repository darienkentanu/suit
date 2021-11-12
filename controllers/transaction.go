package controllers

import (
	"net/http"
	"strconv"

	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/middlewares"
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
	role := middlewares.CurrentRoleLoginUser(c)
	id := middlewares.CurrentLoginUser(c)

	var transactions []models.Transaction
	var err error
	if role == "staff" {
		transactions, err = controllers.transactionModel.GetAllTransaction()
	} else if role == "user" {
		transactions, err = controllers.transactionModel.GetTransactionsByUserID(id)
	}

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
			resCategory.ReceivedPoints = resCategory.Point * resCategory.Weight
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
		resTransaction.TotalReceivedPoints = transaction.Point
		resTransaction.Categories = categories
		if transaction.Status == 1 {
			resTransaction.Status = "transaction succeed"
		} else {
			resTransaction.Status = "transaction is being processed by staff"
		}
		resTransaction.CreatedAt = transaction.CreatedAt
		resTransaction.UpdatedAt = transaction.UpdatedAt
		// resTransaction.TotalPoint = totalPointsUsed
		resAllTransactions = append(resAllTransactions, resTransaction)
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   resAllTransactions,
	})
}

func (controllers *TransactionController) GetTransactionsDropPoint(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id")
	}

	transactions, err := controllers.transactionModel.GetTransactionsByDropPointID(id)
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
			resCategory.ReceivedPoints = resCategory.Point * resCategory.Weight
			categories = append(categories, resCategory)
		}

		var resTransaction models.ResponseGetTransactions
		resTransaction.ID = transaction.ID
		resTransaction.UserID = transaction.UserID
		resTransaction.Method = transaction.Method
		resTransaction.DropPointID = transaction.Drop_PointID
		resTransaction.DropPointAddress = transaction.DropPointAddress
		resTransaction.TotalReceivedPoints = transaction.Point
		resTransaction.Categories = categories
		if transaction.Status == 1 {
			resTransaction.Status = "transaction succeed"
		} else {
			resTransaction.Status = "transaction is being processed by staff"
		}
		// resTransaction.TotalPoint = totalPointsUsed
		resAllTransactions = append(resAllTransactions, resTransaction)
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   resAllTransactions,
	})
}

func (controllers *TransactionController) GetTransactionsWithRangeDate(c echo.Context) error {
	rangeDate := c.Param("range")

	if rangeDate != "daily" && rangeDate != "weekly" && rangeDate != "monthly" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid range")
	}

	role := middlewares.CurrentRoleLoginUser(c)
	id := middlewares.CurrentLoginUser(c)

	transactions, err := controllers.transactionModel.GetTransationsRangeDate(id, role, rangeDate)
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
			resCategory.ReceivedPoints = resCategory.Point * resCategory.Weight
			categories = append(categories, resCategory)
		}

		var resTransaction models.ResponseGetTransactions
		resTransaction.ID = transaction.ID
		resTransaction.UserID = transaction.UserID
		resTransaction.Method = transaction.Method
		resTransaction.DropPointID = transaction.Drop_PointID
		resTransaction.DropPointAddress = transaction.DropPointAddress
		resTransaction.TotalReceivedPoints = transaction.Point
		resTransaction.Categories = categories
		if transaction.Status == 1 {
			resTransaction.Status = "transaction succeed"
		} else {
			resTransaction.Status = "transaction is being processed by staff"
		}
		resAllTransactions = append(resAllTransactions, resTransaction)
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   resAllTransactions,
	})
}

func (controllers *TransactionController) GetTransactionTotalWithRangeDate(c echo.Context) error {
	rangeDate := c.Param("range")

	if rangeDate != "daily" && rangeDate != "weekly" && rangeDate != "monthly" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid range")
	}

	role := middlewares.CurrentRoleLoginUser(c)
	id := middlewares.CurrentLoginUser(c)

	response, err := controllers.transactionModel.GetTransationTotalRangeDate(id, role, rangeDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   response,
	})
}

func (controllers *TransactionController) GetTransactionTotal(c echo.Context) error {
	role := middlewares.CurrentRoleLoginUser(c)
	id := middlewares.CurrentLoginUser(c)

	response, err := controllers.transactionModel.GetTransationTotal(id, role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, M{
		"status": "success",
		"data":   response,
	})
}
