package database

import (
	"database/sql"

	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type TransactionDB struct {
	db    *gorm.DB
	dbSQL *sql.DB
}

func NewTransactionDB(db *gorm.DB, dbSQL *sql.DB) *TransactionDB {
	return &TransactionDB{db: db, dbSQL: dbSQL}
}

type TransactionModel interface {
	GetAllTransaction() ([]models.Transaction, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateStatusTransaction(transactionID int) (models.Transaction, error)
}

func (m *TransactionDB) GetAllTransaction() ([]models.Transaction, error) {
	var transaction []models.Transaction

	if err := m.db.Find(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (m *TransactionDB) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {

	if err := m.db.Save(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (m *TransactionDB) UpdateStatusTransaction(transactionID int) (models.Transaction, error) {
	var transaction models.Transaction
	if err := m.db.First(&transaction, transactionID).Update("status", 1).Error; err != nil {
		return transaction, err
	}
	return transaction, nil
}
