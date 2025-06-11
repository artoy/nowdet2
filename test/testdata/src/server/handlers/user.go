package handlers

import (
	"encoding/json"
	"net/http"

	"server/models"
	"server/usecases"
)

type Handlers struct {
	userUseCase *usecases.UserUseCase
}

func NewHandlers(userUseCase *usecases.UserUseCase) *Handlers {
	return &Handlers{userUseCase: userUseCase}
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.CreateUser(req)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
