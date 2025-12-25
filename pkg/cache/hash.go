package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// HashCache represents a Redis hash cache
type HashCache struct {
	client *redis.Client
}

// NewHashCache creates a new hash cache instance
func NewHashCache(client *redis.Client) *HashCache {
	return &HashCache{client: client}
}

// HSet sets field in the hash stored at key
func (h *HashCache) HSet(ctx context.Context, key, field string, value interface{}) error {
	return h.client.HSet(ctx, key, field, value).Err()
}

// HGet gets the value of a field in the hash stored at key
func (h *HashCache) HGet(ctx context.Context, key, field string) (string, error) {
	return h.client.HGet(ctx, key, field).Result()
}

// HMSet sets multiple fields in the hash stored at key
func (h *HashCache) HMSet(ctx context.Context, key string, values ...interface{}) error {
	return h.client.HMSet(ctx, key, values...).Err()
}

// HMGet gets multiple fields from the hash stored at key
func (h *HashCache) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	return h.client.HMGet(ctx, key, fields...).Result()
}

// HGetAll gets all fields and values in the hash stored at key
func (h *HashCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return h.client.HGetAll(ctx, key).Result()
}

// HDel deletes one or more fields from the hash stored at key
func (h *HashCache) HDel(ctx context.Context, key string, fields ...string) error {
	return h.client.HDel(ctx, key, fields...).Err()
}

// HExists checks if a field exists in the hash stored at key
func (h *HashCache) HExists(ctx context.Context, key, field string) (bool, error) {
	return h.client.HExists(ctx, key, field).Result()
}

// HLen returns the number of fields in the hash stored at key
func (h *HashCache) HLen(ctx context.Context, key string) (int64, error) {
	return h.client.HLen(ctx, key).Result()
}

// HIncrBy increments the value of a field in the hash stored at key
func (h *HashCache) HIncrBy(ctx context.Context, key, field string, value int64) (int64, error) {
	return h.client.HIncrBy(ctx, key, field, value).Result()
}

// HIncrByFloat increments the float value of a field in the hash stored at key
func (h *HashCache) HIncrByFloat(ctx context.Context, key, field string, value float64) (float64, error) {
	return h.client.HIncrByFloat(ctx, key, field, value).Result()
}

// HKeys returns all field names in the hash stored at key
func (h *HashCache) HKeys(ctx context.Context, key string) ([]string, error) {
	return h.client.HKeys(ctx, key).Result()
}

// HVals returns all values in the hash stored at key
func (h *HashCache) HVals(ctx context.Context, key string) ([]string, error) {
	return h.client.HVals(ctx, key).Result()
}

// HSetWithExpiration sets field in the hash stored at key with expiration
func (h *HashCache) HSetWithExpiration(ctx context.Context, key, field string, value interface{}, expiration time.Duration) error {
	pipe := h.client.Pipeline()
	pipe.HSet(ctx, key, field, value)
	pipe.Expire(ctx, key, expiration)
	_, err := pipe.Exec(ctx)
	return err
}
