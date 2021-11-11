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
	"github.com/darienkentanu/suit/gmaps"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func UserVoucherSetup(db *gorm.DB) {
	db.Migrator().DropTable(&models.User_Voucher{})
	db.AutoMigrate(&models.User_Voucher{})
}

func InsertDataUserVoucher(db *gorm.DB) error {
	register := models.RegisterUser{
		Fullname	: "Alika Tania",
		Email		: "alikatania@gmail.com",
		Username	: "alika",
		Password	: "alika123",
		PhoneNumber	: "08123456789",
		Gender		: "female",
		Address		: "Jl. Margonda Raya, Pondok Cina, Kecamatan Beji, Kota Depok, Jawa Barat 16424",
	}

	lat, lng := gmaps.Geocoding(register.Address)
	
	user := models.User{
		Fullname	: register.Fullname,
		PhoneNumber	: register.PhoneNumber,
		Gender		: register.Gender,
		Address		: register.Address,
		Point		: 200,
		Latitude	: lat,
		Longitude	: lng,
	}

	if err := db.Save(&user).Error; err != nil {
		return err
	}

	hashPassword, err := GenerateHashPassword(register.Password)
	if err != nil {
		return err
	}

	login := models.Login{
		Email	: register.Email,
		Username: register.Username,
		Password: hashPassword,
		Role	: "user",
		UserID	: user.ID,
	}

	if err := db.Select("email", "username", "password", "role", "user_id").Create(&login).Error; err != nil {
		return err
	}

	var cart models.Cart
	cart.UserID = 1
	err = db.Select("user_id").Create(&cart).Error
	if err != nil {
		return err
	}

	voucher := models.Voucher{
		Name:  "voucher pulsa 10k",
		Point: 20,
	}
	if err := db.Save(&voucher).Error; err != nil {
		return err
	}

	userVoucher := models.User_Voucher{
		VoucherID: 1,
		UserID: 1,
		Status: "available",
	}

	if err := db.Save(&userVoucher).Error; err != nil {
		return err
	}

	return nil	
}

func InsertDataUserWithPoints(db *gorm.DB) error {
	register := models.RegisterUser{
		Fullname	: "Alika Tania",
		Email		: "alikatania@gmail.com",
		Username	: "alika",
		Password	: "alika123",
		PhoneNumber	: "08123456789",
		Gender		: "female",
		Address		: "Jl. Margonda Raya, Pondok Cina, Kecamatan Beji, Kota Depok, Jawa Barat 16424",
	}

	lat, lng := gmaps.Geocoding(register.Address)
	
	user := models.User{
		Fullname	: register.Fullname,
		PhoneNumber	: register.PhoneNumber,
		Gender		: register.Gender,
		Address		: register.Address,
		Point		: 200,
		Latitude	: lat,
		Longitude	: lng,
	}

	if err := db.Save(&user).Error; err != nil {
		return err
	}

	hashPassword, err := GenerateHashPassword(register.Password)
	if err != nil {
		return err
	}

	login := models.Login{
		Email	: register.Email,
		Username: register.Username,
		Password: hashPassword,
		Role	: "user",
		UserID	: user.ID,
	}

	if err := db.Select("email", "username", "password", "role", "user_id").Create(&login).Error; err != nil {
		return err
	}

	var cart models.Cart
	cart.UserID = 1
	err = db.Select("user_id").Create(&cart).Error
	if err != nil {
		return err
	}

	return nil	
}

func TestClaimVoucher(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCodeLogin int
		expectCode 	int
		response   	string
		login		map[string]interface{}
	}{
		{
			name:       "Claim Voucher",
			path:       "/claim/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusCreated,
			response:   "success",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
	}
	
	e, db := InitEcho()
	UserSetup(db)
	VcSetup(db)
	UserVoucherSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	userVoucherDB := database.NewUserVoucherDB(db)
	voucherDB := database.NewVoucherDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	userVoucherControllers := NewUserVoucherController(userVoucherDB, userDB, voucherDB)
	InsertDataUserWithPoints(db)
	InsertDataVoucher(db)

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

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues("1")

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(userVoucherControllers.ClaimVoucher)(c)){
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

func TestClaimVoucherError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCodeLogin int
		expectCode 	int
		expectError   	string
		paramValues	string
		login		map[string]interface{}
	}{
		{
			name:       "Claim Voucher Invalid ID",
			path:       "/claim/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Invalid input",
			paramValues: "a",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
		{
			name:       "Claim Voucher Not Enough Points",
			path:       "/claim/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Not enough point",
			paramValues: "1",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
		{
			name:       "Claim Voucher Invalid Voucher ID",
			path:       "/claim/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusInternalServerError,
			expectError:   "Internal server error",
			paramValues: "5",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
	}
	
	e, db := InitEcho()
	UserSetup(db)
	VcSetup(db)
	UserVoucherSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	userVoucherDB := database.NewUserVoucherDB(db)
	voucherDB := database.NewVoucherDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	userVoucherControllers := NewUserVoucherController(userVoucherDB, userDB, voucherDB)
	InsertDataUser(db)
	InsertDataVoucher(db)

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

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues(testCase.paramValues)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(userVoucherControllers.ClaimVoucher)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}

func TestClaimVoucherUserError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCodeLogin int
		expectCode 	int
		expectError   	string
		paramValues	string
		login		map[string]interface{}
	}{
		{
			name:       "Claim Voucher Internal server error",
			path:       "/claim/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusInternalServerError,
			expectError:   "Internal server error",
			paramValues: "1",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
	}
	
	e, db := InitEcho()
	UserSetup(db)
	VcSetup(db)
	UserVoucherSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	userVoucherDB := database.NewUserVoucherDB(db)
	voucherDB := database.NewVoucherDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	userVoucherControllers := NewUserVoucherController(userVoucherDB, userDB, voucherDB)
	InsertDataUserVoucher(db)
	InsertDataVoucher(db)
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

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues(testCase.paramValues)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(userVoucherControllers.ClaimVoucher)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}

func TestClaimVoucherUserVoucherError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCodeLogin int
		expectCode 	int
		expectError   	string
		paramValues	string
		login		map[string]interface{}
	}{
		{
			name:       "Claim Voucher Internal server error",
			path:       "/claim/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusInternalServerError,
			expectError:   "Internal server error",
			paramValues: "1",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
	}
	
	e, db := InitEcho()
	UserSetup(db)
	VcSetup(db)
	UserVoucherSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	userVoucherDB := database.NewUserVoucherDB(db)
	voucherDB := database.NewVoucherDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	userVoucherControllers := NewUserVoucherController(userVoucherDB, userDB, voucherDB)
	InsertDataUserVoucher(db)
	InsertDataVoucher(db)
	db.Migrator().DropTable(&models.User_Voucher{})

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

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)
			c.SetParamNames("id")
			c.SetParamValues(testCase.paramValues)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(userVoucherControllers.ClaimVoucher)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}

