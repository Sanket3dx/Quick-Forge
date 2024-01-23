package mongoodb

import (
	"context"
	"fmt"
	"quick_forge/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() (*mongo.Database, error) {
	dbConfig := utils.GetProjectConfig().DB
	// mongodb://localhost:27017/
	// uri := fmt.Sprintf("mongodb://%s:%s@%s/%s",
	// 	dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.DatabaseName)
	uri := fmt.Sprintf("mongodb://%s:%s/",
		dbConfig.Host, dbConfig.Port)

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
