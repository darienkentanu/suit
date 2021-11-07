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
	GetStaffByID(staffID int) (models.ResponseGetStaff, error)
	UpdateStaff(id int, newStaff models.Staff) (models.Staff, error)
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

	rows, err := m.dbSQL.Query("SELECT s.id, s.fullname, l.email, l.username, s.phone_number, l.role, s.drop_point_id, d.address FROM staffs s JOIN logins l ON s.id = l.staff_id JOIN drop_points d ON s.drop_point_id = d.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var staff models.ResponseGetStaff
		if err := rows.Scan(&staff.ID, &staff.Fullname, &staff.Email, &staff.Username, &staff.PhoneNumber, &staff.Role, &staff.DropPointID, &staff.DropPointAddress); err != nil {
			return nil, err
		}
		allStaff = append(allStaff, staff)
	}

	return allStaff, nil
}

func (m *StaffDB) GetStaffByID(staffID int) (models.ResponseGetStaff, error) {
	var staff models.ResponseGetStaff

	rows, err := m.dbSQL.Query("SELECT s.id, s.fullname, l.email, l.username, s.phone_number, l.role, s.drop_point_id, d.address FROM staffs s JOIN logins l ON s.id = l.staff_id JOIN drop_points d ON s.drop_point_id = d.id WHERE l.staff_id = ?", staffID)
	if err != nil {
		return staff, err
	}
	defer rows.Close()
	
	for rows.Next() {
		if err := rows.Scan(&staff.ID, &staff.Fullname, &staff.Email, &staff.Username, &staff.PhoneNumber, &staff.Role, &staff.DropPointID, &staff.DropPointAddress); err != nil {
			return staff, err
		}
	}

	return staff, nil
}

func (m *StaffDB) UpdateStaff(id int, newStaff models.Staff) (models.Staff, error) {
	var staff models.Staff

	if err := m.db.First(&staff, id).Error; err != nil {
		return staff, err
	}

	staff.Fullname = newStaff.Fullname
	staff.PhoneNumber = newStaff.PhoneNumber
	staff.Drop_PointID = newStaff.Drop_PointID
	
	if err := m.db.Save(&staff).Error; err != nil {
		return staff, err
	}

	return staff, nil
}