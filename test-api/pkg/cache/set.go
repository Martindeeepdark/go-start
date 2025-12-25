package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SetCache represents a Redis set cache
type SetCache struct {
	client *redis.Client
}

// NewSetCache creates a new set cache instance
func NewSetCache(client *redis.Client) *SetCache {
	return &SetCache{client: client}
}

// SAdd adds one or more members to a set
func (s *SetCache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return s.client.SAdd(ctx, key, members...).Err()
}

// SRem removes one or more members from a set
func (s *SetCache) SRem(ctx context.Context, key string, members ...interface{}) error {
	return s.client.SRem(ctx, key, members...).Err()
}

// SPop removes and returns a random member from a set
func (s *SetCache) SPop(ctx context.Context, key string) (string, error) {
	return s.client.SPop(ctx, key).Result()
}

// SPopN removes and returns count random members from a set
func (s *SetCache) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	return s.client.SPopN(ctx, key, count).Result()
}

// SMembers returns all members of a set
func (s *SetCache) SMembers(ctx context.Context, key string) ([]string, error) {
	return s.client.SMembers(ctx, key).Result()
}

// SIsMember checks if a member is in a set
func (s *SetCache) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return s.client.SIsMember(ctx, key, member).Result()
}

// SMIsMember checks if multiple members are in a set
func (s *SetCache) SMIsMember(ctx context.Context, key string, members ...interface{}) ([]bool, error) {
	return s.client.SMIsMember(ctx, key, members...).Result()
}

// SCard returns the number of members in a set
func (s *SetCache) SCard(ctx context.Context, key string) (int64, error) {
	return s.client.SCard(ctx, key).Result()
}

// SMove moves a member from one set to another
func (s *SetCache) SMove(ctx context.Context, source, destination string, member interface{}) (bool, error) {
	return s.client.SMove(ctx, source, destination, member).Result()
}

// SDiff returns the difference of multiple sets
func (s *SetCache) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	return s.client.SDiff(ctx, keys...).Result()
}

// SDiffStore stores the difference of multiple sets in destination
func (s *SetCache) SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	return s.client.SDiffStore(ctx, destination, keys...).Result()
}

// SInter returns the intersection of multiple sets
func (s *SetCache) SInter(ctx context.Context, keys ...string) ([]string, error) {
	return s.client.SInter(ctx, keys...).Result()
}

// SInterStore stores the intersection of multiple sets in destination
func (s *SetCache) SInterStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	return s.client.SInterStore(ctx, destination, keys...).Result()
}

// SUnion returns the union of multiple sets
func (s *SetCache) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return s.client.SUnion(ctx, keys...).Result()
}

// SUnionStore stores the union of multiple sets in destination
func (s *SetCache) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	return s.client.SUnionStore(ctx, destination, keys...).Result()
}

// SRandMember returns a random member from a set
func (s *SetCache) SRandMember(ctx context.Context, key string) (string, error) {
	return s.client.SRandMember(ctx, key).Result()
}

// SRandMemberN returns count random members from a set
func (s *SetCache) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	return s.client.SRandMemberN(ctx, key, count).Result()
}
