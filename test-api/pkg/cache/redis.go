package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config represents Redis configuration
type Config struct {
	Host         string `yaml:"host" mapstructure:"host"`
	Port         int    `yaml:"port" mapstructure:"port"`
	Password     string `yaml:"password" mapstructure:"password"`
	DB           int    `yaml:"db" mapstructure:"db"`
	PoolSize     int    `yaml:"pool_size" mapstructure:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns" mapstructure:"min_idle_conns"`
	DialTimeout  int    `yaml:"dial_timeout" mapstructure:"dial_timeout"`   // seconds
	ReadTimeout  int    `yaml:"read_timeout" mapstructure:"read_timeout"`   // seconds
	WriteTimeout int    `yaml:"write_timeout" mapstructure:"write_timeout"` // seconds
	PoolTimeout  int    `yaml:"pool_timeout" mapstructure:"pool_timeout"`   // seconds
}

// Cache represents a Redis cache client
type Cache struct {
	client *redis.Client
}

// New creates a new Redis cache instance
func New(cfg *Config) (*Cache, error) {
	if cfg.PoolSize == 0 {
		cfg.PoolSize = 10
	}
	if cfg.MinIdleConns == 0 {
		cfg.MinIdleConns = 5
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 5
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 3
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 3
	}
	if cfg.PoolTimeout == 0 {
		cfg.PoolTimeout = 4
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	return &Cache{client: client}, nil
}

// Client returns the underlying Redis client
func (c *Cache) Client() *redis.Client {
	return c.client
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// Set sets a key-value pair with expiration
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Del deletes keys
func (c *Cache) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if keys exist
func (c *Cache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

// Expire sets expiration time for a key
func (c *Cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL returns the time to live for a key
func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// Incr increments a key's value
func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// Decr decrements a key's value
func (c *Cache) Decr(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

// IncrBy increments a key's value by delta
func (c *Cache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, key, value).Result()
}

// DecrBy decrements a key's value by delta
func (c *Cache) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.DecrBy(ctx, key, value).Result()
}

// SetNX sets a key-value pair only if the key does not exist
func (c *Cache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, value, expiration).Result()
}

// MGet gets multiple keys' values
func (c *Cache) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return c.client.MGet(ctx, keys...).Result()
}

// MSet sets multiple key-value pairs
func (c *Cache) MSet(ctx context.Context, values ...interface{}) error {
	return c.client.MSet(ctx, values...).Err()
}
