package database

import (
	"context"
	"fmt"
	"time"

	"MysqlToMongo/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(config *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDB.URI))
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("erro ao testar conexão com MongoDB: %v", err)
	}

	return client, nil
}
