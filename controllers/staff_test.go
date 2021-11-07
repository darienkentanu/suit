package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func InsertDataStaff(db *gorm.DB) error {
	register := models.RegisterStaff{
		Fullname	: "Muhammad Azka",
		Email		: "azkam@gmail.com",
		Username	: "mazka",
		Password	: "azka123",
		PhoneNumber	: "08126736271",
		DropPointID : 1,
	}
	
	staff := models.Staff{
		Fullname	: register.Fullname,
		PhoneNumber	: register.PhoneNumber,
		Drop_PointID: register.DropPointID,
	}

	if err := db.Save(&staff).Error; err != nil {
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
		Role	: "staff",
		StaffID	: staff.ID,
	}

	if err := db.Select("email", "username", "password", "role", "staff_id").Create(&login).Error; err != nil {
		return err
	}
	
	return nil	
}

func TestRegisterStaff(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		expectCode 	int
		response   	string
		reqBody		map[string]interface{}
	}{
		{
			name:       "RegisterStaff",
			path:       "/registerstaff",
			expectCode: http.StatusCreated,
			response:   "success",
			reqBody: 	map[string]interface{}{
				"fullname": "Muhammad Haikal",
				"email": "mhaikal@gmail.com",
				"username": "mhaikal",
				"password": "haikal123",
				"phone_number": "0812676718",
				"drop_point_id": 1,
			},
		},
	}

	e, db, dbSQL := InitEcho()
	UserSetup(db)
	staffDB := database.NewStaffDB(db, dbSQL)
	loginDB := database.NewLoginDB(db)
	controllers := NewStaffController(staffDB, loginDB)
	InsertDataDropPoints(db)

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
			if assert.NoError(t, controllers.AddStaff(c)) {
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

func TestGetAllStaff(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "GetAllStaff",
			path:       "/staff",
			expectCode: http.StatusOK,
			response:   "success",
		},
	}
	
	e, db, dbSQL := InitEcho()
	UserSetup(db)
	staffDB := database.NewStaffDB(db, dbSQL)
	loginDB := database.NewLoginDB(db)
	controllers := NewStaffController(staffDB, loginDB)
	InsertDataDropPoints(db)
	InsertDataStaff(db)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, controllers.GetAllStaff(c)) {
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