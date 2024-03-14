package storage

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"go_base/logger"
	"go_base/xerror"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, opts *redis.Options, tlsServer string) (*Cache, error) {
	if tlsServer != "" {
		tlsConf := &tls.Config{
			ServerName:         tlsServer,
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}
		opts.TLSConfig = tlsConf
	}

	rdb := redis.NewClient(opts)

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	logger.L().Infof("Redis initialized: %v", rdb)

	return &Cache{Client: rdb}, nil
}

func (s *Cache) IncreaseStrike(ctx context.Context, key string) (int, error) {
	strikes, err := s.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(strikes), nil
}

func (s *Cache) SetCache(ctx context.Context, key string, value any, expiration time.Duration) error {
	if err := s.Client.Set(ctx, key, value, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Cache) GetCache(ctx context.Context, key string) ([]byte, error) {
	val, err := s.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, xerror.ENotFound()
		}
		return nil, err
	}
	return []byte(val), nil
}

func (s *Cache) SetExpire(ctx context.Context, key string, expiration time.Duration) error {
	if err := s.Client.Expire(ctx, key, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (s *Cache) GetBooleanValue(ctx context.Context, key string) (*bool, error) {
	value := false
	val, err := s.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return &value, nil
		}
		return nil, err
	}
	value, err = strconv.ParseBool(val)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func (s *Cache) GetStringValue(ctx context.Context, key string) (string, error) {
	val, err := s.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", xerror.ENotFound()
		}
		return "", err
	}
	return val, nil
}

func (s *Cache) GetStrikes(ctx context.Context, key string) (int, error) {
	ss, err := s.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	strikes, err := strconv.ParseInt(ss, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(strikes), nil
}

func (s *Cache) DeleteStrikes(ctx context.Context, key string) error {
	return s.ClearCache(ctx, key)
}

func (s *Cache) ClearCache(ctx context.Context, key string) error {
	_, err := s.Client.Del(ctx, key).Result()
	return err
}

func (s *Cache) HealthCheck(ctx context.Context) error {
	result, err := s.Client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	if result != "PONG" {
		return fmt.Errorf("unexpected response for redis ping: %s", result)
	}

	return nil
}
