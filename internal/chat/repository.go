package chat

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	Collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Collection: db.Collection("messages"),
	}
}

func (r *Repository) SaveMessage(ctx context.Context, m *Message) error {
	m.CreatedAt = time.Now()
	_, err := r.Collection.InsertOne(ctx, m)
	return err
}
