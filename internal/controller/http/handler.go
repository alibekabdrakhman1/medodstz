package http

import (
	"github.com/alibekabdrakhman1/medodstz/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type TokenHandler struct {
	Service *service.Service
	logger  *zap.SugaredLogger
}

func (h *TokenHandler) Generate(c echo.Context) error {
	uuid := c.Param("uuid")
	tokens, err := h.Service.Token.Generate(c.Request().Context(), uuid)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	h.logger.Info(tokens)
	return c.JSON(http.StatusCreated, tokens)
}

func (h *TokenHandler) Refresh(c echo.Context) error {
	var r struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := c.Bind(&r)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	tokens, err := h.Service.Token.RefreshToken(c.Request().Context(), r.RefreshToken)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	h.logger.Info(tokens)
	return c.JSON(http.StatusOK, tokens)
}

func NewTokenHandler(s *service.Service, logger *zap.SugaredLogger) *TokenHandler {
	return &TokenHandler{
		Service: s,
		logger:  logger,
	}
}
