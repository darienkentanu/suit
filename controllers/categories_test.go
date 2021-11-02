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

func Setup() {
	_, db, _ := InitEcho()
	db.Migrator().DropTable(&models.Category{})
	db.AutoMigrate(&models.Category{})
}

func InsertDataCategory(db *gorm.DB) error {
	category := models.Category{
		Name:  "pecahan kaca",
		Point: 10,
	}
	if err := db.Save(&category).Error; err != nil {
		return err
	}
	return nil
}

func TestGetCategories(t *testing.T) {
	Setup()
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "GetCategories",
			path:       "/categories",
			expectCode: http.StatusOK,
			response:   "Success",
		},
	}
	e, db, _ := InitEcho()
	cdb := database.NewCategoryDB(db)
	cc := NewCategoryController(cdb)
	InsertDataCategory(db)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	for _, testCase := range testCases {
		ctx.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, cc.GetCategories(ctx)) {
				assert.Equal(t, testCase.expectCode, w.Code)
				body := w.Body.String()

				var response = struct {
					Status string                     `json:"status"`
					Data   []models.Category_Response `json:"data"`
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

func TestAddCategories(t *testing.T) {
	Setup()
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "AddCategories",
			path:       "/categories",
			expectCode: http.StatusCreated,
			response:   "success",
		},
	}

	e, db, _ := InitEcho()
	cdb := database.NewCategoryDB(db)
	cc := NewCategoryController(cdb)
	InsertDataCategory(db)

	reqBody, err := json.Marshal(M{
		"name":  "botol plastik",
		"point": 5,
	})
	if err != nil {
		t.Error(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	for _, testCase := range testCases {
		ctx.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, cc.AddCategories(ctx)) {
				assert.Equal(t, testCase.expectCode, w.Code)
				body := w.Body.String()

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

func TestEditCategories(t *testing.T) {
	Setup()
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "EditCategories",
			path:       "/categories/:id",
			expectCode: http.StatusCreated,
			response:   "success",
		},
	}

	e, db, _ := InitEcho()
	cdb := database.NewCategoryDB(db)
	cc := NewCategoryController(cdb)
	InsertDataCategory(db)

	reqBody, err := json.Marshal(M{
		"name":  "botol plastik",
		"point": 10,
	})
	if err != nil {
		t.Error(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	for _, testCase := range testCases {
		ctx.SetPath(testCase.path)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, cc.EditCategories(ctx)) {
				assert.Equal(t, testCase.expectCode, w.Code)
				body := w.Body.String()

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

func TestDeleteCategories(t *testing.T) {
	Setup()
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "EditCategories",
			path:       "/categories/:id",
			expectCode: http.StatusOK,
			response:   "category succesfully deleted",
		},
	}

	e, db, _ := InitEcho()
	cdb := database.NewCategoryDB(db)
	cc := NewCategoryController(cdb)
	InsertDataCategory(db)

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	for _, testCase := range testCases {
		ctx.SetPath(testCase.path)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, cc.DeleteCategories(ctx)) {
				assert.Equal(t, testCase.expectCode, w.Code)
				body := w.Body.String()

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
