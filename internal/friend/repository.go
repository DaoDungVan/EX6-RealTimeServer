package friend

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	FriendRequests *mongo.Collection
	Friends        *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		FriendRequests: db.Collection("friend_requests"),
		Friends:        db.Collection("friends"),
	}
}

func (r *Repository) SendRequest(
	ctx context.Context, fromID, toID bson.ObjectID) error {
	req := FriendRequest{
		FromUserID: fromID,
		ToUserID:   toID,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	_, err := r.FriendRequests.InsertOne(ctx, req)
	return err
}

func (r *Repository) AcceptRequest(
	ctx context.Context, requestID bson.ObjectID) error {
	_, err := r.FriendRequests.
		UpdateByID(ctx, requestID,
			bson.M{"$set": bson.M{"status": "accepted"}})
	if err != nil {
		return err
	}

	var req FriendRequest
	err = r.FriendRequests.
		FindOne(ctx, bson.M{"_id": requestID}).Decode(&req)
	if err != nil {
		return err
	}

	friend := Friend{
		User1:     req.FromUserID,
		User2:     req.ToUserID,
		CreatedAt: time.Now(),
	}

	_, err = r.Friends.InsertOne(ctx, friend)
	return err
}

func (r *Repository) ListFriends(
	ctx context.Context, userID bson.ObjectID,
) ([]bson.ObjectID, error) {
	cursor, err := r.Friends.Find(ctx, bson.M{
		"$or": []bson.M{
			{"user1": userID},
			{"user2": userID},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []bson.ObjectID
	for cursor.Next(ctx) {
		var f Friend
		if err := cursor.Decode(&f); err == nil {
			if f.User1 == userID {
				ids = append(ids, f.User2)
			} else {
				ids = append(ids, f.User1)
			}
		}
	}
	return ids, nil
}
