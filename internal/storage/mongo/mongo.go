package mongo

import (
	"context"
	"fmt"
	"github.com/alibekabdrakhman1/medodstz/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Dial(ctx context.Context, config *config.Database) (*mongo.Collection, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/", config.User, config.Password, config.Host, config.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, fmt.Errorf("db connecting error: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db ping error: %v", err)
	}
	return client.Database(config.Database).Collection(config.Collection), nil
}
