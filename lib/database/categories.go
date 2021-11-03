package database

import (
	"errors"

	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type CategoryModel interface {
	GetCategories() ([]models.Category, error)
	AddCategories(categories models.Category) (models.Category, error)
	EditCategoriesById(id int, newCategories models.Category) (models.Category, error)
	DeleteCategoriesById(id int) error
	GetCategoryById(id int) (models.Category, error)
}

type CategoryDB struct {
	db *gorm.DB
}

func NewCategoryDB(db *gorm.DB) *CategoryDB {
	return &CategoryDB{db: db}
}

func (cdb *CategoryDB) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := cdb.db.Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}

func (cdb *CategoryDB) AddCategories(categories models.Category) (models.Category, error) {
	if err := cdb.db.Save(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}

func (cdb *CategoryDB) EditCategoriesById(id int, newCategories models.Category) (models.Category, error) {
	var category models.Category
	if err := cdb.db.First(&category, id).Error; err != nil {
		return category, err
	}
	category.Name = newCategories.Name
	category.Point = newCategories.Point
	if err := cdb.db.Save(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (cdb *CategoryDB) DeleteCategoriesById(id int) error {
	rows := cdb.db.Delete(&models.Category{}, id).RowsAffected
	if rows == 0 {
		err := errors.New("categories to be deleted is not found")
		return err
	}
	return nil
}

func (cdb *CategoryDB) GetCategoryById(id int) (models.Category, error) {
	var category models.Category
	err := cdb.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}
