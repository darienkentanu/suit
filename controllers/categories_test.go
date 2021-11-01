package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/darienkentanu/suit/controllers"
	"github.com/darienkentanu/suit/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func InsertDataCategory(db *gorm.DB) error {
	category := models.Category{
		Name:  "gelas kaca",
		Point: 5,
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
		response   int
	}{
		{
			name:       "GetCategories",
			path:       "/categories",
			expectCode: http.StatusOK,
			response:   5,
		},
	}
	e, db, _ := InitEcho()
	InsertDataCategory(db)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(r, w)

	for _, testCase := range testCases {
		ctx.SetPath(testCase.path)

		t.Run(testCase.name, func(t *testing.T) {
			if assert.NoError(t, GetCategories(ctx)) {
				assert.Equal(t, testCase.expectCode, w.Code)
				body := w.Body.String()

				var category models.Category_Response
				err := json.Unmarshal([]byte(body), &category)

				if err != nil {
					assert.Error(t, err, "error")
				}
				assert.Equal(t, testCase.response, 5)
			}
		})
	}
}
