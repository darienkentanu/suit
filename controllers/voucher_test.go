package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// . "github.com/darienkentanu/suit/controllers"
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
	e, db, _ := InitEcho()
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

func TestAddVoucher(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "AddVouchers",
			path:       "/vouchers",
			expectCode: http.StatusCreated,
			response:   "success",
		},
	}

	e, db, _ := InitEcho()
	VcSetup(db)
	cdb := database.NewVoucherDB(db)
	cc := NewVoucherController(cdb)
	InsertDataVoucher(db)

	reqBody, err := json.Marshal(M{
		"name":  "voucher pulsa 20k",
		"point": 50,
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
			if assert.NoError(t, cc.AddVouchers(ctx)) {
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

func TestEditVoucher(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		response   string
	}{
		{
			name:       "EditVouchers",
			path:       "/vouchers/:id",
			expectCode: http.StatusCreated,
			response:   "success",
		},
	}

	e, db, _ := InitEcho()
	VcSetup(db)
	cdb := database.NewVoucherDB(db)
	cc := NewVoucherController(cdb)
	InsertDataVoucher(db)

	reqBody, err := json.Marshal(M{
		"name":  "voucher pulsa 10k",
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
			if assert.NoError(t, cc.EditVouchers(ctx)) {
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

func TestEditVoucherError(t *testing.T) {
	var testCases = []struct {
		name        string
		path        string
		expectCode  int
		expectError string
		paramValues string
		reqBody     map[string]interface{}
	}{
		{
			name:        "Edit Vouchers Invalid ID",
			path:        "/vouchers/:id",
			expectCode:  http.StatusBadRequest,
			expectError: "Invalid id",
			paramValues: "a",
			reqBody: map[string]interface{}{
				"name":  "voucher pulsa 10k",
				"point": 10,
			},
		},
	}

	e, db, _ := InitEcho()
	VcSetup(db)
	voucherDB := database.NewVoucherDB(db)
	voucherControllers := NewVoucherController(voucherDB)
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
		c.SetParamValues(testCase.paramValues)

		t.Run(testCase.name, func(t *testing.T) {
			err := voucherControllers.EditVouchers(c)
			if assert.Error(t, err) {
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

	e, db, _ := InitEcho()
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
		name        string
		path        string
		expectCode  int
		expectError string
		paramValues string
	}{
		{
			name:        "Delete Vouchers Invalid ID",
			path:        "/vouchers/:id",
			expectCode:  http.StatusBadRequest,
			expectError: "Invalid id",
			paramValues: "a",
		},
	}

	e, db, _ := InitEcho()
	VcSetup(db)
	voucherDB := database.NewVoucherDB(db)
	voucherControllers := NewVoucherController(voucherDB)
	InsertDataVoucher(db)

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath(testCase.path)
		c.SetParamNames("id")
		c.SetParamValues(testCase.paramValues)

		t.Run(testCase.name, func(t *testing.T) {
			err := voucherControllers.DeleteVouchers(c)
			if assert.Error(t, err) {
				assert.Containsf(t, err.Error(), testCase.expectError, "expected error containing %q, got %s", testCase.expectError, err)
			}
		})

	}
}
