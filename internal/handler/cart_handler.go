package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"restapi/internal/model"
)

// AddToCart adds a product to the user's cart
func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	// For now, we'll use a hardcoded user ID. In production, extract from JWT token
	userID := 1 // TODO: Extract from JWT token in request context

	var req model.AddToCartRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf("[AddToCart] JSON decode error: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	h.log.Debugf("[AddToCart] user=%d product=%d quantity=%d", userID, req.ProductID, req.Quantity)

	if req.ProductID <= 0 || req.Quantity <= 0 {
		h.log.Warn("[AddToCart] invalid product_id or quantity")
		http.Error(w, "product_id and quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	// Add to cart via service
	item, err := h.cartService.AddToCart(userID, req.ProductID, req.Quantity)
	if err != nil {
		h.log.Errorf("[AddToCart] service error: %v", err)
		http.Error(w, "failed to add item to cart", http.StatusInternalServerError)
		return
	}

	h.log.Debugf("[AddToCart] success - cart_item_id=%d", item.ID)

	resp := model.AddToCartResponse{
		CartID:    item.CartID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Message:   "Item added to cart successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetCart retrieves the user's cart
func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, use hardcoded user ID. In production, extract from JWT token
	userID := 1 // TODO: Extract from JWT token in request context

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		h.log.Errorf("[GetCart] service error: %v", err)
		http.Error(w, "failed to fetch cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

// RemoveFromCart removes an item from the user's cart
func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := 1 // TODO: Extract from JWT token in request context

	// Get cart item ID from URL query or path parameter
	cartItemIDStr := r.URL.Query().Get("item_id")
	if cartItemIDStr == "" {
		http.Error(w, "item_id query parameter required", http.StatusBadRequest)
		return
	}

	cartItemID, err := strconv.Atoi(cartItemIDStr)
	if err != nil {
		http.Error(w, "invalid item_id", http.StatusBadRequest)
		return
	}

	err = h.cartService.RemoveFromCart(userID, cartItemID)
	if err != nil {
		h.log.Errorf("[RemoveFromCart] service error: %v", err)
		http.Error(w, "failed to remove item from cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Item removed from cart",
	})
}

// ClearCart removes all items from the user's cart
func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := 1 // TODO: Extract from JWT token in request context

	err := h.cartService.ClearCart(userID)
	if err != nil {
		h.log.Errorf("[ClearCart] service error: %v", err)
		http.Error(w, "failed to clear cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Cart cleared successfully",
	})
}
