package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"restapi/internal/helper"
	"restapi/internal/model"
	"restapi/internal/service"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type Product struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description"`
}
type ProductsResponse struct {
	Status    string          `json:"status"`
	Timestamp string          `json:"timestamp"`
	Products  []model.Product `json:"products"`
}
type Handler struct {
	service        *service.AuthService
	productService *service.ProductService
	cartService    *service.CartService
	contactService *service.ContactService
	log            *logrus.Logger
}

func NewHandler(s *service.AuthService, ps *service.ProductService, cs *service.CartService, contactService *service.ContactService, log *logrus.Logger) *Handler {

	return &Handler{
		service:        s,
		productService: ps,
		cartService:    cs,
		contactService: contactService,
		log:            log,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
}
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
type SignUpResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
type RegisterPageData struct {
	Error    string
	Email    string
	Password string
	Name     string
}
type LoginPageData struct {
	Error    string
	Email    string
	Password string
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	// DEBUG: log content type so we know how data arrived
	h.log.Debugf("[SignUp] Content-Type: %s", r.Header.Get("Content-Type"))

	var req SignUpRequest

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Errorf("[SignUp] JSON decode error: %v", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
	} else {
		// HTML form submission (application/x-www-form-urlencoded)
		if err := r.ParseForm(); err != nil {
			h.log.Errorf("[SignUp] ParseForm error: %v", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
		req.Name = r.FormValue("name")
	}

	// DEBUG: log parsed values (mask password in real apps)
	h.log.Debugf("[SignUp] parsed — email=%q name=%q password_len=%d", req.Email, req.Name, len(req.Password))

	if req.Email == "" || req.Password == "" {
		h.log.Warn("[SignUp] missing email or password")
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	h.log.Infof("[SignUp] attempting to create user: %s", req.Email)

	user, err := h.service.SignUp(req.Email, req.Password, req.Name)
	if err != nil {
		h.log.Warnf("[SignUp] service error: %v", err)

		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Debugf("[SignUp] user created — id=%d email=%s", user.ID, user.Email)

	resp := SignUpResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	// fix 8: set 201 Created instead of default 200 OK for resource creation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
func (h *Handler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{}
	helper.Render(w, data, "templates/signup.html")
}
func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{}
	helper.Render(w, data, "templates/login.html")
}
func (h *Handler) ProductList(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAll()
	log.Print(products)
	if err != nil {
		h.log.Errorf("[ProductListJson] error fetching products: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// data := map[string]any{
	// 	"products": products,
	// }
	resp := ProductsResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Products:  products,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)

}

func (h *Handler) ProductListPage(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAll()
	log.Print(products)
	if err != nil {
		h.log.Errorf("[ProductListPage] error fetching products: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"products": products,
	}
	helper.Render(w, data, "templates/products/index.html")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	h.log.Debugf("[Login] Content-Type: %s", r.Header.Get("Content-Type"))

	var req LoginRequest

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Errorf("[Login] JSON decode error: %v", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			h.log.Errorf("[Login] ParseForm error: %v", err)
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
	}

	h.log.Debugf("[Login] parsed — email=%q", req.Email)

	if req.Email == "" || req.Password == "" {
		h.log.Warn("[Login] missing email or password")
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	h.log.Infof("[Login] attempting login for: %s", req.Email)

	user, token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		h.log.Warnf("[Login] failed: %v", err)
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	h.log.Debugf("[Login] success for: %s", user.Email)

	resp := LoginResponse{
		Token: token,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
