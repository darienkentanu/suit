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
	GetTransactionsByUserID(userID int) ([]models.Transaction, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateStatusTransaction(transactionID int) (models.Transaction, error)
	GetTransationsRangeDate(id int, role, rangeDate string) ([]models.TransactionSQL, error)
	GetTransationTotalRangeDate(id int, role, rangeDate string) ([]models.ResponseCategoryReport, error)
	GetTransationTotal(id int, role string) ([]models.ResponseCategoryReport, error)
	GetTransactionsByDropPointID(dropPointID int) ([]models.Transaction, error)
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

	if err := m.db.Find(&transactions, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (m *TransactionDB) GetTransactionsByDropPointID(dropPointID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	if err := m.db.Find(&transactions, "drop_point_id = ?", dropPointID).Error; err != nil {
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
	var rows *sql.Rows
	var err error
	
	if role == "staff" {
		if rangeDate == "daily" {
			rows, err = m.dbSQL.Query("SELECT * FROM transactions WHERE created_at >= DATE_ADD(CURDATE(), INTERVAL -1 DAY)")
		} else if rangeDate == "weekly" {
			rows, err = m.dbSQL.Query("SELECT * FROM transactions WHERE created_at >= DATE_ADD(CURDATE(), INTERVAL -7 DAY)")
		} else if rangeDate == "monthly" {
			rows, err = m.dbSQL.Query("SELECT * FROM transactions WHERE created_at >= DATE_ADD(CURDATE(), INTERVAL -30 DAY)")
		}
	} else if role == "user" {
		if rangeDate == "daily" {
			rows, err = m.dbSQL.Query("SELECT * FROM transactions WHERE created_at >= DATE_ADD(CURDATE(), INTERVAL -1 DAY) AND user_id = ?", id)
		} else if rangeDate == "weekly" {
			rows, err = m.dbSQL.Query("SELECT * FROM transactions WHERE created_at >= DATE_ADD(CURDATE(), INTERVAL -7 DAY) AND user_id = ?", id)
		} else if rangeDate == "monthly" {
			rows, err = m.dbSQL.Query("SELECT * FROM transactions WHERE created_at >= DATE_ADD(CURDATE(), INTERVAL -30 DAY) AND user_id = ?", id)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var trans models.TransactionSQL
		if err := rows.Scan(&trans.ID, &trans.UserID, &trans.Status,
			&trans.Point, &trans.Method, &trans.Drop_PointID, &trans.CheckoutID, &trans.CreatedAt, &trans.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, trans)
	}

	return transactions, nil
}

func (m *TransactionDB) GetTransationTotalRangeDate(id int, role, rangeDate string) ([]models.ResponseCategoryReport, error) {
	var categories []models.ResponseCategoryReport
	var rows *sql.Rows
	var err error
	
	if role == "staff" {
		if rangeDate == "daily" {
			rows, err = m.dbSQL.Query("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -1 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 GROUP BY ci.category_id ORDER BY ci.category_id")
		} else if rangeDate == "weekly" {
			rows, err = m.dbSQL.Query("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -7 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 GROUP BY ci.category_id ORDER BY ci.category_id")
		} else if rangeDate == "monthly" {
			rows, err = m.dbSQL.Query("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE t.created_at >= DATE_ADD(CURDATE(), INTERVAL -30 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 GROUP BY ci.category_id ORDER BY ci.category_id")
		}
	} else if role == "user" {
		if rangeDate == "daily" {
			rows, err = m.dbSQL.Query("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE ci.created_at >= DATE_ADD(CURDATE(), INTERVAL -1 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 AND ci.cart_user_id = ? GROUP BY ci.category_id ORDER BY ci.category_id", id)
		} else if rangeDate == "weekly" {
			rows, err = m.dbSQL.Query("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE ci.created_at >= DATE_ADD(CURDATE(), INTERVAL -7 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 AND ci.cart_user_id = ? GROUP BY ci.category_id ORDER BY ci.category_id", id)
		} else if rangeDate == "monthly" {
			rows, err = m.dbSQL.Query("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE ci.created_at >= DATE_ADD(CURDATE(), INTERVAL -30 DAY) AND ci.checkout_id IS NOT NULL AND t.status = 1 AND ci.cart_user_id = ? GROUP BY ci.category_id ORDER BY ci.category_id", id)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var category models.ResponseCategoryReport
		if err := rows.Scan(&category.ID, &category.Name, &category.Weight); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (m *TransactionDB) GetTransationTotal(id int, role string) ([]models.ResponseCategoryReport, error) {
	var categories []models.ResponseCategoryReport
	var rows *sql.Rows
	var err error
	
	if role == "staff" {
		rows, err = m.dbSQL.Query("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE ci.checkout_id IS NOT NULL AND t.status = 1 GROUP BY ci.category_id ORDER BY ci.category_id")
	} else if role == "user" {
		rows, err = m.dbSQL.Query("SELECT ci.category_id, c.name, SUM(ci.weight) FROM cart_items ci JOIN categories c ON ci.category_id = c.id JOIN transactions t ON ci.checkout_id = t.checkout_id WHERE ci.checkout_id IS NOT NULL AND t.status = 1 AND ci.cart_user_id = ? GROUP BY ci.category_id ORDER BY ci.category_id", id)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var category models.ResponseCategoryReport
		if err := rows.Scan(&category.ID, &category.Name, &category.Weight); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}