package usecases

import (
	"fmt"
	"time"

	"server/database"
	"server/models"
)

type UserUseCase struct {
	db *database.SpannerDB
}

func NewUserUseCase(db *database.SpannerDB) *UserUseCase {
	return &UserUseCase{db: db}
}

func (uc *UserUseCase) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		ID:        fmt.Sprintf("%d", time.Now().Unix()),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	if err := uc.db.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}