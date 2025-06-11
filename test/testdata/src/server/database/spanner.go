package database

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"

	"server/models"
)

type SpannerDB struct {
	client *spanner.Client
	ctx    context.Context
}

func NewSpannerDB() (*SpannerDB, error) {
	// This is a fake implementation for testing purposes
	return &SpannerDB{}, nil
}

func (db *SpannerDB) Close() {
	if db.client != nil {
		db.client.Close()
	}
}

func (db *SpannerDB) CreateUser(user *models.User) error {
	mutation := spanner.Insert("users",
		[]string{"id", "name", "email", "created_at", "updated_at"},
		[]any{user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt})

	_, err := db.client.Apply(db.ctx, []*spanner.Mutation{mutation})
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}
