package card

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

// CachedRepository defines caching behavior for card data access operations.
type CachedRepository interface {
	GetDebitCards(ctx context.Context, userID string) ([]DebitCardInfo, error)
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

func (c *cachedRepository) GetDebitCards(
	ctx context.Context,
	userID string,
) ([]DebitCardInfo, error) {
	logger := c.logger.Named("get_debit_cards")
	cacheKey := fmt.Sprintf("cache:debit_cards:user:%s", userID)

	val, err := c.rdb.Get(ctx, cacheKey).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Error("redis get error, falling back to database", zap.Error(err))
	}

	if err != nil {
		logger.Debug(
			"cache miss or error for debit cards, falling back to db",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return c.fetchAndCacheDebitCards(ctx, cacheKey, userID)
	}

	var debitCards []DebitCardInfo
	if unmarshalErr := json.Unmarshal([]byte(val), &debitCards); unmarshalErr != nil {
		logger.Warn(
			"failed to unmarshal debit cards from cache, falling back to db",
			zap.Error(unmarshalErr),
		)
		return c.fetchAndCacheDebitCards(ctx, cacheKey, userID)
	}

	logger.Debug("cache hit for debit cards", zap.String("user_id", userID))
	return debitCards, nil
}

func (c *cachedRepository) fetchAndCacheDebitCards(
	ctx context.Context,
	cacheKey, userID string,
) ([]DebitCardInfo, error) {
	ttl := 5 * time.Minute

	debitCards, err := c.repo.GetDebitCards(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"cachedRepo.GetDebitCards: failed to fetch from database for user %s: %w",
			userID,
			err,
		)
	}

	dataBytes, jsonErr := json.Marshal(debitCards)
	if jsonErr == nil {
		if setErr := c.rdb.Set(ctx, cacheKey, dataBytes, ttl).Err(); setErr != nil {
			c.logger.Error("failed to write accounts to cache", zap.Error(setErr))
		}
	}

	return debitCards, nil
}
