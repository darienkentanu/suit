package database

import (
	"errors"

	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type VoucherModel interface {
	GetVouchers() ([]models.Voucher, error)
	AddVouchers(voucher models.Voucher) (models.Voucher, error)
	EditVouchersById(id int, newVoucher models.Voucher) (models.Voucher, error)
	DeleteVouchersById(id int) error
	GetVoucherByID(voucherID int) (models.Voucher, error)
}

type VoucherDB struct {
	db *gorm.DB
}

func NewVoucherDB(db *gorm.DB) *VoucherDB {
	return &VoucherDB{db: db}
}

func (vdb *VoucherDB) GetVouchers() ([]models.Voucher, error) {
	var vouchers []models.Voucher
	if err := vdb.db.Find(&vouchers).Error; err != nil {
		return vouchers, err
	}
	return vouchers, nil
}

func (vdb *VoucherDB) AddVouchers(vouchers models.Voucher) (models.Voucher, error) {
	if err := vdb.db.Save(&vouchers).Error; err != nil {
		return vouchers, err
	}
	return vouchers, nil
}

func (vdb *VoucherDB) EditVouchersById(id int, newVouchers models.Voucher) (models.Voucher, error) {
	var voucher models.Voucher
	if err := vdb.db.First(&voucher, id).Error; err != nil {
		return voucher, err
	}
	voucher.Name = newVouchers.Name
	voucher.Point = newVouchers.Point
	if err := vdb.db.Save(&voucher).Error; err != nil {
		return voucher, err
	}
	return voucher, nil
}

func (vdb *VoucherDB) DeleteVouchersById(id int) error {
	rows := vdb.db.Delete(&models.Voucher{}, id).RowsAffected
	if rows == 0 {
		err := errors.New("voucher is not found")
		return err
	}
	return nil
}

func (vdb *VoucherDB) GetVoucherByID(voucherID int) (models.Voucher, error) {
	var voucher models.Voucher
	if err := vdb.db.First(&voucher, voucherID).Error; err != nil {
		return voucher, err
	}

	return voucher, nil
}