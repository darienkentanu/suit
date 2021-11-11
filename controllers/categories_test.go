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

func Setup(db *gorm.DB) {
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
	e, db := InitEcho()
	Setup(db)
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

func TestGetCategoriesError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		expectCode 	int
		expectError string
	}{
		{
			name:       "Get Categories empty",
			path:       "/categories",
			expectCode: http.StatusNotFound,
			expectError:   "Not found",
		},
	}

	e, db := InitEcho()
	Setup(db)
	db.Migrator().DropTable(&models.Category{})
	categoryDB := database.NewCategoryDB(db)
	controllers := NewCategoryController(categoryDB)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			err := controllers.GetCategories(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
	}
}

func TestAddCategories(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
		reqBody		map[string]interface{}
	}{
		{
			name:       "AddCategories",
			path:       "/categories",
			expectCode: http.StatusCreated,
			response:   "success",
			reqBody: map[string]interface{}{
				"name":  "botol plastik",
				"point": 5,
			},
		},
	}

	e, db := InitEcho()
	Setup(db)
	categoryDB := database.NewCategoryDB(db)
	controllers := NewCategoryController(categoryDB)

	for _, testCase := range testCases {
		category, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(category))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, controllers.AddCategories(c)) {
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

func TestAddCategoriesError(t *testing.T) {
	var testCases = []struct {
		name       	string
		path       	string
		expectCode 	int
		expectError string
		reqBody		map[string]interface{}
	}{
		{
			name:       "Add Categories Invalid Input",
			path:       "/categories",
			expectCode: http.StatusBadRequest,
			expectError:   "Invalid input",
			reqBody: map[string]interface{}{
				"name":  "botol plastik",
				"point": "5",
			},
		},
		{
			name:       "Add Categories Internal server error",
			path:       "/categories",
			expectCode: http.StatusInternalServerError,
			expectError:   "Internal server error",
			reqBody: map[string]interface{}{
				"name":  "botol plastik",
				"point": 5,
			},
		},
	}

	e, db := InitEcho()
	db.Migrator().DropTable(&models.Category{})
	categoryDB := database.NewCategoryDB(db)
	controllers := NewCategoryController(categoryDB)

	for _, testCase := range testCases {
		category, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(category))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			err := controllers.AddCategories(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
	}
}

func TestEditCategories(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
		reqBody		map[string]interface{}
	}{
		{
			name:       "EditCategories",
			path:       "/categories/:id",
			expectCode: http.StatusOK,
			response:   "success",
			reqBody: map[string]interface{}{
				"name":  "botol plastik",
				"point": 5,
			},
		},
	}

	e, db := InitEcho()
	Setup(db)
	categoryDB := database.NewCategoryDB(db)
	controllers := NewCategoryController(categoryDB)
	InsertDataCategory(db)

	for _, testCase := range testCases {
		category, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(category))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)
		c.SetParamNames("id")
		c.SetParamValues("1")

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, controllers.EditCategories(c)) {
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

func TestEditCategoriesError(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		expectCode 		int
		expectError   	string
		paramValues		string
		reqBody			map[string]interface{}
	}{
		{
			name:       "Edit Categories Invalid ID",
			path:       "/categories/:id",
			expectCode: http.StatusBadRequest,
			expectError: "Invalid id",
			paramValues: "a",
			reqBody: 	map[string]interface{}{
				"name":  "botol plastik",
				"point": 10,
			},
		},
		{
			name:       "Edit Categories Invalid Category ID",
			path:       "/categories/:id",
			expectCode: http.StatusNotFound,
			expectError: "Not found",
			paramValues: "10",
			reqBody: 	map[string]interface{}{
				"name":  "botol plastik",
				"point": 10,
			},
		},
	}
	
	e, db := InitEcho()
	Setup(db)
	categoryDB := database.NewCategoryDB(db)
	categoryControllers := NewCategoryController(categoryDB)

	for _, testCase := range testCases {
		reqBody, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)
		c.SetParamNames("id")
		c.SetParamValues(testCase.paramValues)

		t.Run(testCase.name, func(t *testing.T) {
			err := categoryControllers.EditCategories(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
		
	}
}

func TestDeleteCategories(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "DeleteCategories",
			path:       "/categories/:id",
			expectCode: http.StatusOK,
			response:   "category succesfully deleted",
		},
	}

	e, db := InitEcho()
	CartSetup(db)
	Setup(db)
	cdb := database.NewCategoryDB(db)
	cc := NewCategoryController(cdb)
	InsertDataCategory(db)

	r := httptest.NewRequest(http.MethodDelete, "/", nil)
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

func TestDeleteCategoriesError(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		expectCode 		int
		expectError   	string
		paramValues		string
	}{
		{
			name:       "Delete Categories Invalid ID",
			path:       "/categories/:id",
			expectCode: http.StatusBadRequest,
			expectError: "Invalid id",
			paramValues: "a",
		},
		{
			name:       "Delete Categories Invalid ID",
			path:       "/categories/:id",
			expectCode: http.StatusNotFound,
			expectError: "Not found",
			paramValues: "1",
		},
	}
	
	e, db := InitEcho()
	Setup(db)
	categoryDB := database.NewCategoryDB(db)
	categoryControllers := NewCategoryController(categoryDB)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)
		c.SetParamNames("id")
		c.SetParamValues(testCase.paramValues)

		t.Run(testCase.name, func(t *testing.T) {
			err := categoryControllers.DeleteCategories(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
		
	}
}
