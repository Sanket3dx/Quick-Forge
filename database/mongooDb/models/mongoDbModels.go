package mongooDbModel

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	mongoodb "quick_forge/database/mongooDb"
	"quick_forge/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Counter struct {
	ID    string `bson:"_id"`
	Value int    `bson:"value"`
}

func getNextSequence(db *mongo.Database, collectionName string) (int, error) {

	collection := db.Collection("counters")

	result := collection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": collectionName},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetUpsert(true),
	)

	var counter Counter
	if err := result.Decode(&counter); err != nil {
		return 0, err
	}

	return counter.Value, nil
}

func createCollectionIfNotExists(db *mongo.Database, collectionName string) error {
	collection := db.Collection(collectionName)
	err := collection.Database().Client().Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}

	if _, err := collection.InsertOne(context.Background(), bson.M{}); err != nil {
		return err
	}

	return nil
}

func GetAllData(dbConfig utils.Database, collectionName string, filters []map[string]interface{}) ([]map[string]interface{}, error) {
	db, err := mongoodb.InitDB(dbConfig)
	if err != nil {
		return nil, err
	}
	collection := db.Collection(collectionName)
	fmt.Println(collection)
	// Construct the filter based on the provided filters
	filter := bson.M{}
	for _, f := range filters {
		for k, v := range f {
			filter[k] = v
		}
	}

	var start, limit int

	if startStr, ok := filter["start"].(string); ok {
		if startVal, err := strconv.Atoi(startStr); err == nil {
			start = startVal
		}
	}
	if limitStr, ok := filter["limit"].(string); ok {
		if limitVal, err := strconv.Atoi(limitStr); err == nil {
			limit = limitVal
		}
	}
	delete(filter, "limit")
	delete(filter, "start")

	options := options.Find()
	if start != 0 && limit != 0 {
		// Set up options for Find method
		options = options.SetSkip(int64(start)).SetLimit(int64(limit))
	}

	cursor, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []map[string]interface{}

	for cursor.Next(context.Background()) {
		var result map[string]interface{}

		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func GetData(dbConfig utils.Database, route utils.Route, arg string) (map[string]interface{}, error) {
	db, err := mongoodb.InitDB(dbConfig)
	if err != nil {
		return nil, err
	}

	collection := db.Collection(route.DBTableName)

	primaryKeyField, _ := route.DBTableStruct["primary_key"].(string)
	argInt, err := strconv.Atoi(arg)
	if err != nil {
		return nil, err
	}
	filter := bson.M{primaryKeyField: argInt}
	result := collection.FindOne(context.Background(), filter)
	if err := result.Err(); err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := result.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func InsertData(dbConfig utils.Database, route utils.Route, data map[string]interface{}) (interface{}, error) {
	db, err := mongoodb.InitDB(dbConfig)
	if err != nil {
		return nil, err
	}

	if err := createCollectionIfNotExists(db, route.DBTableName); err != nil {
		return nil, fmt.Errorf("error creating collection: %v", err)
	}

	collection := db.Collection(route.DBTableName)

	nextValue, err := getNextSequence(db, route.DBTableName)
	if err != nil {
		return nil, err
	}
	primaryKeyField, _ := route.DBTableStruct["primary_key"].(string)
	data[primaryKeyField] = nextValue
	result, err := collection.InsertOne(context.Background(), data)
	if err != nil && result.InsertedID == "" {
		return nil, err
	}

	return nextValue, nil
}

func UpdateData(dbConfig utils.Database, route utils.Route, data map[string]interface{}, arg string) error {
	db, err := mongoodb.InitDB(dbConfig)
	if err != nil {
		return err
	}

	collection := db.Collection(route.DBTableName)
	primaryKeyField, _ := route.DBTableStruct["primary_key"].(string)
	// Update a document in the collection
	argInt, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}
	filter := bson.M{primaryKeyField: argInt}
	update := bson.M{"$set": data}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("no documents were updated")
	}

	return nil
}

func DeleteData(dbConfig utils.Database, route utils.Route, arg string) error {
	db, err := mongoodb.InitDB(dbConfig)
	if err != nil {
		return err
	}

	collection := db.Collection(route.DBTableName)
	primaryKeyField, _ := route.DBTableStruct["primary_key"].(string)
	argInt, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}
	filter := bson.M{primaryKeyField: argInt}
	result, err := collection.DeleteMany(context.Background(), filter)

	fmt.Println(result)
	if err != nil {
		fmt.Println("Error deleting document:", err)
		return err
	}
	fmt.Println("DeletedCount:", result.DeletedCount)

	if result.DeletedCount == 0 {
		return fmt.Errorf("no documents were deleted")
	}

	return nil
}

func AuthenticateUser(username, password string) (bson.M, error) {
	db, err := mongoodb.InitDBInternal()
	if err != nil {
		return nil, err
	}
	defer mongoodb.CloseDB()

	collection := db.Collection("user_logins")

	filter := bson.M{"username": username, "password": password}
	result := collection.FindOne(context.Background(), filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // User not found
		}
		return nil, err // Other error occurred
	}

	var user bson.M
	if err := result.Decode(&user); err != nil {
		return nil, err // Error decoding user
	}

	return user, nil
}
