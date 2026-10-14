package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	Collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Collection: db.Collection("users"),
	}
}

func (r *Repository) Create(
	ctx context.Context, user *User) error {
	user.CreatedAt = time.Now()
	_, err := r.Collection.InsertOne(ctx, user)
	return err
}

func (r *Repository) FindByEmail(
	ctx context.Context, email string) (*User, error) {
	var user User
	err := r.Collection.
		FindOne(ctx, bson.M{"email": email}).
		Decode(&user)

	return &user, err
}
