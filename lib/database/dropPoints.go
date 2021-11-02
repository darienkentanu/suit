package database

import (
	"errors"

	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type DropPointsModel interface {
	GetDropPoints() ([]models.Drop_Point, error)
	AddDropPoints(dropPoints models.Drop_Point) (models.Drop_Point, error)
	EditDropPointsById(id int, newDropPoints models.Drop_Point) (models.Drop_Point, error)
	DeleteDropPointsById(id int) error
}

type DropPointsDB struct {
	db *gorm.DB
}

func NewDropPointsDB(db *gorm.DB) *DropPointsDB {
	return &DropPointsDB{db: db}
}

func (dpdb *DropPointsDB) GetDropPoints() ([]models.Drop_Point, error) {
	var dropPoints []models.Drop_Point
	if err := dpdb.db.Find(&dropPoints).Error; err != nil {
		return dropPoints, err
	}
	return dropPoints, nil
}

func (dpdb *DropPointsDB) AddDropPoints(dropPoints models.Drop_Point) (models.Drop_Point, error) {
	if err := dpdb.db.Save(&dropPoints).Error; err != nil {
		return dropPoints, err
	}
	return dropPoints, nil
}

func (dpdb *DropPointsDB) EditDropPointsById(id int, newDropPoints models.Drop_Point) (models.Drop_Point, error) {
	var dropPoints models.Drop_Point
	if err := dpdb.db.First(&dropPoints, id).Error; err != nil {
		return dropPoints, err
	}
	// dropPoints.ID = newDropPoints.ID
	dropPoints.Address = newDropPoints.Address
	dropPoints.Latitude = newDropPoints.Latitude
	dropPoints.Longitude = newDropPoints.Longitude
	if err := dpdb.db.Save(&dropPoints).Error; err != nil {
		return dropPoints, err
	}
	return dropPoints, nil
}

func (dpdb *DropPointsDB) DeleteDropPointsById(id int) error {
	rows := dpdb.db.Delete(&models.Drop_Point{}, id).RowsAffected
	if rows == 0 {
		err := errors.New("drop points to be deleted is not found")
		return err
	}
	return nil
}

// func (cdb *CategoryDB) GetCategoryId(id int) (models.Category, error) {
// 	var category models.Category
// 	err := cdb.db.Where("id = ?", id).First(&category).Error
// 	if err != nil {
// 		return category, err
// 	}
// 	return category, nil
// }
