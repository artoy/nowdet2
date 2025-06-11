package main

import (
	"log"
	"net/http"

	"server/database"
	"server/handlers"
	"server/usecases"
)

func main() {
	// Initialize fake Spanner database
	db, err := database.NewSpannerDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize use cases
	userUseCase := usecases.NewUserUseCase(db)

	// Initialize handlers with use cases
	h := handlers.NewHandlers(userUseCase)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", h.CreateUser)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