func TestRedeemVoucher(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCodeLogin int
		expectCode 	int
		response   	string
		login		map[string]interface{}
	}{
		{
			name:       "Redeem Voucher",
			path:       "/redeem/:id",
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
	VcSetup(db)
	UserVoucherSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	userVoucherDB := database.NewUserVoucherDB(db)
	voucherDB := database.NewVoucherDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	userVoucherControllers := NewUserVoucherController(userVoucherDB, userDB, voucherDB)
	InsertDataUserVoucher(db)

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
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(userVoucherControllers.RedeemVoucher)(c)){
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

func TestRedeemVoucherError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCodeLogin int
		expectCode 	int
		expectError string
		paramValues	string
		login		map[string]interface{}
	}{
		{
			name:       "Redeem Voucher Invalid ID",
			path:       "/redeem/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Invalid input",
			paramValues: "a",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
		{
			name:       "Redeem Voucher Not Available",
			path:       "/redeem/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusBadRequest,
			expectError:   "Not available",
			paramValues: "10",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	VcSetup(db)
	UserVoucherSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	userVoucherDB := database.NewUserVoucherDB(db)
	voucherDB := database.NewVoucherDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	userVoucherControllers := NewUserVoucherController(userVoucherDB, userDB, voucherDB)
	InsertDataUserVoucher(db)

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
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(userVoucherControllers.RedeemVoucher)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}

func TestGetUserVoucher(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		loginPath	string
		expectCodeLogin int
		expectCode 	int
		response   	string
		login		map[string]interface{}
	}{
		{
			name:       "Get User Voucher",
			path:       "/uservouchers",
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
	VcSetup(db)
	UserVoucherSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	userVoucherDB := database.NewUserVoucherDB(db)
	voucherDB := database.NewVoucherDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	userVoucherControllers := NewUserVoucherController(userVoucherDB, userDB, voucherDB)
	InsertDataUserVoucher(db)

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

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				if assert.NoError(t, echoMiddleware.JWT([]byte(constants.JWT_SECRET))(userVoucherControllers.GetUserVoucher)(c)){
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

func TestGetUserVoucherError(t *testing.T) {
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
			name:       "Get User Voucher Internal Server Error",
			path:       "/uservouchers/:id",
			loginPath:	"/login",
			expectCodeLogin: http.StatusOK,
			expectCode: http.StatusInternalServerError,
			expectError:   "Internal server error",
			login:		map[string]interface{}{
				"email"			: "alikatania@gmail.com",
				"password"		: "alika123",
			},
		},
	}
	
	e, db  := InitEcho()
	UserSetup(db)
	VcSetup(db)
	UserVoucherSetup(db)
	db.Migrator().DropTable(&models.User_Voucher{})
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	staffDB := database.NewStaffDB(db)
	dropPointDB := database.NewDropPointsDB(db)
	userVoucherDB := database.NewUserVoucherDB(db)
	voucherDB := database.NewVoucherDB(db)
	loginControllers := NewLoginController(userDB, loginDB, staffDB, dropPointDB)
	userVoucherControllers := NewUserVoucherController(userVoucherDB, userDB, voucherDB)
	InsertDataUserVoucher(db)

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

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath(testCase.path)

			t.Run(testCase.name, func(t *testing.T) {
				err := echoMiddleware.JWT([]byte(constants.JWT_SECRET))(userVoucherControllers.GetUserVoucher)(c)
				if assert.Error(t, err){
					assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
				}
			})
		}
	}
}