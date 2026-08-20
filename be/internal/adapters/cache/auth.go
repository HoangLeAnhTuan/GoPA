package cache

import (
	"context"
	"fmt"
	"time"

	"gopa/internal/constants"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RefreshSessionStore struct {
	client *redis.Client
}

type RateLimiter struct {
	client *redis.Client
}

func NewRefreshSessionStore(client *redis.Client) *RefreshSessionStore {
	return &RefreshSessionStore{client: client}
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{client: client}
}

func (s *RefreshSessionStore) Save(ctx context.Context, userID, sessionID uuid.UUID, tokenHash string, ttl time.Duration) error {
	if err := s.client.Set(ctx, refreshSessionKey(userID, sessionID), tokenHash, ttl).Err(); err != nil {
		return fmt.Errorf("save refresh session: %w", err)
	}
	return nil
}

func (s *RefreshSessionStore) Rotate(ctx context.Context, userID, oldSessionID uuid.UUID, oldTokenHash string, newSessionID uuid.UUID, newTokenHash string, ttl time.Duration) (bool, error) {
	result, err := s.client.Eval(ctx, `
local current = redis.call('GET', KEYS[1])
if not current or current ~= ARGV[1] then return 0 end
redis.call('DEL', KEYS[1])
redis.call('SET', KEYS[2], ARGV[2], 'PX', ARGV[3])
return 1`, []string{refreshSessionKey(userID, oldSessionID), refreshSessionKey(userID, newSessionID)}, oldTokenHash, newTokenHash, ttl.Milliseconds()).Int()
	if err != nil {
		return false, fmt.Errorf("rotate refresh session: %w", err)
	}
	return result == 1, nil
}

func (s *RefreshSessionStore) Delete(ctx context.Context, userID, sessionID uuid.UUID) error {
	if err := s.client.Del(ctx, refreshSessionKey(userID, sessionID)).Err(); err != nil {
		return fmt.Errorf("delete refresh session: %w", err)
	}
	return nil
}

func (l *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	redisKey := constants.CacheRateLimitPrefix + key
	count, err := l.client.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, fmt.Errorf("increment rate limit: %w", err)
	}
	if count == 1 {
		if err := l.client.Expire(ctx, redisKey, window).Err(); err != nil {
			return false, fmt.Errorf("set rate limit expiry: %w", err)
		}
	}
	return count <= int64(limit), nil
}

func refreshSessionKey(userID, sessionID uuid.UUID) string {
	return constants.CacheAuthRefreshPrefix + userID.String() + ":" + sessionID.String()
}
