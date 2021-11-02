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
	db.Migrator().DropTable(&models.Drop_Point{})
	db.Migrator().DropTable(&models.Staff{})
	db.Migrator().DropTable(&models.Cart{})
	db.Migrator().DropTable(&models.User{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Cart{})
	db.AutoMigrate(&models.Staff{})
	db.AutoMigrate(&models.Drop_Point{})
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

	e, db, dbSQL := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db, dbSQL)
	loginDB := database.NewLoginDB(db)
	controllers := NewUserController(userDB, loginDB)
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
	
	e, db, dbSQL := InitEcho()
	UserSetup(db)
	userDB := database.NewUserDB(db, dbSQL)
	loginDB := database.NewLoginDB(db)
	controllers := NewUserController(userDB, loginDB)
	InsertDataUser(db)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath(testCase.path)

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