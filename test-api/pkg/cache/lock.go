package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Lock represents a distributed lock
type Lock struct {
	client     *redis.Client
	key        string
	value      string
	expiration time.Duration
}

// NewLock creates a new distributed lock instance
func NewLock(client *redis.Client, key, value string, expiration time.Duration) *Lock {
	return &Lock{
		client:     client,
		key:        key,
		value:      value,
		expiration: expiration,
	}
}

// Lock acquires the lock
func (l *Lock) Lock(ctx context.Context) (bool, error) {
	return l.client.SetNX(ctx, l.key, l.value, l.expiration).Result()
}

// Unlock releases the lock
func (l *Lock) Unlock(ctx context.Context) error {
	// Use Lua script to ensure only the lock owner can unlock
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	_, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	return err
}

// Extend extends the lock expiration
func (l *Lock) Extend(ctx context.Context, expiration time.Duration) (bool, error) {
	// Use Lua script to ensure only the lock owner can extend
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("expire", KEYS[1], ARGV[2])
	else
		return 0
	end
	`
	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value, int(expiration.Seconds())).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

// TryLock tries to acquire the lock with a timeout
func (l *Lock) TryLock(ctx context.Context, timeout time.Duration) (bool, error) {
	deadline := time.Now().Add(timeout)
	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		acquired, err := l.Lock(ctx)
		if err != nil {
			return false, err
		}
		if acquired {
			return true, nil
		}

		if time.Now().After(deadline) {
			return false, nil
		}

		// Wait a bit before retrying
		time.Sleep(100 * time.Millisecond)
	}
}
