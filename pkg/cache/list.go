package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ListCache represents a Redis list cache
type ListCache struct {
	client *redis.Client
}

// NewListCache creates a new list cache instance
func NewListCache(client *redis.Client) *ListCache {
	return &ListCache{client: client}
}

// LPush inserts all values at the head of the list stored at key
func (l *ListCache) LPush(ctx context.Context, key string, values ...interface{}) error {
	return l.client.LPush(ctx, key, values...).Err()
}

// RPush inserts all values at the tail of the list stored at key
func (l *ListCache) RPush(ctx context.Context, key string, values ...interface{}) error {
	return l.client.RPush(ctx, key, values...).Err()
}

// LPop removes and returns the first element of the list stored at key
func (l *ListCache) LPop(ctx context.Context, key string) (string, error) {
	return l.client.LPop(ctx, key).Result()
}

// RPop removes and returns the last element of the list stored at key
func (l *ListCache) RPop(ctx context.Context, key string) (string, error) {
	return l.client.RPop(ctx, key).Result()
}

// LLen returns the length of the list stored at key
func (l *ListCache) LLen(ctx context.Context, key string) (int64, error) {
	return l.client.LLen(ctx, key).Result()
}

// LRange returns the specified elements of the list stored at key
func (l *ListCache) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return l.client.LRange(ctx, key, start, stop).Result()
}

// LIndex returns the element at index in the list stored at key
func (l *ListCache) LIndex(ctx context.Context, key string, index int64) (string, error) {
	return l.client.LIndex(ctx, key, index).Result()
}

// LSet sets the list element at index to value
func (l *ListCache) LSet(ctx context.Context, key string, index int64, value interface{}) error {
	return l.client.LSet(ctx, key, index, value).Err()
}

// LTrim trims an existing list so that it will contain only the specified range of elements
func (l *ListCache) LTrim(ctx context.Context, key string, start, stop int64) error {
	return l.client.LTrim(ctx, key, start, stop).Err()
}

// LRem removes the first count occurrences of elements equal to value from the list stored at key
func (l *ListCache) LRem(ctx context.Context, key string, count int64, value interface{}) error {
	return l.client.LRem(ctx, key, count, value).Err()
}

// LInsertBefore inserts value before the pivot value in the list stored at key
func (l *ListCache) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) error {
	return l.client.LInsertBefore(ctx, key, pivot, value).Err()
}

// LInsertAfter inserts value after the pivot value in the list stored at key
func (l *ListCache) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) error {
	return l.client.LInsertAfter(ctx, key, pivot, value).Err()
}

// LPopCount removes and returns the first count elements of the list stored at key
func (l *ListCache) LPopCount(ctx context.Context, key string, count int) ([]string, error) {
	return l.client.LPopCount(ctx, key, count).Result()
}

// RPopCount removes and returns the last count elements of the list stored at key
func (l *ListCache) RPopCount(ctx context.Context, key string, count int) ([]string, error) {
	return l.client.RPopCount(ctx, key, count).Result()
}
