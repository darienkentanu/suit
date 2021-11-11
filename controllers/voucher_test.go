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

func VcSetup(db *gorm.DB) {
	db.Migrator().DropTable(&models.Voucher{})
	db.AutoMigrate(&models.Voucher{})
}

func InsertDataVoucher(db *gorm.DB) error {
	voucher := models.Voucher{
		Name:  "voucher pulsa 10k",
		Point: 20,
	}
	if err := db.Save(&voucher).Error; err != nil {
		return err
	}
	return nil
}

func TestGetVoucher(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "GetAllVouchers",
			path:       "/vouchers",
			expectCode: http.StatusOK,
			response:   "Success",
		},
	}
	e, db  := InitEcho()
	UserVoucherSetup(db)
	VcSetup(db)
	cdb := database.NewVoucherDB(db)
	cc := NewVoucherController(cdb)
	InsertDataVoucher(db)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	for _, testCase := range testCases {
		ctx.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, cc.GetVouchers(ctx)) {
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

func TestGetVoucherError(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		expectCode 		int
		expectError   	string
	}{
		{
			name:       "Get Vouchers Not found",
			path:       "/vouchers/:id",
			expectCode: http.StatusInternalServerError,
			expectError: "Not found",
		},
	}
	
	e, db  := InitEcho()
	db.Migrator().DropTable(&models.Voucher{})
	voucherDB := database.NewVoucherDB(db)
	voucherControllers := NewVoucherController(voucherDB)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(nil))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			err := voucherControllers.GetVouchers(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
		
	}
}

func TestAddVoucher(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
		reqBody    map[string]interface{}
	}{
		{
			name:       "AddVouchers",
			path:       "/vouchers",
			expectCode: http.StatusCreated,
			response:   "success",
			reqBody: 	map[string]interface{}{
				"name":  "voucher pulsa 10k",
				"point": 10,
			},
		},
	}

	e, db  := InitEcho()
	VcSetup(db)
	cdb := database.NewVoucherDB(db)
	cc := NewVoucherController(cdb)

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
			if assert.NoError(t, cc.AddVouchers(c)) {
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

func TestAddVoucherError(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		expectCode 		int
		expectError   	string
		reqBody			map[string]interface{}
	}{
		{
			name:       "Add Vouchers Invalid Input",
			path:       "/vouchers/:id",
			expectCode: http.StatusBadRequest,
			expectError: "Invalid input",
			reqBody: 	map[string]interface{}{
				"name":  "voucher pulsa 10k",
				"point": "10",
			},
		},
		{
			name:       "Add Vouchers Internal server error",
			path:       "/vouchers/:id",
			expectCode: http.StatusInternalServerError,
			expectError: "Internal server error",
			reqBody: 	map[string]interface{}{
				"name":  "voucher pulsa 10k",
				"point": 10,
			},
		},
	}
	
	e, db  := InitEcho()
	db.Migrator().DropTable(&models.Voucher{})
	voucherDB := database.NewVoucherDB(db)
	voucherControllers := NewVoucherController(voucherDB)

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
			err := voucherControllers.AddVouchers(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
		
	}
}

func TestEditVoucher(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
		reqBody	   map[string]interface{}
	}{
		{
			name:       "EditVouchers",
			path:       "/vouchers/:id",
			expectCode: http.StatusOK,
			response:   "success",
			reqBody: 	map[string]interface{}{
				"name":  "voucher pulsa 10k",
				"point": 15,
			},
		},
	}

	e, db  := InitEcho()
	VcSetup(db)
	cdb := database.NewVoucherDB(db)
	cc := NewVoucherController(cdb)
	InsertDataVoucher(db)

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
			if assert.NoError(t, cc.EditVouchers(c)) {
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

func TestEditVoucherError(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		expectCode 		int
		expectError   	string
		paramValues		string
		reqBody			map[string]interface{}
	}{
		{
			name:       "Edit Vouchers Invalid ID",
			path:       "/vouchers/:id",
			expectCode: http.StatusBadRequest,
			expectError: "Invalid id",
			paramValues: "a",
			reqBody: 	map[string]interface{}{
				"name":  "voucher pulsa 10k",
				"point": 10,
			},
		},
		{
			name:       "Edit Vouchers Invalid Input",
			path:       "/vouchers/:id",
			expectCode: http.StatusBadRequest,
			expectError: "Invalid input",
			paramValues: "1",
			reqBody: 	map[string]interface{}{
				"name":  "voucher pulsa 10k",
				"point": "10",
			},
		},
		{
			name:       "Edit Vouchers Not found",
			path:       "/vouchers/:id",
			expectCode: http.StatusNotFound,
			expectError: "Not found",
			paramValues: "20",
			reqBody: 	map[string]interface{}{
				"name":  "voucher pulsa 10k",
				"point": 10,
			},
		},
	}
	
	e, db  := InitEcho()
	VcSetup(db)
	voucherDB := database.NewVoucherDB(db)
	voucherControllers := NewVoucherController(voucherDB)

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
			err := voucherControllers.EditVouchers(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
		
	}
}

func TestDeleteVoucher(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "DeleteVouchers",
			path:       "/vouchers/:id",
			expectCode: http.StatusOK,
			response:   "voucher succesfully deleted",
		},
	}

	e, db  := InitEcho()
	VcSetup(db)
	cdb := database.NewVoucherDB(db)
	cc := NewVoucherController(cdb)
	InsertDataVoucher(db)

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	for _, testCase := range testCases {
		ctx.SetPath(testCase.path)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, cc.DeleteVouchers(ctx)) {
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

func TestDeleteVoucherError(t *testing.T) {
	var testCases = []struct {
		name       		string
		path       		string
		expectCode 		int
		expectError   	string
		paramValues		string
	}{
		{
			name:       "Delete Vouchers Invalid ID",
			path:       "/vouchers/:id",
			expectCode: http.StatusBadRequest,
			expectError: "Invalid id",
			paramValues: "a",
		},
		{
			name:       "Delete Vouchers Invalid Voucher ID",
			path:       "/vouchers/:id",
			expectCode: http.StatusNotFound,
			expectError: "Not found",
			paramValues: "1",
		},
	}
	
	e, db  := InitEcho()
	VcSetup(db)
	voucherDB := database.NewVoucherDB(db)
	voucherControllers := NewVoucherController(voucherDB)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		c.SetPath(testCase.path)
		c.SetParamNames("id")
		c.SetParamValues(testCase.paramValues)

		t.Run(testCase.name, func(t *testing.T) {
			err := voucherControllers.DeleteVouchers(c)
			if assert.Error(t, err){
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})
		
	}
}
