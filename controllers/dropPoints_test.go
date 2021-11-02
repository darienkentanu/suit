package controllers_test

import (
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

// func TestAddDropPoints(t *testing.T) {
// 	var testCases = []struct {
// 		name       string
// 		path       string
// 		expectCode int
// 		response   string
// 	}{
// 		{
// 			name:       "AddDropPoints",
// 			path:       "/droppoints",
// 			expectCode: http.StatusCreated,
// 			response:   "success",
// 		},
// 	}

// 	e, db, _ := InitEcho()
// 	DPSetup(db)
// 	dpdb := database.NewDropPointsDB(db)
// 	dpc := NewDropPointsController(dpdb)
// 	InsertDataDropPoints(db)
// 	reqBody, err := json.Marshal(M{
// 		"address": "universitas padjadjaran",
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
// 	w := httptest.NewRecorder()
// 	ctx := e.NewContext(r, w)

// 	for _, testCase := range testCases {
// 		ctx.SetPath(testCase.path)

// 		t.Run(testCase.name, func(t *testing.T) {
// 			if assert.NoError(t, dpc.AddDropPoints(ctx)) {
// 				assert.Equal(t, testCase.expectCode, w.Code)
// 				body := w.Body.String()

// 				var response = struct {
// 					Status string `json:"status"`
// 					Data   M      `json:"data"`
// 				}{}
// 				err := json.Unmarshal([]byte(body), &response)
// 				if err != nil {
// 					assert.Error(t, err, "error")
// 				}
// 				assert.Equal(t, testCase.response, response.Status)
// 			}
// 		})
// 	}
// }

// func TestEditDropPoints(t *testing.T) {
// 	var testCases = []struct {
// 		name       string
// 		path       string
// 		expectCode int
// 		response   string
// 	}{
// 		{
// 			name:       "EditDropPoints",
// 			path:       "/droppoints/:id",
// 			expectCode: http.StatusCreated,
// 			response:   "success",
// 		},
// 	}

// 	e, db, _ := InitEcho()
// 	DPSetup(db)
// 	dpdb := database.NewDropPointsDB(db)
// 	dpc := NewDropPointsController(dpdb)
// 	InsertDataDropPoints(db)

// 	reqBody, err := json.Marshal(M{
// 		"address": "universitas brawijaya",
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
// 	w := httptest.NewRecorder()
// 	ctx := e.NewContext(r, w)

// 	for _, testCase := range testCases {
// 		ctx.SetPath(testCase.path)
// 		ctx.SetParamNames("id")
// 		ctx.SetParamValues("1")

// 		t.Run(testCase.name, func(t *testing.T) {
// 			if assert.NoError(t, dpc.EditDropPoints(ctx)) {
// 				assert.Equal(t, testCase.expectCode, w.Code)
// 				body := w.Body.String()

// 				var response = struct {
// 					Status string `json:"status"`
// 					Data   M      `json:"data"`
// 				}{}
// 				err := json.Unmarshal([]byte(body), &response)
// 				if err != nil {
// 					assert.Error(t, err, "error")
// 				}
// 				assert.Equal(t, testCase.response, response.Status)
// 			}
// 		})
// 	}
// }

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
	DPSetup(db)
	dpdb := database.NewDropPointsDB(db)
	dpc := NewDropPointsController(dpdb)
	InsertDataDropPoints(db)

	r := httptest.NewRequest(http.MethodPost, "/", nil)
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
