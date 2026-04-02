package router

import (
	"net/http"
	"restapi/internal/handler"
)

func New(authHandler *handler.AuthHandler) http.Handler {
	mux := http.NewServeMux()

	// serve static files (css, js, images)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("GET /login", authHandler.LoginPage)
	mux.HandleFunc("GET /signup", authHandler.RegisterPage)
	mux.HandleFunc("POST /signup", authHandler.SignUp)
	mux.HandleFunc("POST /login", authHandler.Login)

	return mux
}
