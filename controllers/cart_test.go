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

func CartSetup(db *gorm.DB) {
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.Drop_Point{})
	db.Migrator().DropTable(&models.CartItem{})
	db.Migrator().DropTable(&models.Checkout{})
	db.AutoMigrate(&models.Checkout{})
	db.AutoMigrate(&models.CartItem{})
	db.AutoMigrate(&models.Drop_Point{})
	db.AutoMigrate(&models.Transaction{})
}

func InsertDataCartItem(db *gorm.DB) error {
	cartItemInput := models.CartItem_Input{
		CategoryID: 1,
		Weight:     3,
	}

	var cartItems models.CartItem
	cartItems.CartUserID = 1
	cartItems.CategoryID = cartItemInput.CategoryID
	cartItems.Weight = cartItemInput.Weight

	err := db.Select("cart_user_id", "category_id", "weight").Create(&cartItems).Error
	if err != nil {
		return err
	}

	return nil
}

func TestAddToCart(t *testing.T) {
	var testCases = []struct {
		name            string
		path            string
		loginPath       string
		expectCodeLogin int
		expectCode      int
		response        string
		login           map[string]interface{}
		reqBody         map[string]interface{}
	}{
		{
			name:            "AddCartItem",
			path:            "/cart",
			loginPath:       "/login",
			expectCodeLogin: http.StatusOK,
			expectCode:      http.StatusCreated,
			response:        "success",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
			},
			reqBody: map[string]interface{}{
				"category_id": 1,
				"weight":      3,
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
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	cartController := NewCartController(cartDB)
	InsertDataUser(db)
	InsertDataCategory(db)

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
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
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
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(cartController.AddToCart)(c)) {
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

func TestAddWeightCartItem(t *testing.T) {
	var testCases = []struct {
		name            string
		path            string
		loginPath       string
		expectCodeLogin int
		expectCode      int
		response        string
		login           map[string]interface{}
		reqBody         map[string]interface{}
	}{
		{
			name:            "AddCartItem",
			path:            "/cart",
			loginPath:       "/login",
			expectCodeLogin: http.StatusOK,
			expectCode:      http.StatusCreated,
			response:        "success",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
			},
			reqBody: map[string]interface{}{
				"category_id": 1,
				"weight":      3,
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
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	cartController := NewCartController(cartDB)
	InsertDataUser(db)
	InsertDataCategory(db)
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
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
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
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(cartController.AddToCart)(c)) {
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

func TestGetCartItem(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		loginPath  string
		expectCode int
		response   string
		login      map[string]interface{}
	}{
		{
			name:       "GetCartItem",
			path:       "/cart",
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
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
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	cartController := NewCartController(cartDB)
	InsertDataUser(db)
	InsertDataCategory(db)
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
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(cartController.GetCartItem)(c)) {
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

func TestEditCartItem(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		loginPath  string
		expectCode int
		response   string
		login      map[string]interface{}
		reqBody    map[string]interface{}
	}{
		{
			name:       "EditCartItem",
			path:       "/cart",
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "success",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
			},
			reqBody: map[string]interface{}{
				"category_id": 1,
				"weight":      5,
			},
		},
	}

	e, db := InitEcho()
	Setup(db)
	UserSetup(db)
	CartSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	cartDB := database.NewCartDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	cartController := NewCartController(cartDB)
	InsertDataUser(db)
	InsertDataCategory(db)
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

			reqBody, err := json.Marshal(testCase.reqBody)
			if err != nil {
				t.Error(err)
			}

			req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues("1")

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(cartController.EditCartItem)(c)) {
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

func TestEditCartItemError(t *testing.T) {
	var testCases = []struct {
		name            string
		path            string
		loginPath       string
		expectCodeLogin int
		expectCode      int
		expectError     string
		paramValues     string
		login           map[string]interface{}
		reqBody         map[string]interface{}
	}{
		{
			name:            "Edit Cart Item Invalid ID",
			path:            "/cart",
			loginPath:       "/login",
			expectCodeLogin: http.StatusOK,
			expectCode:      http.StatusBadRequest,
			expectError:     "invalid cart item id",
			paramValues:     "a",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
			},
			reqBody: map[string]interface{}{
				"category_id": 1,
				"weight":      5,
			},
		},
	}

	e, db := InitEcho()
	Setup(db)
	UserSetup(db)
	CartSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	cartDB := database.NewCartDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	cartController := NewCartController(cartDB)
	InsertDataUser(db)
	InsertDataCategory(db)
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
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
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

			req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues(testCase.paramValues)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(cartController.EditCartItem)(c)
				if assert.Error(t, err) {
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}

func TestDeleteCartItem(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		loginPath  string
		expectCode int
		response   string
		login      map[string]interface{}
	}{
		{
			name:       "DeleteCartItem",
			path:       "/cart",
			loginPath:  "/login",
			expectCode: http.StatusOK,
			response:   "cart item succesfully deleted",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
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
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	cartController := NewCartController(cartDB)
	InsertDataUser(db)
	InsertDataCategory(db)
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

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues("1")

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(cartController.DeleteCartItem)(c)) {
					assert.Equal(t, testCase.expectCode, rec.Code)
					body := rec.Body.String()

					var response = struct {
						Message string `json:"message"`
					}{}
					err := json.Unmarshal([]byte(body), &response)

					if err != nil {
						assert.Error(t, err, "error")
					}
					assert.Equal(t, testCase.response, response.Message)
				}
			})
		}
	}
}

func TestDeleteCartItemError(t *testing.T) {
	var testCases = []struct {
		name            string
		path            string
		loginPath       string
		expectCodeLogin int
		expectCode      int
		expectError     string
		paramValues     string
		login           map[string]interface{}
	}{
		{
			name:            "Delete Cart Item Invalid ID",
			path:            "/cart",
			loginPath:       "/login",
			expectCodeLogin: http.StatusOK,
			expectCode:      http.StatusBadRequest,
			expectError:     "invalid cart item id",
			paramValues:     "a",
			login: map[string]interface{}{
				"email":    "alikatania@gmail.com",
				"password": "alika123",
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
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	cartController := NewCartController(cartDB)
	InsertDataUser(db)
	InsertDataCategory(db)
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
				Status string               `json:"status"`
				Data   models.ResponseLogin `json:"data"`
			}{}
			err := json.Unmarshal([]byte(body), &responseLogin)
			if err != nil {
				assert.Error(t, err, "error")
			}

			assert.NotEmpty(t, responseLogin.Data.Token)
			token := responseLogin.Data.Token

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues(testCase.paramValues)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(cartController.DeleteCartItem)(c)
				if assert.Error(t, err) {
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}
