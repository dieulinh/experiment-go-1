package router

import (
	"net/http"
	"restapi/internal/handler"
)

func New(authHandler *handler.Handler) http.Handler {
	mux := http.NewServeMux()

	// serve static files (css, js, images)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("GET /login", authHandler.LoginPage)
	mux.HandleFunc("GET /signup", authHandler.RegisterPage)
	mux.HandleFunc("POST /signup", authHandler.SignUp)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("GET /products", authHandler.ProductListPage)
	mux.HandleFunc("GET /api/health", authHandler.APIHealth)
	mux.HandleFunc("GET /product_list", authHandler.ProductList)
	return mux
}
