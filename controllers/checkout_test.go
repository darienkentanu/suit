package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darienkentanu/suit/constants"
	. "github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func InsertDataCheckout(db *gorm.DB) error {
	checkoutInput := models.Checkout_Input_DropOff{
		CategoryID	: []int{1},
		DropPointID	: 1,
	}

	var checkout models.Checkout
	if err := db.Save(&checkout).Error; err != nil {
		return err
	}

	var cartItem models.CartItem

	if err := db.Where("cart_user_id = ? and category_id = ? and checkout_id IS NULL", 1, 1).First(&cartItem).Error; err != nil {
		return err
	}

	if err := db.Model(&cartItem).Update("checkout_id", 1).Error; err != nil {
		return err
	}

	var transaction models.Transaction
	transaction.UserID = 1
	transaction.Point = 30
	transaction.Method = "dropoff"
	transaction.Drop_PointID = checkoutInput.DropPointID
	transaction.CheckoutID = 1

	if err := db.Save(&transaction).Error; err != nil {
		return err
	}
	
	return nil	
}

func InsertDataCheckoutVerification(db *gorm.DB) error {
	checkoutInput := models.Checkout_Input_DropOff{
		CategoryID	: []int{1},
		DropPointID	: 1,
	}

	var checkout models.Checkout
	if err := db.Save(&checkout).Error; err != nil {
		return err
	}

	var cartItem models.CartItem

	if err := db.Where("cart_user_id = ? and category_id = ? and checkout_id IS NULL", 1, 1).First(&cartItem).Error; err != nil {
		return err
	}

	if err := db.Model(&cartItem).Update("checkout_id", 1).Error; err != nil {
		return err
	}

	var transaction models.Transaction
	transaction.UserID = 1
	transaction.Point = 30
	transaction.Method = "dropoff"
	transaction.Drop_PointID = checkoutInput.DropPointID
	transaction.CheckoutID = 1
	transaction.Status = 1

	if err := db.Save(&transaction).Error; err != nil {
		return err
	}
	
	return nil	
}

