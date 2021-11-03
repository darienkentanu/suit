package database

import (
	"errors"

	"github.com/darienkentanu/suit/models"
	"gorm.io/gorm"
)

type CartDB struct {
	db *gorm.DB
}

func NewCartDB(db *gorm.DB) *CartDB {
	return &CartDB{db: db}
}

type CartModel interface {
	CreateCart(cart models.Cart) error
	AddToCart(userID int, input models.CartItem_Input) (models.CartItem, error)
	GetCartItem(userID int) ([]models.CartItem, error)
	EditCartItem(cartID int, input models.CartItem_Input) (models.CartItem, error)
	DeleteCartItem(cardID int) error
	CheckCartByCategoryID(int, int) bool
	AddCartWeight(int, models.CartItem_Input) (models.CartItem, error)
}

func (cdb *CartDB) CreateCart(cart models.Cart) error {
	if err := cdb.db.Save(&cart).Error; err != nil {
		return err
	}
	return nil
}

func (cdb *CartDB) AddToCart(userID int, input models.CartItem_Input) (models.CartItem, error) {
	var cartItems models.CartItem
	cartItems.CartUserID = userID
	cartItems.CategoryID = input.CategoryID
	cartItems.Weight = input.Weight
	if err := cdb.db.Select("cart_user_id", "category_id", "weight").Create(&cartItems).Error; err != nil {
		return cartItems, err
	}
	return cartItems, nil
}

func (cdb *CartDB) GetCartItem(userID int) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	if err := cdb.db.Find(&cartItems, "cart_user_id = ? and checkout_id IS NULL", userID).Error; err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (cdb *CartDB) EditCartItem(cartID int, input models.CartItem_Input) (models.CartItem, error) {
	var cartItems models.CartItem
	if err := cdb.db.First(&cartItems, cartID).Update("weight", input.Weight).Error; err != nil {
		return cartItems, err
	}
	return cartItems, nil
}

func (cdb *CartDB) DeleteCartItem(cardID int) error {
	rows := cdb.db.Delete(&models.CartItem{}, cardID).RowsAffected
	if rows == 0 {
		err := errors.New("cart item to be deleted is not found")
		return err
	}
	return nil
}

func (cdb *CartDB) CheckCartByCategoryID(Userid int, categoryID int) bool {
	var cartItem models.CartItem
	row := cdb.db.First(&cartItem, "cart_user_id = ?", Userid, "category_id = ?", categoryID, "checkout_id = null").RowsAffected
	if row == 1 {
		return true
	}
	return false
}

func (cdb *CartDB) AddCartWeight(userID int, input models.CartItem_Input) (models.CartItem, error) {
	var cartItem models.CartItem
	if err := cdb.db.Where("cart_user_id = ? and category_id = ? and checkout_id IS NULL", userID, input.CategoryID).First(&cartItem).Error; err != nil {
		return cartItem, err
	}
	updatedWeight := input.Weight + cartItem.Weight
	if err := cdb.db.Model(&cartItem).Update("weight", updatedWeight).Error; err != nil {
		return cartItem, err
	}
	return cartItem, nil
}
