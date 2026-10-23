package transaction

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

// CachedRepository defines caching behavior for transaction data access operations.
type CachedRepository interface {
	GetRecentTransactions(ctx context.Context, userID string) ([]Transaction, error)
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

func (c *cachedRepository) GetRecentTransactions(
	ctx context.Context,
	userID string,
) ([]Transaction, error) {
	logger := c.logger.Named("get_recent_transactions")
	cacheKey := fmt.Sprintf("cache:recent_transactions:user:%s", userID)

	val, err := c.rdb.Get(ctx, cacheKey).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Error("redis get error, falling back to database", zap.Error(err))
	}

	if err != nil {
		logger.Debug(
			"cache miss or error for recent transactions, falling back to db",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return c.fetchAndCacheRecentTransactions(ctx, cacheKey, userID)
	}

	var recentTransactions []Transaction
	if unmarshalErr := json.Unmarshal([]byte(val), &recentTransactions); unmarshalErr != nil {
		logger.Warn(
			"failed to unmarshal recent transactions from cache, falling back to db",
			zap.Error(unmarshalErr),
		)
		return c.fetchAndCacheRecentTransactions(ctx, cacheKey, userID)
	}

	logger.Debug("cache hit for recent transactions", zap.String("user_id", userID))
	return recentTransactions, nil
}

func (c *cachedRepository) fetchAndCacheRecentTransactions(
	ctx context.Context,
	cacheKey, userID string,
) ([]Transaction, error) {
	ttl := 30 * time.Second

	recentTransactions, err := c.repo.GetRecentTransactions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"cachedRepo.GetRecentTransactions: failed to fetch from database for user %s: %w",
			userID,
			err,
		)
	}

	dataBytes, jsonErr := json.Marshal(recentTransactions)
	if jsonErr == nil {
		if setErr := c.rdb.Set(ctx, cacheKey, dataBytes, ttl).Err(); setErr != nil {
			c.logger.Error("failed to write accounts to cache", zap.Error(setErr))
		}
	}

	return recentTransactions, nil
}
