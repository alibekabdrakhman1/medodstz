package service

import (
	"github.com/alibekabdrakhman1/medodstz/internal/storage"
	"go.uber.org/zap"
)

type Service struct {
	Token ITokenService
}

func NewManager(repository *storage.Repository, jwtKey string, logger *zap.SugaredLogger) *Service {
	return &Service{
		Token: NewTokenService(repository, jwtKey, logger),
	}
}
