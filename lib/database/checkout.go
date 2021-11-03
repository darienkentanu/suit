package database

import (
	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type CheckoutDB struct {
	db *gorm.DB
}

func NewCheckoutDB(db *gorm.DB) *CheckoutDB {
	return &CheckoutDB{db: db}
}

type CheckoutModel interface {
	AddCheckout() (int, error)
}

func (m *CheckoutDB) AddCheckout() (int, error) {
	var checkout models.Checkout
	if err := m.db.Save(&checkout).Error; err != nil {
		return 0, err
	}
	return checkout.ID, nil
}