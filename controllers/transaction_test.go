package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darienkentanu/suit/constants"
	// . "github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TransactionSetup(db *gorm.DB) {
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.CartItem{})
	db.Migrator().DropTable(&models.Category{})
	db.Migrator().DropTable(&models.Checkout{})
	db.Migrator().DropTable(&models.Login{})
	db.Migrator().DropTable(&models.Staff{})
	db.Migrator().DropTable(&models.Drop_Point{})
	db.Migrator().DropTable(&models.Cart{})
	db.Migrator().DropTable(&models.User{})

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Cart{})
	db.AutoMigrate(&models.Drop_Point{})
	db.AutoMigrate(&models.Staff{})
	db.AutoMigrate(&models.Login{})
	db.AutoMigrate(&models.Checkout{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.CartItem{})
	db.AutoMigrate(&models.Transaction{})
}

func TestGetTransaction(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		loginPath  string
		expectCode int
		response   string
		login      map[string]interface{}
	}{
		{
			name:       "Get Transaction User",
			path:       "/transactions",
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
			},
		},
		{
			name: "Get All Transaction (staff)",
			path: "/transactions",
			// login staff
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "azkam@gmail.com",
				"password": "azka123",
			},
		},
	}

	e, db := InitEcho()
	// UserSetup(db)
	// CartSetup(db)
	TransactionSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	transDB := database.NewTransactionDB(db)
	cartDB := database.NewCartDB(db)
	categoryDB := database.NewCategoryDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	transControllers := NewTransactionController(transDB, categoryDB, cartDB, dropPointDB)

	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)
	InsertDataCategory(db)
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
			assert.Equal(t, testCase.expectCode, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(transControllers.GetTransactions)(c)) {
					assert.Equal(t, testCase.expectCode, rec.Code)
					body := rec.Body.String()

					var response = struct {
						Status string                 `json:"status"`
						Data   models.ResponseGetUser `json:"data"`
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

func TestGetTransactionByDropPoints(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		loginPath  string
		expectCode int
		response   string
		login      map[string]interface{}
	}{
		{
			name: "Get Transaction by drop points",
			path: "/transactionsbydroppoint/:id",
			// login staff
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "azkam@gmail.com",
				"password": "azka123",
			},
		},
	}

	e, db := InitEcho()
	// UserSetup(db)
	// CartSetup(db)
	TransactionSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	transDB := database.NewTransactionDB(db)
	cartDB := database.NewCartDB(db)
	categoryDB := database.NewCategoryDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	transControllers := NewTransactionController(transDB, categoryDB, cartDB, dropPointDB)

	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)
	InsertDataCategory(db)
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
			assert.Equal(t, testCase.expectCode, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues("1")

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(transControllers.GetTransactionsDropPoint)(c)) {
					assert.Equal(t, testCase.expectCode, rec.Code)
					body := rec.Body.String()

					var response = struct {
						Status string                         `json:"status"`
						Data   models.ResponseGetTransactions `json:"data"`
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

func TestGetTransactionWithRange(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		loginPath  string
		expectCode int
		response   string
		login      map[string]interface{}
	}{
		{
			name:       "Get Transaction With Range (user)",
			path:       "/transactionreport/:range",
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
			},
		},
		{
			name: "Get Transaction with range (staff)",
			path: "/transactionreport/:range",
			// login staff
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "azkam@gmail.com",
				"password": "azka123",
			},
		},
	}

	e, db := InitEcho()
	// UserSetup(db)
	// CartSetup(db)
	TransactionSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	transDB := database.NewTransactionDB(db)
	cartDB := database.NewCartDB(db)
	categoryDB := database.NewCategoryDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	transControllers := NewTransactionController(transDB, categoryDB, cartDB, dropPointDB)

	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)
	InsertDataCategory(db)
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
			assert.Equal(t, testCase.expectCode, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("range")
			c.SetParamValues("daily")

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(transControllers.GetTransactionsWithRangeDate)(c)) {
					assert.Equal(t, testCase.expectCode, rec.Code)
					body := rec.Body.String()

					var response = struct {
						Status string                         `json:"status"`
						Data   models.ResponseGetTransactions `json:"data"`
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

func TestGetTransactionWithRangeError(t *testing.T) {
	var testCases = []struct {
		name        string
		path        string
		loginPath   string
		expectCode  int
		expectError string
		login       map[string]interface{}
	}{
		{
			name:        "Get Transaction With Range Invalid Param",
			path:        "/transactionreport/:range",
			loginPath:   "/login",
			expectCode:  http.StatusOK,
			expectError: "Invalid range",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
			},
		},
	}

	e, db := InitEcho()
	// UserSetup(db)
	// CartSetup(db)
	TransactionSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	transDB := database.NewTransactionDB(db)
	cartDB := database.NewCartDB(db)
	categoryDB := database.NewCategoryDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	transControllers := NewTransactionController(transDB, categoryDB, cartDB, dropPointDB)

	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)
	InsertDataCategory(db)
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
			assert.Equal(t, testCase.expectCode, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("range")
			c.SetParamValues("Weeekly")

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(transControllers.GetTransactionsWithRangeDate)(c)
				if assert.Error(t, err) {
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}

func TestGetTransactionTotal(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		loginPath  string
		expectCode int
		response   string
		login      map[string]interface{}
	}{
		{
			name:       "Get Transaction Total User",
			path:       "/totaltransaction",
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
			},
		},
		{
			name: "Get Transaction Total Staff",
			path: "/totaltransaction",
			// login staff
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "azkam@gmail.com",
				"password": "azka123",
			},
		},
	}

	e, db := InitEcho()
	// UserSetup(db)
	// CartSetup(db)
	TransactionSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	transDB := database.NewTransactionDB(db)
	cartDB := database.NewCartDB(db)
	categoryDB := database.NewCategoryDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	transControllers := NewTransactionController(transDB, categoryDB, cartDB, dropPointDB)

	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)
	InsertDataCategory(db)
	InsertDataCartItem(db)
	InsertDataCheckoutVerification(db)

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
			assert.Equal(t, testCase.expectCode, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(transControllers.GetTransactionTotal)(c)) {
					assert.Equal(t, testCase.expectCode, rec.Code)
					body := rec.Body.String()

					var response = struct {
						Status string                 `json:"status"`
						Data   models.ResponseGetUser `json:"data"`
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

func TestGetTransactionTotalWithRangeDate(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		loginPath  string
		expectCode int
		response   string
		login      map[string]interface{}
	}{
		{
			name:       "Get Transaction Total With Range Date User",
			path:       "/totaltransaction/:range",
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
			},
		},
		{
			name: "Get Transaction Total With Range Date Staff",
			path: "/totaltransaction/:range",
			// login staff
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "azkam@gmail.com",
				"password": "azka123",
			},
		},
	}

	e, db := InitEcho()
	// UserSetup(db)
	// CartSetup(db)
	TransactionSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	transDB := database.NewTransactionDB(db)
	cartDB := database.NewCartDB(db)
	categoryDB := database.NewCategoryDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	transControllers := NewTransactionController(transDB, categoryDB, cartDB, dropPointDB)

	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)
	InsertDataCategory(db)
	InsertDataCartItem(db)
	InsertDataCheckoutVerification(db)

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
			assert.Equal(t, testCase.expectCode, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("range")
			c.SetParamValues("monthly")

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(transControllers.GetTransactionTotalWithRangeDate)(c)) {
					assert.Equal(t, testCase.expectCode, rec.Code)
					body := rec.Body.String()

					var response = struct {
						Status string                 `json:"status"`
						Data   models.ResponseGetUser `json:"data"`
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

func TestGetTransactionTotalWithRangeDateError(t *testing.T) {
	var testCases = []struct {
		name        string
		path        string
		loginPath   string
		expectCode  int
		expectError string
		login       map[string]interface{}
	}{
		{
			name:        "Get Transaction Total With Range Date User Invalid Param",
			path:        "/totaltransaction/:range",
			loginPath:   "/login",
			expectCode:  http.StatusOK,
			expectError: "Invalid range",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
			},
		},
	}

	e, db := InitEcho()
	// UserSetup(db)
	// CartSetup(db)
	TransactionSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	transDB := database.NewTransactionDB(db)
	cartDB := database.NewCartDB(db)
	categoryDB := database.NewCategoryDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	transControllers := NewTransactionController(transDB, categoryDB, cartDB, dropPointDB)

	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)
	InsertDataCategory(db)
	InsertDataCartItem(db)
	InsertDataCheckoutVerification(db)

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
			assert.Equal(t, testCase.expectCode, loginRec.Code)
			body := loginRec.Body.String()

			var responseLogin = struct {
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("range")
			c.SetParamValues("Monthlyy")

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(transControllers.GetTransactionTotalWithRangeDate)(c)
				if assert.Error(t, err) {
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}
