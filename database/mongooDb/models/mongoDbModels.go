package mongooDbModel

import (
	"context"
	"fmt"
	"strconv"

	mongoodb "quick_forge/database/mongooDb"
	"quick_forge/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetAllData(collectionName string) ([]map[string]interface{}, error) {
	db, err := mongoodb.InitDB()
	if err != nil {
		return nil, err
	}
	defer mongoodb.CloseDB(db.Client())

	collection := db.Collection(collectionName)
	cursor, err := collection.Find(context.Background(), bson.M{})
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

func GetData(route utils.Route, arg string) (map[string]interface{}, error) {
	db, err := mongoodb.InitDB()
	if err != nil {
		return nil, err
	}
	defer mongoodb.CloseDB(db.Client())

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

func InsertData(route utils.Route, data map[string]interface{}) (interface{}, error) {
	db, err := mongoodb.InitDB()
	if err != nil {
		return nil, err
	}
	defer mongoodb.CloseDB(db.Client())

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

func UpdateData(route utils.Route, data map[string]interface{}, arg string) error {
	db, err := mongoodb.InitDB()
	if err != nil {
		return err
	}
	defer mongoodb.CloseDB(db.Client())

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

func DeleteData(route utils.Route, arg string) error {
	db, err := mongoodb.InitDB()
	if err != nil {
		return err
	}
	defer mongoodb.CloseDB(db.Client())

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
