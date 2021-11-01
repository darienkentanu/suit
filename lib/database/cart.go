package database

import (
	"github.com/darienkentanu/suit/config"
	"github.com/darienkentanu/suit/models"
)

func CreateCart(cart models.Cart) error {
	if err := config.InitDB().Save(&cart).Error; err != nil {
		return err
	}
	return nil
}