func TestCreateCheckoutDropoff(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		loginPath		string
		expectCodeLogin	int
		expectCode 		int
		response   		string
		login			map[string]interface{}
		reqBody			models.Checkout_Input_DropOff
	}{
		{
			name:       "CreateCheckoutDropoff",
			path:       "/checkoutbydropoff",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusOK,
			response:   "success",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
			reqBody: 	models.Checkout_Input_DropOff{
				CategoryID: []int{1},
				DropPointID: 1,
			},
		},
	}
	
	e, db := InitEcho()
	UserSetup(db)
	Setup(db)
	CartSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	cartDB := database.NewCartDB(db)
	checkoutDB := database.NewCheckoutDB(db)
	categoryDB := database.NewCategoryDB(db)
	transactionDB := database.NewTransactionDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	checkoutController := NewCheckoutController(checkoutDB, cartDB, categoryDB, dropPointDB, userDB, transactionDB, loginDB)
	InsertDataUser(db)
	InsertDataCategory(db)
	InsertDataDropPoints(db)
	InsertDataCartItem(db)

	for _, testCase := range testCases {
		login, err := json.Marshal(testCase.login)
		if err != nil {
			t.Error(err)
		}

		loginReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(login))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		loginC := e.NewContext(loginReq, loginRec)
		
		loginC.SetPath(testCase.loginPath)

		if assert.NoError(t, loginControllers.Login(loginC)) {
			assert.Equal(t, testCase.expectCodeLogin, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string					`json:"status"`
				Data   models.ResponseLogin 	`json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			reqBody, err := json.Marshal(testCase.reqBody)
			if err != nil {
				t.Error(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(checkoutController.CreateCheckoutDropOff)(c)){
					assert.Equal(t, testCase.expectCode, rec.Code)
					body := rec.Body.String()

					var response = struct {
						Status string					`json:"status"`
						Data   models.ResponseGetUser 	`json:"data"`
					}{}
					err := json.Unmarshal([]byte(body), &response)

					if err != nil {
						assert.Error(t, err, "error")
					}
					assert.Equal(t, testCase.response, response.Status)
				}
			})
		}
	}
}

func TestCreateCheckoutDropoffError(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		loginPath		string
		expectCodeLogin	int
		expectCode 		int
		expectError   	string
		login			map[string]interface{}
		reqBody			map[string]interface{}
	}{
		{
			name:       "Create Checkout Dropoff cart is empty",
			path:       "/checkoutbydropoff",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Category is not exist in cart",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
			reqBody: 	map[string]interface{}{
				"category_id": []int{2, 3},
				"drop_point_id": 1,
			},
		},
	}
	
	e, db := InitEcho()
	UserSetup(db)
	Setup(db)
	CartSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	cartDB := database.NewCartDB(db)
	checkoutDB := database.NewCheckoutDB(db)
	categoryDB := database.NewCategoryDB(db)
	transactionDB := database.NewTransactionDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	checkoutController := NewCheckoutController(checkoutDB, cartDB, categoryDB, dropPointDB, userDB, transactionDB, loginDB)
	InsertDataUser(db)
	InsertDataCategory(db)
	InsertDataDropPoints(db)

	for _, testCase := range testCases {
		login, err := json.Marshal(testCase.login)
		if err != nil {
			t.Error(err)
		}

		loginReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(login))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		loginC := e.NewContext(loginReq, loginRec)
		
		loginC.SetPath(testCase.loginPath)

		if assert.NoError(t, loginControllers.Login(loginC)) {
			assert.Equal(t, testCase.expectCodeLogin, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string					`json:"status"`
				Data   models.ResponseLogin 	`json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			reqBody, err := json.Marshal(testCase.reqBody)
			if err != nil {
				t.Error(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(checkoutController.CreateCheckoutDropOff)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}

func TestCreateCheckoutPickup(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		loginPath		string
		expectCodeLogin	int
		expectCode 		int
		response   		string
		login			map[string]interface{}
		reqBody			models.Checkout_Input_PickUp
	}{
		{
			name:       "CreateCheckoutPickup",
			path:       "/checkoutbypickup",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusOK,
			response:   "success",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
			reqBody: 	models.Checkout_Input_PickUp{
				CategoryID: []int{1},
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	Setup(db)
	CartSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	cartDB := database.NewCartDB(db)
	checkoutDB := database.NewCheckoutDB(db)
	categoryDB := database.NewCategoryDB(db)
	transactionDB := database.NewTransactionDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	checkoutController := NewCheckoutController(checkoutDB, cartDB, categoryDB, dropPointDB, userDB, transactionDB, loginDB)
	InsertDataUser(db)
	InsertDataCategory(db)
	InsertDataDropPoints(db)
	InsertDataCartItem(db)

	for _, testCase := range testCases {
		login, err := json.Marshal(testCase.login)
		if err != nil {
			t.Error(err)
		}

		loginReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(login))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		loginC := e.NewContext(loginReq, loginRec)
		
		loginC.SetPath(testCase.loginPath)

		if assert.NoError(t, loginControllers.Login(loginC)) {
			assert.Equal(t, testCase.expectCodeLogin, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string					`json:"status"`
				Data   models.ResponseLogin 	`json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			reqBody, err := json.Marshal(testCase.reqBody)
			if err != nil {
				t.Error(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(checkoutController.CreateCheckoutPickup)(c)){
					assert.Equal(t, testCase.expectCode, rec.Code)
					body := rec.Body.String()

					var response = struct {
						Status string					`json:"status"`
						Data   models.ResponseGetUser 	`json:"data"`
					}{}
					err := json.Unmarshal([]byte(body), &response)

					if err != nil {
						assert.Error(t, err, "error")
					}
					assert.Equal(t, testCase.response, response.Status)
				}
			})
		}
	}
}

func TestCreateCheckoutPickupError(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		loginPath		string
		expectCodeLogin	int
		expectCode 		int
		expectError   	string
		login			map[string]interface{}
		reqBody			map[string]interface{}
	}{
		{
			name:       "Create Checkout Pickup Invalid Input",
			path:       "/checkoutbypickup",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Invalid input",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
			reqBody: 	map[string]interface{}{
				"category_id": "a",
			},
		},
		{
			name:       "Create Checkout Pickup cart is empty",
			path:       "/checkoutbypickup",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Category is not exist in cart",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
			reqBody: 	map[string]interface{}{
				"category_id": []int{2, 3},
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	Setup(db)
	CartSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	cartDB := database.NewCartDB(db)
	checkoutDB := database.NewCheckoutDB(db)
	categoryDB := database.NewCategoryDB(db)
	transactionDB := database.NewTransactionDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	checkoutController := NewCheckoutController(checkoutDB, cartDB, categoryDB, dropPointDB, userDB, transactionDB, loginDB)
	InsertDataUser(db)
	InsertDataCategory(db)
	InsertDataDropPoints(db)

	for _, testCase := range testCases {
		login, err := json.Marshal(testCase.login)
		if err != nil {
			t.Error(err)
		}

		loginReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(login))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		loginC := e.NewContext(loginReq, loginRec)
		
		loginC.SetPath(testCase.loginPath)

		if assert.NoError(t, loginControllers.Login(loginC)) {
			assert.Equal(t, testCase.expectCodeLogin, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string					`json:"status"`
				Data   models.ResponseLogin 	`json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			reqBody, err := json.Marshal(testCase.reqBody)
			if err != nil {
				t.Error(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(checkoutController.CreateCheckoutPickup)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}

func TestVerification(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		loginPath		string
		expectCodeLogin	int
		expectCode 		int
		response   		string
		login			map[string]interface{}
	}{
		{
			name:       "Verification",
			path:       "/verification/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusOK,
			response:   "success",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	Setup(db)
	CartSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	cartDB := database.NewCartDB(db)
	checkoutDB := database.NewCheckoutDB(db)
	categoryDB := database.NewCategoryDB(db)
	transactionDB := database.NewTransactionDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	checkoutController := NewCheckoutController(checkoutDB, cartDB, categoryDB, dropPointDB, userDB, transactionDB, loginDB)
	InsertDataUser(db)
	InsertDataCategory(db)
	InsertDataDropPoints(db)
	InsertDataCartItem(db)
	InsertDataCheckout(db)

	for _, testCase := range testCases {
		login, err := json.Marshal(testCase.login)
		if err != nil {
			t.Error(err)
		}

		loginReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(login))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		loginC := e.NewContext(loginReq, loginRec)
		
		loginC.SetPath(testCase.loginPath)

		if assert.NoError(t, loginControllers.Login(loginC)) {
			assert.Equal(t, testCase.expectCodeLogin, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string					`json:"status"`
				Data   models.ResponseLogin 	`json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			req := httptest.NewRequest(http.MethodPut, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues("1")

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(checkoutController.Verification)(c)){
					assert.Equal(t, testCase.expectCode, rec.Code)
					body := rec.Body.String()

					var response = struct {
						Status string					`json:"status"`
						Data   models.ResponseGetUser 	`json:"data"`
					}{}
					err := json.Unmarshal([]byte(body), &response)

					if err != nil {
						assert.Error(t, err, "error")
					}
					assert.Equal(t, testCase.response, response.Status)
				}
			})
		}
	}
}

func TestVerificationError(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		loginPath		string
		expectCodeLogin	int
		expectCode 		int
		expectError   	string
		paramValues		string
		login			map[string]interface{}
	}{
		{
			name:       "Verification Invalid ID",
			path:       "/verification/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Invalid id",
			paramValues: "a",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
		{
			name:       "Verification Invalid Trans ID",
			path:       "/verification/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusInternalServerError,
			expectError:   "Internal server error",
			paramValues: "50",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	Setup(db)
	CartSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	cartDB := database.NewCartDB(db)
	checkoutDB := database.NewCheckoutDB(db)
	categoryDB := database.NewCategoryDB(db)
	transactionDB := database.NewTransactionDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	checkoutController := NewCheckoutController(checkoutDB, cartDB, categoryDB, dropPointDB, userDB, transactionDB, loginDB)
	InsertDataUser(db)
	InsertDataCategory(db)
	InsertDataDropPoints(db)
	InsertDataCartItem(db)
	InsertDataCheckout(db)

	for _, testCase := range testCases {
		login, err := json.Marshal(testCase.login)
		if err != nil {
			t.Error(err)
		}

		loginReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(login))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		loginC := e.NewContext(loginReq, loginRec)
		
		loginC.SetPath(testCase.loginPath)

		if assert.NoError(t, loginControllers.Login(loginC)) {
			assert.Equal(t, testCase.expectCodeLogin, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string					`json:"status"`
				Data   models.ResponseLogin 	`json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			req := httptest.NewRequest(http.MethodPut, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues(testCase.paramValues)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(checkoutController.Verification)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}