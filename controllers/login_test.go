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
	"github.com/stretchr/testify/assert"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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
	
	e, db, dbSQL := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db, dbSQL)
	loginDB := database.NewLoginDB(db)
	loginControllers := NewLoginController(userDB, loginDB)
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

func TestGetProfile(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCode 	int
		response   	string
		login		map[string]string
	}{
		{
			name:       "GetProfile",
			path:       "/profile",
			loginPath:	"/login",
			expectCode: http.StatusOK,
			response:   "success",
			login:		map[string]string{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
	}
	
	e, db, dbSQL := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db, dbSQL)
	loginDB := database.NewLoginDB(db)
	loginControllers := NewLoginController(userDB, loginDB)
	// userControllers := NewUserController(userDB, loginDB)
	InsertDataUser(db)

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