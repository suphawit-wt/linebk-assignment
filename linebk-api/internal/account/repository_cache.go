package account

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

// CachedRepository defines caching behavior for account data access operations.
type CachedRepository interface {
	GetAccountsByUserID(ctx context.Context, userID string) ([]Info, error)
}

type cachedRepository struct {
	logger *zaplog.Logger
	repo   Repository
	rdb    *redis.Client
}

// NewCachedRepository creates a new cached decorator instance
func NewCachedRepository(
	logger *zaplog.Logger,
	repo Repository,
	rdb *redis.Client,
) CachedRepository {
	return &cachedRepository{
		repo:   repo,
		rdb:    rdb,
		logger: logger.Named("cache_repo"),
	}
}

func (c *cachedRepository) GetAccountsByUserID(ctx context.Context, userID string) ([]Info, error) {
	logger := c.logger.Named("get_accounts_by_user_id")
	cacheKey := fmt.Sprintf("cache:accounts:user:%s", userID)

	val, err := c.rdb.Get(ctx, cacheKey).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Error("redis get error, falling back to database", zap.Error(err))
	}

	if err != nil {
		logger.Debug(
			"cache miss or error for accounts, falling back to db",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return c.fetchAndCacheAccounts(ctx, cacheKey, userID)
	}

	var accounts []Info
	if unmarshalErr := json.Unmarshal([]byte(val), &accounts); unmarshalErr != nil {
		logger.Warn(
			"failed to unmarshal accounts from cache, falling back to db",
			zap.Error(unmarshalErr),
		)
		return c.fetchAndCacheAccounts(ctx, cacheKey, userID)
	}

	logger.Debug("cache hit for accounts", zap.String("user_id", userID))
	return accounts, nil
}

func (c *cachedRepository) fetchAndCacheAccounts(
	ctx context.Context,
	cacheKey, userID string,
) ([]Info, error) {
	ttl := 10 * time.Second

	accounts, err := c.repo.GetAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"cachedRepo.GetAccountsByUserID: failed to fetch from database for user %s: %w",
			userID,
			err,
		)
	}

	dataBytes, jsonErr := json.Marshal(accounts)
	if jsonErr == nil {
		if setErr := c.rdb.Set(ctx, cacheKey, dataBytes, ttl).Err(); setErr != nil {
			c.logger.Error("failed to write accounts to cache", zap.Error(setErr))
		}
	}

	return accounts, nil
}
