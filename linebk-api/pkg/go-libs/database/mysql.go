package database

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
)

// NewMySQLClient creates and configures a new MySQL connection pool using sqlx.
func NewMySQLClient(appMode string, cfg *config.MySQLConfig) (*sqlx.DB, error) {
	appMode = strings.ToLower(appMode)

	queryParams := url.Values{}
	queryParams.Set("parseTime", "true")
	queryParams.Set("charset", "utf8mb4")
	queryParams.Set("collation", "utf8mb4_unicode_ci")

	sslMode := cfg.SSLMode
	if sslMode == "" {
		if appMode == envProd || appMode == envProduction || appMode == envStaging {
			queryParams.Set("tls", "true")
		} else {
			queryParams.Set("tls", "false")
		}
	} else {
		queryParams.Set("tls", sslMode)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
		queryParams.Encode(),
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("mysql connection error: %w", err)
	}

	maxOpenConns := 25
	maxIdleConns := 25
	if appMode == envProd || appMode == envProduction || appMode == envStaging {
		maxOpenConns = 100
		maxIdleConns = 100
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("mysql ping error: %w", err)
	}

	return db, nil
}
