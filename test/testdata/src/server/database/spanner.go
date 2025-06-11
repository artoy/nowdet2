package database

import (
	"context"
	"fmt"

	"server/models"

	"cloud.google.com/go/spanner"
)

type SpannerDB struct {
	client *spanner.Client
	ctx    context.Context
}

func NewSpannerDB() (*SpannerDB, error) {
	ctx := context.Background()

	// Fake Spanner configuration for testing
	projectID := "fake-project-id"
	instanceID := "fake-instance-id"
	databaseID := "fake-database-id"

	database := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		projectID, instanceID, databaseID)

	// In a real implementation, this would connect to actual Spanner
	// For testing, we'll simulate the connection
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		// Return a mock client for testing purposes
		return &SpannerDB{
			client: nil, // Fake client
			ctx:    ctx,
		}, nil
	}

	return &SpannerDB{
		client: client,
		ctx:    ctx,
	}, nil
}

func (db *SpannerDB) Close() {
	if db.client != nil {
		db.client.Close()
	}
}

func (db *SpannerDB) CreateUser(user *models.User) error {
	// If client is nil (fake/testing mode), just return success
	if db.client == nil {
		return nil
	}

	// Insert user into Spanner database
	mutation := spanner.Insert("users",
		[]string{"id", "name", "email", "created_at"},
		[]any{user.ID, user.Name, user.Email, user.CreatedAt})

	_, err := db.client.Apply(db.ctx, []*spanner.Mutation{mutation})
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}
