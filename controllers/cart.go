package controllers

import "github.com/darienkentanu/suit/models"

type CartDB interface {
	CreateCart(cart models.Cart) (error)
}

type CartController struct {
	db CartDB
}