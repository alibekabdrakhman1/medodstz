package storage

import (
	"context"
	"github.com/alibekabdrakhman1/medodstz/internal/config"
	"github.com/alibekabdrakhman1/medodstz/internal/storage/mongo"
	"go.uber.org/zap"
)

type Repository struct {
	Token ITokenRepository
}

func NewRepository(ctx context.Context, config *config.Config, logger *zap.SugaredLogger) (*Repository, error) {
	mongoDB, err := mongo.Dial(ctx, &config.Database)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	token := mongo.NewTokenRepository(mongoDB, logger)

	return &Repository{
		Token: token,
	}, nil
}

type ITokenRepository interface {
	CreateToken(ctx context.Context, uuid, token string) error
	UpdateToken(ctx context.Context, uuid, token string) error
	GetToken(ctx context.Context, uuid string) (string, error)
}
