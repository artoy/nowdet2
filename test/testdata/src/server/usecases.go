package main

import (
	"fmt"
	"time"

)

type UserUsecase struct {
	db *SpannerDB
}

func NewUserUsecase(db *SpannerDB) *UserUsecase {
	return &UserUsecase{db: db}
}

func (uc *UserUsecase) CreateUser(req CreateUserRequest) (*User, error) {
	user := &User{
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
