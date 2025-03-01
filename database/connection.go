package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"wgnalvian.com/payment-server/config"
)

var client *mongo.Client

func ConnectMongo() *mongo.Database {

	if client == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		conn, _ := mongo.Connect(ctx, options.Client().ApplyURI(config.LoadConfig().MONGO_URL))
		defer cancel()

		err := conn.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

		client = conn
	}

	db := client.Database(config.LoadConfig().DATABASE_NAME)
	return db
}

func DisconnectMongo() {
	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to disconnect MongoDB: %v", err)
		}
		log.Println("MongoDB connection closed.")
	}
}
