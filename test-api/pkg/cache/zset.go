package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SortedSetCache represents a Redis sorted set cache
type SortedSetCache struct {
	client *redis.Client
}

// NewSortedSetCache creates a new sorted set cache instance
func NewSortedSetCache(client *redis.Client) *SortedSetCache {
	return &SortedSetCache{client: client}
}

// ZAdd adds one or more members to a sorted set
func (z *SortedSetCache) ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	// Convert []*redis.Z to []redis.Z
	zs := make([]redis.Z, len(members))
	for i, m := range members {
		if m != nil {
			zs[i] = *m
		}
	}
	return z.client.ZAdd(ctx, key, zs...).Err()
}

// ZRem removes one or more members from a sorted set
func (z *SortedSetCache) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return z.client.ZRem(ctx, key, members...).Err()
}

// ZRange returns members in a sorted set by range
func (z *SortedSetCache) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return z.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores returns members with scores in a sorted set by range
func (z *SortedSetCache) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return z.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZRangeByScore returns members in a sorted set by score range
func (z *SortedSetCache) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return z.client.ZRangeByScore(ctx, key, opt).Result()
}

// ZRangeByScoreWithScores returns members with scores in a sorted set by score range
func (z *SortedSetCache) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	return z.client.ZRangeByScoreWithScores(ctx, key, opt).Result()
}

// ZRevRange returns members in a sorted set by range in descending order
func (z *SortedSetCache) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return z.client.ZRevRange(ctx, key, start, stop).Result()
}

// ZRevRangeWithScores returns members with scores in a sorted set by range in descending order
func (z *SortedSetCache) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return z.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

// ZCard returns the number of members in a sorted set
func (z *SortedSetCache) ZCard(ctx context.Context, key string) (int64, error) {
	return z.client.ZCard(ctx, key).Result()
}

// ZCount returns the number of members in a sorted set with scores between min and max
func (z *SortedSetCache) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	return z.client.ZCount(ctx, key, min, max).Result()
}

// ZScore returns the score of a member in a sorted set
func (z *SortedSetCache) ZScore(ctx context.Context, key, member string) (float64, error) {
	return z.client.ZScore(ctx, key, member).Result()
}

// ZIncrBy increments the score of a member in a sorted set
func (z *SortedSetCache) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	return z.client.ZIncrBy(ctx, key, increment, member).Result()
}

// ZRank returns the rank of a member in a sorted set
func (z *SortedSetCache) ZRank(ctx context.Context, key, member string) (int64, error) {
	return z.client.ZRank(ctx, key, member).Result()
}

// ZRevRank returns the rank of a member in a sorted set in descending order
func (z *SortedSetCache) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	return z.client.ZRevRank(ctx, key, member).Result()
}

// ZRemRangeByRank removes members in a sorted set by rank range
func (z *SortedSetCache) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) error {
	return z.client.ZRemRangeByRank(ctx, key, start, stop).Err()
}

// ZRemRangeByScore removes members in a sorted set by score range
func (z *SortedSetCache) ZRemRangeByScore(ctx context.Context, key, min, max string) error {
	return z.client.ZRemRangeByScore(ctx, key, min, max).Err()
}

// ZPopMax removes and returns the member with the highest score in a sorted set
func (z *SortedSetCache) ZPopMax(ctx context.Context, key string) ([]redis.Z, error) {
	return z.client.ZPopMax(ctx, key).Result()
}

// ZPopMin removes and returns the member with the lowest score in a sorted set
func (z *SortedSetCache) ZPopMin(ctx context.Context, key string) ([]redis.Z, error) {
	return z.client.ZPopMin(ctx, key).Result()
}
