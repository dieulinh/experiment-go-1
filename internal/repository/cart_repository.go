package repository

import (
	"database/sql"
	"errors"
	"restapi/internal/model"
)

// GetOrCreateCart retrieves or creates a cart for a user
func (r *Repository) GetOrCreateCart(userID int) (*model.Cart, error) {
	var cartID int

	// Try to get existing cart
	query := `SELECT id FROM carts WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&cartID)

	if err == sql.ErrNoRows {
		// Cart doesn't exist, create it
		insertQuery := `INSERT INTO carts (user_id, created_at, updated_at) 
			VALUES ($1, NOW(), NOW()) RETURNING id`
		err = r.db.QueryRow(insertQuery, userID).Scan(&cartID)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// Fetch the full cart with items
	return r.GetCartByID(cartID)
}

// GetCartByID retrieves a cart with all its items
func (r *Repository) GetCartByID(cartID int) (*model.Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE id = $1`
	
	var cart model.Cart
	err := r.db.QueryRow(query, cartID).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Fetch cart items
	itemsQuery := `
		SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
		       p.id, p.name, p.description, p.price
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = $1
	`

	rows, err := r.db.Query(itemsQuery, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart.Items = []model.CartItem{}
	for rows.Next() {
		var item model.CartItem
		var product model.Product

		err := rows.Scan(
			&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
			&product.ID, &product.Name, &product.Description, &product.Price,
		)
		if err != nil {
			return nil, err
		}

		item.Product = &product
		cart.Items = append(cart.Items, item)
	}

	return &cart, rows.Err()
}

// AddToCart adds a product to the cart
func (r *Repository) AddToCart(cartID int, productID int, quantity int) (*model.CartItem, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	// Check if product already in cart
	checkQuery := `SELECT id, quantity FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	var existingID int
	var existingQty int

	err := r.db.QueryRow(checkQuery, cartID, productID).Scan(&existingID, &existingQty)

	if err == sql.ErrNoRows {
		// Item doesn't exist, insert new
		insertQuery := `
			INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
			RETURNING id, cart_id, product_id, quantity, created_at, updated_at
		`

		var item model.CartItem
		err = r.db.QueryRow(insertQuery, cartID, productID, quantity).Scan(
			&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		return &item, nil
	} else if err != nil {
		return nil, err
	}

	// Item exists, update quantity
	updateQuery := `
		UPDATE cart_items 
		SET quantity = quantity + $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, cart_id, product_id, quantity, created_at, updated_at
	`

	var item model.CartItem
	err = r.db.QueryRow(updateQuery, quantity, existingID).Scan(
		&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// RemoveFromCart removes item from cart
func (r *Repository) RemoveFromCart(cartID int, cartItemID int) error {
	query := `DELETE FROM cart_items WHERE id = $1 AND cart_id = $2`
	result, err := r.db.Exec(query, cartItemID, cartID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("cart item not found")
	}

	return nil
}

// ClearCart removes all items from cart
func (r *Repository) ClearCart(cartID int) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err := r.db.Exec(query, cartID)
	return err
}