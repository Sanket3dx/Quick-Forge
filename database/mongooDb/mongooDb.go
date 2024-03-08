package mongoodb

import (
	"context"
	"fmt"
	"quick_forge/utils"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	clientOnce  sync.Once
)

// InitDB initializes the MongoDB client and returns the database instance.
func InitDB() (*mongo.Database, error) {
	clientOnce.Do(func() {
		dbConfig := utils.GetProjectConfig().DB
		uri := dbConfig.MongoString
		clientOptions := options.Client().ApplyURI(uri)

		clientOptions.SetMinPoolSize(5)                    // Minimum number of connections in the pool
		clientOptions.SetMaxPoolSize(20)                   // Maximum number of connections in the pool
		clientOptions.SetMaxConnIdleTime(30 * time.Second) // Maximum idle time before a connection is closed

		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			fmt.Printf("Error connecting to MongoDB: %v\n", err)
			return
		}

		err = client.Ping(context.Background(), nil)
		if err != nil {
			fmt.Printf("Error pinging MongoDB: %v\n", err)
			return
		}

		mongoClient = client
		fmt.Println("Connected to MongoDB!")
	})

	return mongoClient.Database(utils.GetProjectConfig().DB.DatabaseName), nil
}

// CloseDB disconnects the MongoDB client.
func CloseDB() {
	if mongoClient != nil {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			fmt.Printf("Error disconnecting from MongoDB: %v\n", err)
			return
		}
		fmt.Println("Disconnected from MongoDB.")
	}
}

/*
this functions creates new db conn with mongo each time


func InitDB() (*mongo.Database, error) {
	dbConfig := utils.GetProjectConfig().DB
	uri := dbConfig.MongoString
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error pinging MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB!")

	db := client.Database(dbConfig.DatabaseName)

	return db, nil
}

func CloseDB(client *mongo.Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		fmt.Printf("Error disconnecting from MongoDB: %v\n", err)
	} else {
		fmt.Println("Disconnected from MongoDB.")
	}
}
*/
