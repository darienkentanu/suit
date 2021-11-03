package database

import (
	"database/sql"

	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type TransactionDB struct {
	db *gorm.DB
	dbSQL *sql.DB
}

func NewTransactionDB(db *gorm.DB, dbSQL *sql.DB) *TransactionDB {
	return &TransactionDB{db: db, dbSQL: dbSQL}
}

type TransactionModel interface {
	GetAllTransaction() ([]models.Transaction, error)
}

func (m *TransactionDB) GetAllTransaction() ([]models.Transaction, error) {
	var transaction []models.Transaction

	if err := m.db.Find(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}