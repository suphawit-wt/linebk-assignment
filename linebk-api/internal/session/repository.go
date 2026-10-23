package session

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Metadata represents the metadata associated with an auth session
type Metadata struct {
	SessionID string `json:"sessionId"`
	TokenHash string `json:"tokenHash"`
	IPAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// Repository defines the data access contract for jwt tokens sessions.
type Repository interface {
	GetMetadata(ctx context.Context, sessionKey string) (*Metadata, error)
	Create(
		ctx context.Context,
		sessionID, sessionKey, accSessionsKey string,
		metadata *Metadata,
		ttl time.Duration,
	) error
	Revoke(ctx context.Context, sessionID, sessionKey, accSessionsKey string) (bool, error)
	TrimSessionsByCapacity(
		ctx context.Context,
		accID, accSessionsKey, clientKey string,
		maxAllowed int64,
	) error
}

type repositoryDB struct {
	rdb *redis.Client
}

// NewRepositoryDB creates a new instance of sessionRepositoryDB with the provided Redis client.
func NewRepositoryDB(rdb *redis.Client) Repository {
	return &repositoryDB{rdb: rdb}
}

func (r *repositoryDB) GetMetadata(
	ctx context.Context,
	sessionKey string,
) (*Metadata, error) {
	var metadata Metadata

	val, err := r.rdb.Get(ctx, sessionKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrTokenNotFound
		}

		return nil, fmt.Errorf(
			"sessionRepo.GetMetadata: redis GET failed for key %s: %w",
			sessionKey,
			err,
		)
	}

	if err := json.Unmarshal([]byte(val), &metadata); err != nil {
		return nil, fmt.Errorf(
			"sessionRepo.GetMetadata: failed to unmarshal JSON for key %s: %w",
			sessionKey,
			err,
		)
	}

	return &metadata, nil
}

func (r *repositoryDB) Create(
	ctx context.Context,
	sessionID, sessionKey, accSessionsKey string,
	metadata *Metadata,
	ttl time.Duration,
) error {
	score := float64(time.Now().Unix())

	payload, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf(
			"sessionRepo.Create: failed to marshal payload for key %s: %w",
			sessionKey,
			err,
		)
	}

	if err = r.rdb.Set(ctx, sessionKey, payload, ttl).Err(); err != nil {
		return fmt.Errorf("sessionRepo.Create: redis SET failed for key %s: %w", sessionKey, err)
	}

	pipe := r.rdb.Pipeline()
	pipe.ZAdd(ctx, accSessionsKey, redis.Z{
		Score:  score,
		Member: sessionID,
	})

	if ttl > 0 {
		pipe.Expire(ctx, accSessionsKey, ttl)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf(
			"sessionRepo.Create: failed to create session zset for key %s: %w",
			accSessionsKey,
			err,
		)
	}

	return nil
}

func (r *repositoryDB) Revoke(
	ctx context.Context,
	sessionID, sessionKey, accSessionsKey string,
) (bool, error) {
	deletedCount, err := r.rdb.Del(ctx, sessionKey).Result()
	if err != nil {
		return false, err
	}

	err = r.rdb.ZRem(ctx, accSessionsKey, sessionID).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, fmt.Errorf(
			"sessionRepo.Revoke: failed for key %s: %w",
			accSessionsKey,
			err,
		)
	}

	return deletedCount > 0, nil
}

func (r *repositoryDB) TrimSessionsByCapacity(
	ctx context.Context,
	accID,
	accSessionsKey,
	clientKey string,
	maxAllowed int64,
) error {
	// 1. Get current count of active sessions from Redis ZSET
	count, err := r.rdb.ZCard(ctx, accSessionsKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf(
			"sessionRepo.TrimSessionsByCapacity: failed for key %s: %w",
			accSessionsKey,
			err,
		)
	}

	// 2. FIFO Eviction if session limit reached
	if count >= maxAllowed {
		overflow := count - maxAllowed + 1

		// Fetch the oldest session IDs (lowest scores / timestamps)
		oldestSessionIDs, err := r.rdb.ZRange(ctx, accSessionsKey, 0, overflow-1).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf(
				"sessionRepo.TrimSessionsByCapacity: failed for key %s: %w",
				accSessionsKey,
				err,
			)
		}

		for _, oldSessionID := range oldestSessionIDs {
			oldSessionKey := fmt.Sprintf("%s:%s:session:%s", clientKey, accID, oldSessionID)
			_, err = r.Revoke(ctx, oldSessionID, oldSessionKey, accSessionsKey)
			if err != nil {
				return fmt.Errorf(
					"sessionRepo.TrimSessionsByCapacity: failed for key %s: %w",
					accSessionsKey,
					err,
				)
			}
		}
	}

	return nil
}
