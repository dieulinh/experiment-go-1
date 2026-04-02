package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"restapi/internal/service"

	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	service *service.AuthService
	log     *logrus.Logger
	tmpl    *template.Template
}

func NewAuthHandler(s *service.AuthService, log *logrus.Logger) *AuthHandler {
	tmpl := template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))

	return &AuthHandler{
		service: s,
		log:     log,
		tmpl:    tmpl,
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
	Error string
	Email string
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// fix 1: only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignUpRequest // fix 2: use SignUpRequest not LoginRequest — wrong struct
	defer r.Body.Close()  // fix 3: always close body

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("invalid request:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// fix 4: log says "Login attempt" but this is SignUp
	h.log.Info("SignUp attempt:", req.Email)

	// fix 5: basic validation before hitting the service
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.service.SignUp(req.Email, req.Password, req.Name)
	if err != nil {
		h.log.Warn("Create user failed:", err)

		// fix 6: handle specific errors with proper status codes
		// instead of blindly returning err.Error() to the client
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "email already exists", http.StatusConflict) // 409
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError) // 500
		return
	}

	resp := SignUpResponse{
		ID:    user.ID, // fix 7: user.id → user.ID (must be exported field)
		Email: user.Email,
	}

	// fix 8: set 201 Created instead of default 200 OK for resource creation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	h.tmpl.ExecuteTemplate(w, "base", RegisterPageData{})
}
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.tmpl.ExecuteTemplate(w, "base", LoginPageData{})
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("invalid request:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	h.log.Info("Login attempt:", req.Email)

	user, token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		h.log.Warn("login failed:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := LoginResponse{
		Token: token,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
