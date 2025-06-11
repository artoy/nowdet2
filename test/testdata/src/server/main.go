package main

import (
	"log"
	"net/http"

	"server/database"
	"server/handlers"
	"server/usecases"
)

func main() {
	db, err := database.NewSpannerDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	userUseCase := usecases.NewUserUsecase(db)

	h := handlers.NewHandlers(userUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", h.CreateUser)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
