package mongodb

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	testmg "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	User     string
	Password string
	Database string
}

var (
	mongoConfig = Config{
		User:     "root",
		Password: "password",
		Database: "marketplace",
	}
)

func newMongoContainer(ctx context.Context) (*testmg.MongoDBContainer, error) {
	return testmg.RunContainer(
		ctx,
		testmg.WithUsername(mongoConfig.User),
		testmg.WithPassword(mongoConfig.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Waiting for connections")))
}

func newMongoDB(ctx context.Context, url string) (*mongo.Database, error) {
	opts := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect mongo db: %s", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to mongo postgres db: %s", err)
	}

	return client.Database(mongoConfig.Database), nil
}
