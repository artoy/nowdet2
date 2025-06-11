package usecases

import (
	"fmt"
	"time"

	"server/database"
	"server/models"
)

type UserUsecase struct {
	db *database.SpannerDB
}

func NewUserUsecase(db *database.SpannerDB) *UserUsecase {
	return &UserUsecase{db: db}
}

func (uc *UserUsecase) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		ID:        fmt.Sprintf("%d", time.Now().Unix()),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.db.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
