package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/alibekabdrakhman1/medodsTZ/internal/model"
	"github.com/alibekabdrakhman1/medodsTZ/internal/storage"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"time"
)

type TokenService struct {
	repository   *storage.Repository
	jwtSecretKey string
	logger       *zap.SugaredLogger
}

func NewTokenService(repo *storage.Repository, jwtKey string, logger *zap.SugaredLogger) *TokenService {
	return &TokenService{
		repository:   repo,
		jwtSecretKey: jwtKey,
		logger:       logger,
	}
}

func (s *TokenService) Generate(ctx context.Context, uuid string) (*model.Response, error) {
	tokens, err := s.generateToken(ctx, uuid)
	if err != nil {
		s.logger.Errorf("generating token err: %v", err)
		return nil, fmt.Errorf("generating token err: %w", err)
	}
	err = s.repository.Token.CreateToken(ctx, uuid, tokens.RefreshToken)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *TokenService) RefreshToken(ctx context.Context, refreshToken string) (*model.Response, error) {
	token, err := jwt.Parse(
		refreshToken,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.jwtSecretKey), nil
		},
	)

	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			if validationErr.Errors&jwt.ValidationErrorExpired > 0 {
				return nil, errors.New("expiration date validation error")
			}
		}

		s.logger.Error(err)
		return nil, fmt.Errorf("failed to parse jwt err: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error(err)
		return nil, fmt.Errorf("unexpected type %T", claims)
	}

	t, err := s.repository.Token.GetToken(ctx, claims["uuid"].(string))
	if err != nil {
		return nil, err
	}
	fmt.Println(t, refreshToken)
	if t != refreshToken {
		s.logger.Error("token does not match with token in db")
		return nil, errors.New("token does not match with token in db")
	}

	tokens, err := s.generateToken(ctx, claims["uuid"].(string))
	if err != nil {
		s.logger.Error(err)
		return nil, fmt.Errorf("generating token err: %w", err)
	}

	err = s.repository.Token.UpdateToken(ctx, claims["uuid"].(string), tokens.RefreshToken)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *TokenService) generateToken(ctx context.Context, uuid string) (*model.Response, error) {
	accessTokenExpirationTime := time.Now().Add(time.Hour)
	refreshTokenExpirationTime := time.Now().Add(24 * time.Hour)

	accessTokenClaims := &model.JWTClaim{
		Uuid: uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpirationTime.Unix(),
		},
	}

	secretKey := []byte(s.jwtSecretKey)
	accessClaimToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessTokenString, err := accessClaimToken.SignedString(secretKey)
	if err != nil {
		s.logger.Errorf("AccessToken: SignedStrign err: %v", err)
		return nil, fmt.Errorf("AccessToken: SignedString err: %w", err)
	}

	refreshTokenClaims := &model.JWTClaim{
		Uuid: uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpirationTime.Unix(),
		},
	}

	refreshClaimToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	refreshTokenString, err := refreshClaimToken.SignedString(secretKey)
	if err != nil {
		s.logger.Errorf("RefreshToken: SignedString err: %v", err)
		return nil, fmt.Errorf("RefreshToken: SignedString err: %w", err)
	}

	jwtToken := &model.Response{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	return jwtToken, nil
}
