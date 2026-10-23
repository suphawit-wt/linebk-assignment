package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
)

// NewRedisClient creates and initializes a new Redis client instance.
func NewRedisClient(appMode string, cfg *config.RedisConfig) (*redis.Client, error) {
	if cfg.Addr == "" {
		return nil, fmt.Errorf("redis setup error: address cannot be empty")
	}

	appMode = strings.ToLower(appMode)

	poolSize := 20
	minIdleConns := 10
	if appMode == envProd || appMode == envProduction || appMode == envStaging {
		poolSize = 100
		minIdleConns = 50
	}

	opts := &redis.Options{
		Addr:         cfg.Addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,

		DialTimeout:  2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,

		ConnMaxIdleTime: 5 * time.Minute,
		ConnMaxLifetime: 30 * time.Minute,
	}

	client := redis.NewClient(opts)

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping error: %w", err)
	}

	return client, nil
}
