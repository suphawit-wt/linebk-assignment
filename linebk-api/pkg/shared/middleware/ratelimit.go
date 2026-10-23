package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/dto"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/enums/api"
)

// RateLimiter limits requests based on client IP using Redis.
func RateLimiter(rdb *redis.Client, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		clientIP := c.ClientIP()

		key := fmt.Sprintf("ratelimit:signin:%s", clientIP)

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			_ = rdb.Expire(ctx, key, window).Err()
		}

		if count > int64(maxRequests) {
			c.Header("Retry-After", fmt.Sprintf("%.0f", window.Seconds()))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.ApiResponse{
				Status:     api.StatusTooManyRequests,
				StatusCode: http.StatusTooManyRequests,
				Message:    "too many sign-in attempts, please try again later.",
			})
			return
		}

		c.Next()
	}
}
