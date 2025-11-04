package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
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
	// Log the request ID from context for tracing
	requestID, ok := ctx.Value(common.RequestIDContextKey).(string)
	if !ok {
		r.logger.Warn("Redis AddSession: failed to get requestID from context")
		requestID = "unknown"
	}
	r.logger.Info("Redis AddSession called",
		r.logger.ToString("requestID", requestID),
		r.logger.ToInt("userID", int(userID)),
	)
	// Generate access key and value
	key := generateAccessKey(accessToken)
	value := userID

	// Set session with expiration
	err := r.client.Set(r.context, key, value, expiration).Err()
	if err != nil {
		r.logger.Error("Redis set session error", r.logger.ToString("key", key),
			r.logger.ToString("accessToken", accessToken),
			r.logger.ToString("error", err.Error()))
		return err
	}

	// Update the set of tokens in Redis with current token
	userKey := generateUserKey(userID)

	err = r.client.SAdd(r.context, userKey, accessToken).Err()
	if err != nil {
		r.logger.Error("redis update user tokens error", r.logger.ToString("userKey", userKey),
			r.logger.ToString("error", err.Error()))
		return err
	}

	r.logger.Info("redis added session", r.logger.ToString("key", key), r.logger.ToString("accessToken", accessToken))
	return nil
}

// GetUserIDByToken retrieves userID associated with the given token from Redis
func (r *Redis) GetUserIDByToken(ctx context.Context, accessToken string) (uint, error) {
	// Generate access key
	key := generateAccessKey(accessToken)

	// Get userID from Redis
	result, err := r.client.Get(r.context, key).Result()
	if err != nil {
		r.logger.Warn("redis get userID by token error", r.logger.ToString("key", key),
			r.logger.ToString("error", err.Error()))
		return 0, err
	}

	// Convert result to uint
	userIDUint64, err := strconv.ParseUint(result, 10, 64)
	if err != nil {
		r.logger.Error("redis parse userID error", r.logger.ToString("key", key),
			r.logger.ToString("error", err.Error()))
		return 0, err
	}
	userID := uint(userIDUint64)

	r.logger.Info("redis get userID by token", r.logger.ToString("key", key), r.logger.ToInt("userID", int(userID)))
	return userID, nil
}

// DeleteSession removes all session data associated with the given userID from Redis
func (r *Redis) DeleteSession(ctx context.Context, userID uint) error {
	// Get user access keys
	userKey := generateUserKey(userID)
	tokens, err := r.client.SMembers(r.context, userKey).Result()
	if err != nil {
		r.logger.Error("redis get user tokens for deletion error", r.logger.ToString("userKey", userKey),
			r.logger.ToString("error", err.Error()))
		return err
	}

	// Delete each access key
	for _, token := range tokens {
		accessKey := generateAccessKey(token)
		err := r.client.Del(r.context, accessKey).Err()
		if err != nil {
			r.logger.Error("redis delete access key error", r.logger.ToString("accessKey", accessKey),
				r.logger.ToString("error", err.Error()))
			return err
		}
		r.logger.Info("redis deleted access key", r.logger.ToString("accessKey", accessKey))
	}

	// Delete user key
	err = r.client.Del(r.context, userKey).Err()
	if err != nil {
		r.logger.Error("redis delete user key error", r.logger.ToString("userKey", userKey),
			r.logger.ToString("error", err.Error()))
		return err
	}

	r.logger.Info("redis deleted user key", r.logger.ToString("userKey", userKey))
	return nil
}
