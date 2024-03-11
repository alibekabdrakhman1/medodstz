package service

import (
	"context"
	"github.com/alibekabdrakhman1/medodsTZ/internal/model"
)

type ITokenService interface {
	Generate(ctx context.Context, uuid string) (*model.Response, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.Response, error)
}
