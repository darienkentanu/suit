package database

import (
	"database/sql"

	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type StaffDB struct {
	db *gorm.DB
	dbSQL *sql.DB
}

func NewStaffDB(db *gorm.DB, dbSQL *sql.DB) *StaffDB {
	return &StaffDB{db: db, dbSQL: dbSQL}
}

type StaffModel interface {
	CreateStaff(staff models.Staff)	(models.Staff, error)
	GetPhoneNumberStaff(number string) (int)
	GetAllStaff() ([]models.ResponseGetStaff, error)
}

func (m *StaffDB) CreateStaff(staff models.Staff) (models.Staff, error) {
	if err := m.db.Save(&staff).Error; err != nil {
		return staff, err
	}

	return staff, nil
}

func (m *StaffDB) GetPhoneNumberStaff(number string) int {
	var staff models.Staff
	row := m.db.Where("phone_number = ?", number).Find(&staff).RowsAffected
	return int(row)
}

func (m *StaffDB) GetAllStaff() ([]models.ResponseGetStaff, error) {
	var allStaff []models.ResponseGetStaff

	rows, err := m.dbSQL.Query("SELECT s.id, s.fullname, l.email, l.username, s.phone_number, l.role, s.drop_point_id FROM staffs s JOIN logins l ON s.id = l.staff_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var staff models.ResponseGetStaff
		if err := rows.Scan(&staff.ID, &staff.Fullname, &staff.Email, &staff.Username, &staff.PhoneNumber, &staff.Role, &staff.DropPointID); err != nil {
			return nil, err
		}
		allStaff = append(allStaff, staff)
	}

	return allStaff, nil
}