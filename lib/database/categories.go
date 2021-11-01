package database

import (
	"errors"

	"github.com/darienkentanu/suit/models"
)

func GetCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}

func AddCategories(categories models.Category) (models.Category, error) {
	if err := db.Save(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}

func EditCategoriesById(id int, newCategories models.Category) (models.Category, error) {
	var category models.Category
	if err := db.First(&category, id).Error; err != nil {
		return category, err
	}
	category.Name = newCategories.Name
	category.Point = newCategories.Point
	if err := db.Save(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func DeleteCategoriesById(id int) error {
	rows := db.Delete(&models.Category{}, id).RowsAffected
	if rows == 0 {
		err := errors.New("categories to be deleted is not found")
		return err
	}
	return nil
}

// func GetCategoryId(id int) (models.Category, error) {
// 	var category models.Category
// 	err := db.Where("id = ?", id).First(&category).Error
// 	if err != nil {
// 		return category, err
// 	}
// 	return category, nil
// }
