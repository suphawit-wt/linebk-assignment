package config

import (
	"time"

	"github.com/gin-contrib/cors"
)

// GetCORSConfig maps the internal CORSConfig struct to a gin-contrib/cors Config.
func GetCORSConfig(cfg *CORSConfig) cors.Config {
	return cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     cfg.AllowedMethods,
		AllowHeaders:     cfg.AllowedHeaders,
		ExposeHeaders:    cfg.ExposedHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Duration(cfg.MaxAge) * time.Second,
	}
}
