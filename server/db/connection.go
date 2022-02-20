package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// closes MongoDB connection and cancels context
func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	// cancels context
	defer cancel()

	defer func() {

		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	return client, ctx, cancel, err

}

func Ping(client *mongo.Client, ctx context.Context) error {

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	fmt.Println("connected successfully")
	return nil

}

func InsertOne(client *mongo.Client, ctx context.Context, database, col string, doc interface{}) (*mongo.InsertOneResult, error) {

	collection := client.Database(database).Collection(col)

	result, err := collection.InsertOne(ctx, doc)
	return result, err

}

func InsertMany(client *mongo.Client, ctx context.Context, database, col string, docs []interface{}) (*mongo.InsertManyResult, error) {

	collection := client.Database(database).Collection(col)

	result, err := collection.InsertMany(ctx, docs)
	return result, err

}

func QueryOne(client *mongo.Client, ctx context.Context, database, col string, query, field interface{}) *mongo.SingleResult {

	collection := client.Database(database).Collection(col)

	result := collection.FindOne(ctx, query)

	// result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return result
}
