package dbtest

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StartTestContainer(ctx context.Context) (*mongo.Client, string, func()) {
	container, err := mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:4.4.8"))
	if err != nil {
		log.Panicf("failed to start test container MongoDB: %s", err)
	}

	endpoint, err := container.ConnectionString(ctx)
	if err != nil {
		log.Panicf("failed to get connection string for test container MongoDB: %s", err)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		log.Panicf("failed to connect to test container MongoDB: %s", err)
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Panicf("failed to ping to test container MongoDB: %s", err)
	}

	containerIP, err := container.Host(ctx)
	if err != nil {
		log.Panicf("failed to get test container IP address: %s", err)
	}

	stopContainer := func() {
		if err := container.Terminate(ctx); err != nil {
			log.Panicf("failed to stop test container MongoDB: %s", err)
		}
	}

	return mongoClient, containerIP, stopContainer
}

func InsertMockData(ctx context.Context, collection *mongo.Collection, data []interface{}) error {
	_, err := collection.InsertMany(ctx, data)
	return err
}

func DeleteMockData(ctx context.Context, collection *mongo.Collection) error {
	_, err := collection.DeleteMany(ctx, bson.M{})
	return err
}

func FindMockData(ctx context.Context, collection *mongo.Collection, filter interface{}) ([]interface{}, error) {
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []interface{}
	for cur.Next(ctx) {
		var result interface{}
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
