package mongoinit

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Init(ctx context.Context, URI string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))

	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	fmt.Fprintln(os.Stdout, "Successfully Connected to The Database")

	return client, nil
}

func Disconnect(mc *mongo.Client, ctx context.Context) {
	if err := mc.Disconnect(context.Background()); err != nil {
		fmt.Printf("mongo error: %v", err)
		os.Exit(1)
	}
}
