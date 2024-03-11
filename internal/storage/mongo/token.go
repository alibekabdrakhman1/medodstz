package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type TokenRepository struct {
	DB     *mongo.Collection
	logger *zap.SugaredLogger
}

func NewTokenRepository(db *mongo.Collection, logger *zap.SugaredLogger) *TokenRepository {
	return &TokenRepository{
		DB:     db,
		logger: logger,
	}
}

func (r *TokenRepository) CreateToken(ctx context.Context, uuid, token string) error {
	check, err := r.uniqueCheck(ctx, uuid)
	if err != nil {
		r.logger.Error("unique checking error:", err)
		return err
	}
	if check {
		return errors.New("uuid already exists")
	}

	_, err = r.DB.InsertOne(ctx, bson.M{"uuid": uuid, "token": token})
	if err != nil {
		r.logger.Error("inserting token error: ", err)
		return err
	}
	r.logger.Info("token created successfully for uuid:", uuid)
	return nil
}

func (r *TokenRepository) UpdateToken(ctx context.Context, uuid, token string) error {
	filter := bson.M{"uuid": uuid}
	update := bson.M{"$set": bson.M{"token": token}}
	_, err := r.DB.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("updating token error:", err)
		return err
	}
	r.logger.Info("token updated successfully for uuid:", uuid)
	return nil
}

func (r *TokenRepository) GetToken(ctx context.Context, uuid string) (string, error) {
	var res struct {
		Token string `bson:"token"`
	}
	filter := bson.M{"uuid": uuid}
	err := r.DB.FindOne(ctx, filter).Decode(&res)
	fmt.Println(err)
	if err != nil {
		r.logger.Error("getting token error:", err)
		return "", err
	}
	r.logger.Info("token retrieved successfully for uuid:", uuid)
	return res.Token, nil
}

func (r *TokenRepository) uniqueCheck(ctx context.Context, uuid string) (bool, error) {
	_, err := r.GetToken(ctx, uuid)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
