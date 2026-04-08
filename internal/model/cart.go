package model

import (
	"time"
)

type Cart struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID        int             `json:"id"`
	CartID    int             `json:"cart_id"`
	ProductID int             `json:"product_id"`
	Product   *Product        `json:"product,omitempty"`
	Quantity  int             `json:"quantity"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type AddToCartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type AddToCartResponse struct {
	CartID   int     `json:"cart_id"`
	ProductID int     `json:"product_id"`
	Quantity int     `json:"quantity"`
	Message  string  `json:"message"`
}
