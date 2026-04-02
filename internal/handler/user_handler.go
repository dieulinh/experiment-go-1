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
	service      *service.AuthService
	log          *logrus.Logger
	loginTmpl    *template.Template
	registerTmpl *template.Template
}

func NewAuthHandler(s *service.AuthService, log *logrus.Logger) *AuthHandler {
	base := filepath.Join("templates", "base.html")
	loginTmpl := template.Must(template.ParseFiles(base, filepath.Join("templates", "login.html")))
	registerTmpl := template.Must(template.ParseFiles(base, filepath.Join("templates", "signup.html")))

	return &AuthHandler{
		service:      s,
		log:          log,
		loginTmpl:    loginTmpl,
		registerTmpl: registerTmpl,
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

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
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
func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	h.registerTmpl.ExecuteTemplate(w, "base", RegisterPageData{})
}
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.loginTmpl.ExecuteTemplate(w, "base", LoginPageData{})
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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
