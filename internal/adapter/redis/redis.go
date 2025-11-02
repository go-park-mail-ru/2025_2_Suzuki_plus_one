package redis

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	logger  logger.Logger
	client  *redis.Client
	context context.Context
}

func NewRedis(logger logger.Logger, connString string, password string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     connString,
		Password: password,
		DB:       0, // Use default DB
		Protocol: 2, // Connection protocol
	})
	return &Redis{
		logger:  logger,
		client:  client,
		context: context.Background(),
	}
}

func (r *Redis) CheckConnection() error {
	_, err := r.client.Ping(r.context).Result()
	if err != nil {
		r.logger.Error("Failed to connect to Redis: " + err.Error())
		return err
	}
	r.logger.Info("Connected to Redis successfully")
	return nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Get(key string) (string, error) {
	result, err := r.client.Get(r.context, key).Result()
	if err != nil {
		r.logger.Warn("redis can't see the key", r.logger.ToString("key", key))
		return "", err
	}

	r.logger.Info("redis get key", r.logger.ToString("key", key), r.logger.ToString("value", result))
	return result, nil
}

func (r *Redis) Set(key string, value string) error {
	err := r.client.Set(r.context, key, value, 0).Err()
	if err != nil {
		r.logger.Error("redis set error", r.logger.ToString("key", key),
			r.logger.ToString("value", value),
			r.logger.ToString("error", err.Error()))
		return err
	}

	r.logger.Info("redis set key", r.logger.ToString("key", key), r.logger.ToString("value", value))
	return nil
}
