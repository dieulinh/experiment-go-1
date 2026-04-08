package service

import (
	"restapi/internal/model"
	"restapi/internal/repository"
)

type CartService struct {
	repo *repository.Repository
}

func NewCartService(repo *repository.Repository) *CartService {
	return &CartService{
		repo: repo,
	}
}

// GetCart retrieves the user's cart
func (cs *CartService) GetCart(userID int) (*model.Cart, error) {
	return cs.repo.GetOrCreateCart(userID)
}

// AddToCart adds a product to the user's cart
func (cs *CartService) AddToCart(userID int, productID int, quantity int) (*model.CartItem, error) {
	// Get or create cart for user
	cart, err := cs.repo.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// Add item to cart
	return cs.repo.AddToCart(cart.ID, productID, quantity)
}

// RemoveFromCart removes an item from the cart
func (cs *CartService) RemoveFromCart(userID int, cartItemID int) error {
	// Get user's cart
	cart, err := cs.repo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	return cs.repo.RemoveFromCart(cart.ID, cartItemID)
}

// ClearCart removes all items from the user's cart
func (cs *CartService) ClearCart(userID int) error {
	// Get user's cart
	cart, err := cs.repo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	return cs.repo.ClearCart(cart.ID)
}
