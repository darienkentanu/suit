package database

import (
	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type CartDB struct {
	db *gorm.DB
}

func NewCartDB(db *gorm.DB) *CartDB {
	return &CartDB{db: db}
}

type CartModel interface {
	CreateCart(cart models.Cart) (error)
}

func (m *CartDB) CreateCart(cart models.Cart) error {
	if err := m.db.Save(&cart).Error; err != nil {
		return err
	}
	return nil
}