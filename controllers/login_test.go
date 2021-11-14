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
)

func TestLogin(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCode 	int
		response   	string
		reqBody		map[string]string
	}{
		{
			name:       "Login",
			path:       "/login",
			expectCode: http.StatusOK,
			response:   "success",
			reqBody:	map[string]string{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	InsertDataUser(db)

	for _, testCase := range testCases {
		login, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		loginReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(login))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		loginC := e.NewContext(loginReq, loginRec)
		
		loginC.SetPath(testCase.loginPath)

		assert.Equal(t, testCase.expectCode, loginRec.Code)
		body := loginRec.Body.String()

		var responseLogin = struct {
			Status string					`json:"status"`
			Data   models.ResponseLogin 	`json:"data"`
		}{}

		err = json.Unmarshal([]byte(body), &responseLogin)
		if err != nil {
			assert.Error(t, err, "error")
		}

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, loginControllers.Login(loginC)){
				assert.Equal(t, testCase.expectCode, loginRec.Code)
				body := loginRec.Body.String()

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

func TestLoginError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCode 	int
		expectError	string
		reqBody		map[string]interface{}
	}{
		{
			name:       "Login Bind Error",
			path:       "/login",
			expectCode: http.StatusBadRequest,
			expectError: "Invalid input",
			reqBody:	map[string]interface{}{
				"email"			: "alikataniaaa@gmail.com",
				"password"		: 12345,
			},
		},
		{
			name:       "Login Incorrect Email",
			path:       "/login",
			expectCode: http.StatusBadRequest,
			expectError: "Incorrect email or username",
			reqBody:	map[string]interface{}{
				"email"			: "alikataniaaa@gmail.com",
				"password"		: "alika123",
			},
		},
		{
			name:       "Login Incorrect Password",
			path:       "/login",
			expectCode: http.StatusBadRequest,
			expectError: "Incorrect password",
			reqBody:	map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika12345",
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	InsertDataUser(db)

	for _, testCase := range testCases {
		login, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		loginReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(login))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		loginC := e.NewContext(loginReq, loginRec)
		
		loginC.SetPath(testCase.loginPath)

		t.Run(testCase.name, func(t *testing.T) {
			err := loginControllers.Login(loginC)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
		
	}
}

func TestGetProfile(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCode 	int
		response   	string
		login		map[string]interface{}
	}{
		{
			name:       "Get Profile User",
			path:       "/profile",
			loginPath:	"/login",
			expectCode: http.StatusOK,
			response:   "success",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
		{
			name:       "Get Profile Staff",
			path:       "/profile",
			loginPath:	"/login",
			expectCode: http.StatusOK,
			response:   "success",
			login:		map[string]interface{}{
				"email"			: "azkam@gmail.com",
				"password"		: "azka123",
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)

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
				Status string					`json:"status"`
				Data   models.ResponseLogin 	`json:"data"`
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
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(loginControllers.GetProfile)(c)){
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

func TestGetProfileError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCodeLogin int
		expectCode 	int
		expectError string
		login		map[string]interface{}
	}{
		{
			name:       "Get Profile User Internal Server Error",
			path:       "/profile",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusInternalServerError,
			expectError:   "Internal server error",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
		{
			name:       "Get Profile Staff Internal Server Error",
			path:       "/profile",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusInternalServerError,
			expectError:   "Internal server error",
			login:		map[string]interface{}{
				"email"			: "azkam@gmail.com",
				"password"		: "azka123",
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)
	db.Migrator().DropTable(&models.Staff{})
	db.Migrator().DropTable(&models.User{})

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

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(loginControllers.GetProfile)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}

func TestUpdateProfile(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCode 	int
		response   	string
		login		map[string]interface{}
		reqBody		map[string]interface{}
	}{
		{
			name:       "Update Profile User",
			path:       "/profile",
			loginPath:	"/login",
			expectCode: http.StatusOK,
			response:   "success",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
			reqBody: 	map[string]interface{}{
				"fullname"		: "Alika Tania P.",
				"email"			: "alikataniap@gmail.com",
				"username"		: "alikap",
				"password"		: "alika123",
				"phone_number"	: "0812783781",
				"gender"		: "female",
				"address"		: "Jl. Kebon Jeruk Raya No. 27, Kebon Jeruk, Jakarta Barat 11530",
			},
		},
		{
			name:       "Update Profile Staff",
			path:       "/profile",
			loginPath:	"/login",
			expectCode: http.StatusOK,
			response:   "success",
			login:		map[string]interface{}{
				"email"			: "azkam@gmail.com",
				"password"		: "azka123",
			},
			reqBody: 	map[string]interface{}{
				"fullname": "Muhammad Azka R.",
				"email": 	"azkamr@gmail.com",
				"username": "mazkar",
				"password": "azka1234",
				"phone_number": "08126736171",
				"drop_point_id": 1,
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)

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

			req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(loginControllers.UpdateProfile)(c)){
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

func TestUpdateProfileError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCodeLogin int
		expectCode 	int
		expectError string
		login		map[string]interface{}
		reqBody		map[string]interface{}
	}{
		{
			name:       "Update Profile User Invalid Input",
			path:       "/profile",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Invalid input",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
			reqBody: 	map[string]interface{}{
				"fullname"		: "Alika Tania P.",
				"email"			: "alikataniap@gmail.com",
				"username"		: "alikap",
				"password"		: "alika123",
				"phone_number"	: 812783781,
				"gender"		: "female",
				"address"		: "Jl. Kebon Jeruk Raya No. 27, Kebon Jeruk, Jakarta Barat 11530",
			},
		},
		{
			name:       "Update Profile Staff Invalid Input",
			path:       "/profile",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Invalid input",
			login:		map[string]interface{}{
				"email"			: "azkam@gmail.com",
				"password"		: "azka123",
			},
			reqBody: 	map[string]interface{}{
				"fullname": "Muhammad Azka R.",
				"email": 	"azkamr@gmail.com",
				"username": "mazkar",
				"password": "azka1234",
				"phone_number": 8126736171,
				"drop_point_id": 1,
			},
		},
		{
			name:       "Update Profile User Internal server error",
			path:       "/profile",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusInternalServerError,
			expectError:   "Internal server error",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
			reqBody: 	map[string]interface{}{
				"fullname"		: "Alika Tania P.",
				"email"			: "alikataniap@gmail.com",
				"username"		: "alikap",
				"password"		: "alika123",
				"phone_number"	: "08112783781",
				"gender"		: "female",
				"address"		: "Jl. Kebon Jeruk Raya No. 27, Kebon Jeruk, Jakarta Barat 11530",
			},
		},
		{
			name:       "Update Profile Staff Internal server error",
			path:       "/profile",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Internal server error",
			login:		map[string]interface{}{
				"email"			: "azkam@gmail.com",
				"password"		: "azka123",
			},
			reqBody: 	map[string]interface{}{
				"fullname": "Muhammad Azka R.",
				"email": 	"azkamr@gmail.com",
				"username": "mazkar",
				"password": "azka1234",
				"phone_number": "08126736171",
				"drop_point_id": 1,
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	InsertDataUser(db)
	InsertDataDropPoints(db)
	InsertDataStaff(db)
	db.Migrator().DropTable(&models.Staff{})
	db.Migrator().DropTable(&models.User{})

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

			req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(loginControllers.UpdateProfile)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}

func TestUpdateProfileDropPointError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCodeLogin int
		expectCode 	int
		expectError string
		login		map[string]interface{}
		reqBody		map[string]interface{}
	}{
		{
			name:       "Update Profile Staff Internal server error",
			path:       "/profile",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Internal server error",
			login:		map[string]interface{}{
				"email"			: "azkam@gmail.com",
				"password"		: "azka123",
			},
			reqBody: 	map[string]interface{}{
				"fullname": "Muhammad Azka R.",
				"email": 	"azkamr@gmail.com",
				"username": "mazkar",
				"password": "azka1234",
				"phone_number": "08126736171",
				"drop_point_id": 1,
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	InsertDataDropPoints(db)
	InsertDataStaff(db)
	db.Migrator().DropTable(&models.Drop_Point{})

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

			req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(loginControllers.UpdateProfile)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}