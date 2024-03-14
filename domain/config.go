package domain

import "context"

type Config struct {
	Services        *AllServices
	ENV             string
	Version         string
	AdminAuthSecret string
	UserAuthSecret  string
	CacheFunc       func(ctx context.Context, key string) (string, error)
}
