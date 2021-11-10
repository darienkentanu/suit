package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// . "github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/gmaps"
	"github.com/darienkentanu/suit/lib/database"
	"github.com/darienkentanu/suit/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func DPSetup(db *gorm.DB) {
	db.Migrator().DropTable(&models.Drop_Point{})
	db.AutoMigrate(&models.Drop_Point{})
}

func InsertDataDropPoints(db *gorm.DB) error {
	dropPoints := models.Drop_Point{
		Address: "universitas padjadjaran",
	}
	dropPoints.Latitude, dropPoints.Longitude = gmaps.Geocoding(dropPoints.Address)
	if err := db.Save(&dropPoints).Error; err != nil {
		return err
	}
	return nil
}

func TestGetDropPoints(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "GetDropPoints",
			path:       "/droppoints",
			expectCode: http.StatusOK,
			response:   "Success",
		},
	}
	e, db, _ := InitEcho()
	DPSetup(db)
	dpdb := database.NewDropPointsDB(db)
	dpc := NewDropPointsController(dpdb)
	InsertDataDropPoints(db)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	for _, testCase := range testCases {
		ctx.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, dpc.GetDropPoints(ctx)) {
				assert.Equal(t, testCase.expectCode, w.Code)
				body := w.Body.String()

				var response = struct {
					Status string                        `json:"status"`
					Data   []models.Drop_Points_Response `json:"data"`
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

func TestAddDropPoints(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
		reqBody    map[string]string
	}{
		{
			name:       "AddDropPoints",
			path:       "/droppoints",
			expectCode: http.StatusCreated,
			response:   "success",
			reqBody: map[string]string{
				"address": "universitas padjadjaran",
			},
		},
	}

	e, db, _ := InitEcho()
	DPSetup(db)
	dpdb := database.NewDropPointsDB(db)
	dpc := NewDropPointsController(dpdb)

	for _, testCase := range testCases {
		reqBody, err := json.Marshal(testCase.reqBody)
		if err != nil {
			t.Error(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, dpc.AddDropPoints(c)) {
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

func TestEditDropPoints(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
		reqBody    map[string]string
	}{
		{
			name:       "EditDropPoints",
			path:       "//droppoints/:id",
			expectCode: http.StatusOK,
			response:   "success",
			reqBody: map[string]string{
				"address": "universitas brawijaya",
			},
		},
	}

	e, db, _ := InitEcho()
	DPSetup(db)
	dpdb := database.NewDropPointsDB(db)
	dpc := NewDropPointsController(dpdb)
	InsertDataDropPoints(db)

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
		c.SetParamValues("1")

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, dpc.EditDropPoints(c)) {
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

func TestEditDropPointError(t *testing.T) {
	var testCases = []struct {
		name        string
		path        string
		expectCode  int
		expectError string
		paramValues string
		reqBody     map[string]interface{}
	}{
		{
			name:        "Edit Drop Point Invalid ID",
			path:        "/droppoints/:id",
			expectCode:  http.StatusBadRequest,
			expectError: "Invalid id",
			paramValues: "a",
			reqBody: map[string]interface{}{
				"address": "universitas brawijaya",
			},
		},
	}

	e, db, _ := InitEcho()
	DPSetup(db)
	dropPointDB := database.NewDropPointsDB(db)
	dropPointControllers := NewDropPointsController(dropPointDB)
	InsertDataDropPoints(db)

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
			err := dropPointControllers.EditDropPoints(c)
			if assert.Error(t, err) {
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})

	}
}

func TestDeleteDropPoints(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "DeleteDropPoints",
			path:       "/droppoints/:id",
			expectCode: http.StatusOK,
			response:   "drop point succesfully deleted",
		},
	}

	e, db, _ := InitEcho()
	CartSetup(db)
	// DPSetup(db)
	dpdb := database.NewDropPointsDB(db)
	dpc := NewDropPointsController(dpdb)
	InsertDataDropPoints(db)

	r := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	for _, testCase := range testCases {
		ctx.SetPath(testCase.path)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, dpc.DeleteDropPoints(ctx)) {
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

func TestDeleteDropPointError(t *testing.T) {
	var testCases = []struct {
		name        string
		path        string
		expectCode  int
		expectError string
		paramValues string
	}{
		{
			name:        "Delete Drop Point Invalid ID",
			path:        "/droppoints/:id",
			expectCode:  http.StatusBadRequest,
			expectError: "Invalid id",
			paramValues: "a",
		},
	}

	e, db, _ := InitEcho()
	DPSetup(db)
	dropPointDB := database.NewDropPointsDB(db)
	dropPointControllers := NewDropPointsController(dropPointDB)
	InsertDataDropPoints(db)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath(testCase.path)
		c.SetParamNames("id")
		c.SetParamValues(testCase.paramValues)

		t.Run(testCase.name, func(t *testing.T) {
			err := dropPointControllers.DeleteDropPoints(c)
			if assert.Error(t, err) {
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})

	}
}
