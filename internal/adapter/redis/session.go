package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

func generateAccessKey(token string) string {
	return "access:" + token
}

func generateUserKey(userID uint) string {
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	return "access:user:" + userIDStr
}

// Set into Redis session data with expiration
//
// Key-Value structure:
//
//	access:<token> -> userID
//	access:user:<userID> -> Set(token1, token2, token3 ...)
func (r *Redis) AddSession(ctx context.Context, userID uint, accessToken string, expiration time.Duration) error {
	// Bind logger with request ID
	log := logger.LoggerWithKey(r.logger, ctx, common.ContexKeyRequestID)

	log.Info("Redis AddSession called",
		log.ToInt("userID", int(userID)),
	)

	// Generate access key and value
	key := generateAccessKey(accessToken)
	value := userID

	// Set session with expiration
	err := r.client.Set(r.context, key, value, expiration).Err()
	if err != nil {
		log.Error("Redis set session error", log.ToString("key", key),
			log.ToString("accessToken", accessToken),
			log.ToString("error", err.Error()))
		return err
	}

	// Update the set of tokens in Redis with current token
	userKey := generateUserKey(userID)

	// Add the new token to the user's set
	err = r.client.SAdd(r.context, userKey, accessToken).Err()
	if err != nil {
		log.Error("redis sadd user token error", log.ToString("userKey", userKey),
			log.ToString("error", err.Error()))
		return err
	}

	// Clean up expired tokens from the user's set
	tokens, err := r.client.SMembers(r.context, userKey).Result()
	if err != nil {
		log.Error("redis smembers user tokens error", log.ToString("userKey", userKey),
			log.ToString("error", err.Error()))
		return err
	}

	// Check each token for relevance (because Redis does not auto-remove expired keys in sets)
	for _, token := range tokens {
		accessKey := generateAccessKey(token)
		exists, err := r.client.Exists(r.context, accessKey).Result()
		if err != nil {
			log.Error("redis exists access key error", log.ToString("accessKey", accessKey),
				log.ToString("error", err.Error()))
			return err
		}
		if exists == 0 {
			// Token expired, remove from set
			err := r.client.SRem(r.context, userKey, token).Err()
			if err != nil {
				log.Error("redis srem expired token error", log.ToString("userKey", userKey),
					log.ToString("token", token),
					log.ToString("error", err.Error()))
				return err
			}
			log.Info("redis removed expired token from user set", log.ToString("userKey", userKey), log.ToString("token", token))
		}
	}

	log.Info("redis added session", log.ToString("key", key), log.ToString("accessToken", accessToken))
	return nil
}

// GetUserIDByToken retrieves userID associated with the given token from Redis
func (r *Redis) GetUserIDByToken(ctx context.Context, accessToken string) (uint, error) {
	log := logger.LoggerWithKey(r.logger, ctx, common.ContexKeyRequestID)
	log.Debug("Redis GetUserIDByToken called", log.ToString("accessToken", accessToken))

	// Generate access key
	key := generateAccessKey(accessToken)

	// Get userID from Redis
	result, err := r.client.Get(r.context, key).Result()
	if err != nil {
		log.Warn("redis get userID by token error", log.ToString("key", key),
			log.ToString("error", err.Error()))
		return 0, err
	}

	// Convert result to uint
	userIDUint64, err := strconv.ParseUint(result, 10, 64)
	if err != nil {
		log.Error("redis parse userID error", log.ToString("key", key),
			log.ToString("error", err.Error()))
		return 0, err
	}
	userID := uint(userIDUint64)

	log.Debug("redis get userID by token", log.ToString("key", key), log.ToInt("userID", int(userID)))
	return userID, nil
}

// DeleteAllSession removes all session data associated with the given userID from Redis
func (r *Redis) DeleteAllSession(ctx context.Context, userID uint) error {
	log := logger.LoggerWithKey(r.logger, ctx, common.ContexKeyRequestID)
	log.Info("Redis DeleteAllSession called", log.ToInt("userID", int(userID)))

	// Get user access keys
	userKey := generateUserKey(userID)
	tokens, err := r.client.SMembers(r.context, userKey).Result()
	if err != nil {
		log.Error("redis get user tokens for deletion error", log.ToString("userKey", userKey),
			log.ToString("error", err.Error()))
		return err
	}

	// Delete each access key
	for _, token := range tokens {
		accessKey := generateAccessKey(token)
		err := r.client.Del(r.context, accessKey).Err()
		if err != nil {
			log.Error("redis delete access key error", log.ToString("accessKey", accessKey),
				log.ToString("error", err.Error()))
			return err
		}
		log.Info("redis deleted access key", log.ToString("accessKey", accessKey))
	}

	// Delete user key
	err = r.client.Del(r.context, userKey).Err()
	if err != nil {
		log.Error("redis delete user key error", log.ToString("userKey", userKey),
			log.ToString("error", err.Error()))
		return err
	}

	log.Info("redis deleted user key", log.ToString("userKey", userKey))
	return nil
}

// Remove ONE session (access token) for the given userID from Redis
func (r *Redis) DeleteSession(ctx context.Context, userID uint, accessToken string) error {
	log := logger.LoggerWithKey(r.logger, ctx, common.ContexKeyRequestID)
	log.Debug("Redis DeleteSession called",
		log.ToInt("userID", int(userID)),
		log.ToString("accessToken", accessToken),
	)

	// Delete access key
	accessKey := generateAccessKey(accessToken)
	err := r.client.Del(r.context, accessKey).Err()
	if err != nil {
		log.Error("redis delete access key error", log.ToString("accessKey", accessKey),
			log.ToString("error", err.Error()))
		return err
	}
	log.Debug("redis deleted access key", log.ToString("accessKey", accessKey))

	// Remove token from user set
	userKey := generateUserKey(userID)
	err = r.client.SRem(r.context, userKey, accessToken).Err()
	if err != nil {
		log.Error("redis remove token from user set error", log.ToString("userKey", userKey),
			log.ToString("error", err.Error()))
		return err
	}
	log.Info("redis removed token from user set", log.ToString("userKey", userKey))

	return nil
}
