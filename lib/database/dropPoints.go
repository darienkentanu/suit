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
	GetDropPointsByID(id int) (models.Drop_Point, error)
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
	} else if len(dropPoints) == 0 {
		err := errors.New("is empty")
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
	} else if dropPoints.ID == 0 {
		err := errors.New("not found")
		return dropPoints, err
	}
	return dropPoints, nil
}

func (dpdb *DropPointsDB) DeleteDropPointsById(id int) error {
	var dropPoint models.Drop_Point

	row := dpdb.db.Delete(&dropPoint, id).RowsAffected
	if row == 0 {
		err := errors.New("not found")
		return err
	}

	return nil
}

func (dpdb *DropPointsDB) GetDropPointsByID(id int) (models.Drop_Point, error) {
	var dropPoint models.Drop_Point
	err := dpdb.db.Where("id = ?", id).First(&dropPoint).Error
	if err != nil {
		return dropPoint, err
	}
	return dropPoint, nil
}
