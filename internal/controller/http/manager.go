package http

import (
	"github.com/alibekabdrakhman1/medodstz/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Manager struct {
	Token ITokenHandler
}

func NewManager(srv *service.Service, logger *zap.SugaredLogger) *Manager {
	return &Manager{NewTokenHandler(srv, logger)}
}

type ITokenHandler interface {
	Generate(c echo.Context) error
	Refresh(c echo.Context) error
}
