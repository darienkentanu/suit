package database

import (
	"database/sql"

	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type UserVoucherDB struct {
	db *gorm.DB
	dbSQL *sql.DB
}

func NewUserVoucherDB(db *gorm.DB, dbSQL *sql.DB) *UserVoucherDB {
	return &UserVoucherDB{db: db, dbSQL: dbSQL}
}

type UserVoucherModel interface {
	AddUserVoucher(userID, voucherID int) (models.User_Voucher, error)
	GetAllVoucher(userID int) ([]models.ResponseGetUserVoucher, error)
	UpdateStatusUserVoucher(userID, voucherID int) (models.User_Voucher, error)
	GetUserVoucherByID(id int) (models.ResponseGetUserVoucher, error)
	GetVoucherByUserAndVoucherID(userID, voucherID int) int
}

func (m *UserVoucherDB) AddUserVoucher(userID, voucherID int) (models.User_Voucher, error) {
	var userVoucher models.User_Voucher
	userVoucher.UserID = userID
	userVoucher.VoucherID = voucherID
	userVoucher.Status = "available"

	if err := m.db.Save(&userVoucher).Error; err != nil {
		return userVoucher, err
	}

	return userVoucher, nil
}

func (m *UserVoucherDB) GetAllVoucher(userID int) ([]models.ResponseGetUserVoucher, error) {
	var userVouchers []models.ResponseGetUserVoucher

	rows, err := m.dbSQL.Query("SELECT uv.id, uv.user_id, uv.voucher_id, v.name, v.point, uv.status FROM user_vouchers uv JOIN vouchers v ON uv.voucher_id = v.id WHERE user_id = ? ORDER BY uv.status DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var userVoucher models.ResponseGetUserVoucher
		if err := rows.Scan(&userVoucher.ID, &userVoucher.UserID, &userVoucher.VoucherID, &userVoucher.VoucherName, &userVoucher.Point, &userVoucher.Status); err != nil {
			return nil, err
		}
		userVouchers = append(userVouchers, userVoucher)
	}

	return userVouchers, nil
}

func (m *UserVoucherDB) UpdateStatusUserVoucher(userID, voucherID int) (models.User_Voucher, error) {
	var userVoucher models.User_Voucher

	status := "available"

	if err := m.db.Where("user_id = ? and voucher_id = ? and status = ?", userID, voucherID, status).First(&userVoucher).Error; err != nil {
		return userVoucher, err
	}

	newStatus := "used"

	if err := m.db.Model(&userVoucher).Update("status", newStatus).Error; err != nil {
		return userVoucher, err
	}

	return userVoucher, nil
}

func (m *UserVoucherDB) GetUserVoucherByID(userVoucherID int) (models.ResponseGetUserVoucher, error) {
	var userVoucher models.ResponseGetUserVoucher

	rows, err := m.dbSQL.Query("SELECT uv.id, uv.user_id, uv.voucher_id, v.name, v.point, uv.status FROM user_vouchers uv JOIN vouchers v ON uv.voucher_id = v.id WHERE uv.id = ?", userVoucherID)
	if err != nil {
		return userVoucher, err
	}
	defer rows.Close()
	
	for rows.Next() {
		if err := rows.Scan(&userVoucher.ID, &userVoucher.UserID, &userVoucher.VoucherID, &userVoucher.VoucherName, &userVoucher.Point, &userVoucher.Status); err != nil {
			return userVoucher, err
		}
	}

	return userVoucher, nil
}

func (m *UserVoucherDB) GetVoucherByUserAndVoucherID(userID, voucherID int) int {
	var userVoucher models.User_Voucher

	status := "available"

	row := m.db.Where("user_id = ? and voucher_id = ? and status = ?", userID, voucherID, status).Find(&userVoucher).RowsAffected
	
	return int(row)
}