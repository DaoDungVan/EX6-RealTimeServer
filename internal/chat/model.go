package chat

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Message struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SenderID  bson.ObjectID `bson:"sender_id" json:"sender_id"`
	RoomID    string        `bson:"room_id"   json:"room_id"`
	Content   string        `bson:"content"   json:"content"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
