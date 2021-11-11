package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/gmaps"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func UserSetup(db *gorm.DB) {
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
}

func InsertDataUser(db *gorm.DB) error {
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

func TestRegisterUser(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		expectCode 	int
		response   	string
		reqBody		map[string]string
	}{
		{
			name:       "RegisterUser",
			path:       "/register",
			expectCode: http.StatusCreated,
			response:   "success",
			reqBody: 	map[string]string{
				"fullname"		: "Ara Alifia",
				"email"			: "araalifia@gmail.com",
				"username"		: "alifia",
				"password"		: "alifia123",
				"phone_number"	: "0827873486",
				"gender"		: "female",
				"address"		: "Jl. Kebon Jeruk Raya No. 27, Kebon Jeruk, Jakarta Barat 11530",
			},
		},
	}

	e, db := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	cartDB := database.NewCartDB(db)
	controllers := NewUserController(userDB, loginDB, cartDB)

	for _, testCase := range testCases {
		register, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(register))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, controllers.RegisterUsers(c)) {
				assert.Equal(t, testCase.expectCode, rec.Code)
				body := rec.Body.String()

				var response = struct {
					Status string `json:"status"`
					Data   M      `json:"data"`
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

func TestRegisterUserError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		expectCode 	int
		expectError string
		reqBody		map[string]interface{}
	}{
		{
			name:       "Register User Invalid input",
			path:       "/register",
			expectCode: http.StatusBadRequest,
			expectError:   "Invalid input",
			reqBody: 	map[string]interface{}{
				"fullname"		: "Ara Alifia",
				"email"			: "araalifia@gmail.com",
				"username"		: "alifia",
				"password"		: "alifia123",
				"phone_number"	: 827873486,
				"gender"		: "female",
				"address"		: "Jl. Kebon Jeruk Raya No. 27, Kebon Jeruk, Jakarta Barat 11530",
			},
		},
		{
			name:       "Register User Duplicate Email",
			path:       "/register",
			expectCode: http.StatusBadRequest,
			expectError:   "Email is already registered",
			reqBody: 	map[string]interface{}{
				"fullname"		: "Ara Alifia",
				"email"			: "alikatania@gmail.com",
				"username"		: "alifia",
				"password"		: "alifia123",
				"phone_number"	: "0827873486",
				"gender"		: "female",
				"address"		: "Jl. Kebon Jeruk Raya No. 27, Kebon Jeruk, Jakarta Barat 11530",
			},
		},
		{
			name:       "Register User Duplicate Phone Number",
			path:       "/register",
			expectCode: http.StatusBadRequest,
			expectError:   "Phone number is already registered",
			reqBody: 	map[string]interface{}{
				"fullname"		: "Ara Alifia",
				"email"			: "araalifia@gmail.com",
				"username"		: "alifia",
				"password"		: "alifia123",
				"phone_number"	: "08123456789",
				"gender"		: "female",
				"address"		: "Jl. Kebon Jeruk Raya No. 27, Kebon Jeruk, Jakarta Barat 11530",
			},
		},
		{
			name:       "Register User Duplicate Username",
			path:       "/register",
			expectCode: http.StatusBadRequest,
			expectError:   "Username is already registered",
			reqBody: 	map[string]interface{}{
				"fullname"		: "Ara Alifia",
				"email"			: "araalifia@gmail.com",
				"username"		: "alika",
				"password"		: "alifia123",
				"phone_number"	: "0827873486",
				"gender"		: "female",
				"address"		: "Jl. Kebon Jeruk Raya No. 27, Kebon Jeruk, Jakarta Barat 11530",
			},
		},
	}

	e, db := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	cartDB := database.NewCartDB(db)
	controllers := NewUserController(userDB, loginDB, cartDB)
	InsertDataUser(db)

	for _, testCase := range testCases {
		register, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(register))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			err := controllers.RegisterUsers(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
	}
}

func TestRegisterUserInternalServerError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		expectCode 	int
		expectError string
		reqBody		map[string]interface{}
	}{
		{
			name:       "Register User Internal server error",
			path:       "/register",
			expectCode: http.StatusInternalServerError,
			expectError:   "Cannot create user",
			reqBody: 	map[string]interface{}{
				"fullname"		: "Ara Alifia",
				"email"			: "araalifia@gmail.com",
				"username"		: "alifia",
				"password"		: "alifia123",
				"phone_number"	: "0827873486",
				"gender"		: "female",
				"address"		: "Jl. Kebon Jeruk Raya No. 27, Kebon Jeruk, Jakarta Barat 11530",
			},
		},
	}

	e, db := InitEcho()
	UserSetup(db)
	db.Migrator().DropTable(&models.Login{})
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	cartDB := database.NewCartDB(db)
	controllers := NewUserController(userDB, loginDB, cartDB)

	for _, testCase := range testCases {
		register, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(register))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			err := controllers.RegisterUsers(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
	}
}

func TestRegisterUserCartError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		expectCode 	int
		expectError string
		reqBody		map[string]interface{}
	}{
		{
			name:       "Register User Internal server error",
			path:       "/register",
			expectCode: http.StatusInternalServerError,
			expectError:   "Cannot create cart",
			reqBody: 	map[string]interface{}{
				"fullname"		: "Ara Alifia",
				"email"			: "araalifia@gmail.com",
				"username"		: "alifia",
				"password"		: "alifia123",
				"phone_number"	: "0827873486",
				"gender"		: "female",
				"address"		: "Jl. Kebon Jeruk Raya No. 27, Kebon Jeruk, Jakarta Barat 11530",
			},
		},
	}

	e, db := InitEcho()
	UserSetup(db)
	db.Migrator().DropTable(&models.Cart{})
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	cartDB := database.NewCartDB(db)
	controllers := NewUserController(userDB, loginDB, cartDB)

	for _, testCase := range testCases {
		register, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(register))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			err := controllers.RegisterUsers(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "GetAllUsers",
			path:       "/users",
			expectCode: http.StatusOK,
			response:   "success",
		},
	}
	
	e, db := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	cartDB := database.NewCartDB(db)
	controllers := NewUserController(userDB, loginDB, cartDB)
	InsertDataUser(db)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, controllers.GetAllUsers(c)) {
				assert.Equal(t, testCase.expectCode, rec.Code)
				body := rec.Body.String()

				var response = struct {
					Status string					`json:"status"`
					Data   []models.ResponseGetUser `json:"data"`
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

func TestGetAllUsersError(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		expectError   string
	}{
		{
			name:       "Get All Users Error",
			path:       "/users",
			expectCode: http.StatusNotFound,
			expectError:   "Not found",
		},
	}
	
	e, db := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db)
	loginDB := database.NewLoginDB(db)
	cartDB := database.NewCartDB(db)
	controllers := NewUserController(userDB, loginDB, cartDB)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for _, testCase := range testCases {
		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			err := controllers.GetAllUsers(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
	}
}

