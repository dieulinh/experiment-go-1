package main

import (
	"fmt"
	"log"
	"net/http"
	"restapi/internal/db"
	"restapi/internal/handler"
	"restapi/internal/logger"
	"restapi/internal/repository"
	"restapi/internal/router"
	"restapi/internal/service"
)

func main() {
	// init logger
	var myStr = "Hello"
	fmt.Println(len(myStr))
	logger.Init()
	logger.Log.Info("App starting now")

	// init database
	db.Init()

	// init layers
	userRepo := repository.NewRepository(db.DB)
	authService := service.NewAuthService(userRepo)
	productService := service.NewProductService(userRepo)
	cartService := service.NewCartService(userRepo)
	handler := handler.NewHandler(authService, productService, cartService, logger.Log)

	// routes
	mux := router.New(handler)
	// http.HandleFunc("/login", authHandler.Login)
	// http.HandleFunc("/signup", authHandler.SignUp)

	logger.Log.Info("🚀 Server running at :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func TestDB() {
	var result int
	err := db.DB.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		log.Fatal("❌ Test query failed:", err)
	}

	log.Println("✅ DB test query works:", result)
}
