package db

import (
	"context"
	"fmt"
	"realtime-server/internal/common/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client

func InitMongo() *mongo.Database {
	ctx, cancel := context.
		WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.
		Connect(options.Client().
			ApplyURI(config.
				GetEnv("MONGO_URI", "")),
		)
	if err != nil {
		panic(fmt.Errorf("MongoDB connection error: %v", err))
	}
	if err := client.Ping(ctx, nil); err != nil {
		panic(fmt.Errorf("MongoDB ping failed: %v", err))
	}

	MongoClient = client
	fmt.Println("Connected to MongoDB successfully!")
	return client.Database("realtime_db")
}
