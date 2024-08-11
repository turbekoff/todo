package mongo

import (
	"context"
	"time"

	"github.com/turbekoff/todo/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewConnection(config *config.MongoConfig) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(config.URI)

	if config.User != "" && config.Password != "" {
		opts.SetAuth(options.Credential{
			Username: config.User,
			Password: config.Password,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
