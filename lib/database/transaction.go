package database

import (
	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type TransactionDB struct {
	db    *gorm.DB
}

func NewTransactionDB(db *gorm.DB) *TransactionDB {
	return &TransactionDB{db: db}
}

type TransactionModel interface {
	GetAllTransaction() ([]models.Transaction, error)
	GetTransactionsByUserID(userID int) ([]models.Transaction, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateStatusTransaction(transactionID int) (models.Transaction, error)
	GetTransationsRangeDate(id int, role, rangeDate string) ([]models.TransactionSQL, error)
	GetTransationTotalRangeDate(id int, role, rangeDate string) ([]models.ResponseCategoryReport, error)
	GetTransationTotal(id int, role string) ([]models.ResponseCategoryReport, error)
	GetTransactionsByDropPointID(dropPointID int) ([]models.TransactionSQL, error)
}

func (m *TransactionDB) GetAllTransaction() ([]models.Transaction, error) {
	var transactions []models.Transaction

	if err := m.db.Find(&transactions).Error; err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (m *TransactionDB) GetTransactionsByUserID(userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	if err := m.db.Raw("SELECT * FROM transactions WHERE created_at >= DATE_ADD(CURDATE(), INTERVAL -1 DAY)").Scan(&transactions).Error; err != nil {
		return nil, err
	}
	if err := m.db.Find(&transactions, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (m *TransactionDB) GetTransactionsByDropPointID(dropPointID int) ([]models.TransactionSQL, error) {
	var transactions []models.TransactionSQL
	
	if err := m.db.Raw("SELECT t.*, dp.address as drop_point_address FROM transactions t JOIN drop_points dp ON t.drop_point_id = dp.id WHERE t.drop_point_id = ?", dropPointID).Scan(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
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

func (m *TransactionDB) GetTransationsRangeDate(id int, role, rangeDate string) ([]models.TransactionSQL, error) {
	var transactions []models.TransactionSQL
	
	if role == "staff" {
		if rangeDate == "daily" {
			if err := m.db.Raw("SELECT t.*, dp.address as drop_point_address FROM transactions t JOIN drop_points dp ON t.drop_point_id = dp.id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -1 DAY)").Scan(&transactions).Error; err != nil {
				return nil, err
			}
		} else if rangeDate == "weekly" {
			if err := m.db.Raw("SELECT t.*, dp.address as drop_point_address FROM transactions t JOIN drop_points dp ON t.drop_point_id = dp.id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -7 DAY)").Scan(&transactions).Error; err != nil {
				return nil, err
			}
		} else if rangeDate == "monthly" {
			if err := m.db.Raw("SELECT t.*, dp.address as drop_point_address FROM transactions t JOIN drop_points dp ON t.drop_point_id = dp.id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -30 DAY)").Scan(&transactions).Error; err != nil {
				return nil, err
			}
		}
	} else if role == "user" {
		if rangeDate == "daily" {
			if err := m.db.Raw("SELECT t.*, dp.address as drop_point_address FROM transactions t JOIN drop_points dp ON t.drop_point_id = dp.id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -1 DAY) AND user_id = ?", id).Scan(&transactions).Error; err != nil {
				return nil, err
			}
		} else if rangeDate == "weekly" {
			if err := m.db.Raw("SELECT t.*, dp.address as drop_point_address FROM transactions t JOIN drop_points dp ON t.drop_point_id = dp.id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -7 DAY) AND user_id = ?", id).Scan(&transactions).Error; err != nil {
				return nil, err
			}
		} else if rangeDate == "monthly" {
			if err := m.db.Raw("SELECT t.*, dp.address as drop_point_address FROM transactions t JOIN drop_points dp ON t.drop_point_id = dp.id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -30 DAY) AND user_id = ?", id).Scan(&transactions).Error; err != nil {
				return nil, err
			}
		}
	}

	return transactions, nil
}

func (m *TransactionDB) GetTransationTotalRangeDate(id int, role, rangeDate string) ([]models.ResponseCategoryReport, error) {
	var categories []models.ResponseCategoryReport
	
	if role == "staff" {
		if rangeDate == "daily" {
			if err := m.db.Raw("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -1 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 GROUP BY ci.category_id ORDER BY ci.category_id").Scan(&categories).Error; err != nil {
				return nil, err
			}
		} else if rangeDate == "weekly" {
			if err := m.db.Raw("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -7 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 GROUP BY ci.category_id ORDER BY ci.category_id").Scan(&categories).Error; err != nil {
				return nil, err
			}
		} else if rangeDate == "monthly" {
			if err := m.db.Raw("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -30 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 GROUP BY ci.category_id ORDER BY ci.category_id").Scan(&categories).Error; err != nil {
				return nil, err
			}
		}
	} else if role == "user" {
		if rangeDate == "daily" {
			if err := m.db.Raw("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE ci.created_at >= DATE_ADD(CURDATE(), INTERVAL -1 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 AND ci.cart_user_id = ? GROUP BY ci.category_id ORDER BY ci.category_id", id).Scan(&categories).Error; err != nil {
				return nil, err
			}
		} else if rangeDate == "weekly" {
			if err := m.db.Raw("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE ci.created_at >= DATE_ADD(CURDATE(), INTERVAL -7 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 AND ci.cart_user_id = ? GROUP BY ci.category_id ORDER BY ci.category_id", id).Scan(&categories).Error; err != nil {
				return nil, err
			}
		} else if rangeDate == "monthly" {
			if err := m.db.Raw("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE ci.created_at >= DATE_ADD(CURDATE(), INTERVAL -30 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 AND ci.cart_user_id = ? GROUP BY ci.category_id ORDER BY ci.category_id", id).Scan(&categories).Error; err != nil {
				return nil, err
			}
		}
	}

	return categories, nil
}

func (m *TransactionDB) GetTransationTotal(id int, role string) ([]models.ResponseCategoryReport, error) {
	var categories []models.ResponseCategoryReport
	
	if role == "staff" {
		if err := m.db.Raw("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE ci.checkout_id IS NOT NULL AND t.status = 1 GROUP BY ci.category_id ORDER BY ci.category_id").Scan(&categories).Error; err != nil {
			return nil, err
		}
	} else if role == "user" {
		if err := m.db.Raw("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE ci.checkout_id IS NOT NULL AND t.status = 1 AND ci.cart_user_id = ? GROUP BY ci.category_id ORDER BY ci.category_id", id).Scan(&categories).Error; err != nil {
			return nil, err
		}
	}

	return categories, nil
}