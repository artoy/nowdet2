package main

import (
	"log"
	"net/http"

)

func main() {
	db, err := NewSpannerDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	userUseCase := NewUserUsecase(db)

	h := NewHandlers(userUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", h.CreateUser)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
