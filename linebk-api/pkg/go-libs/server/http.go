package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/middleware"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/validator"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
)

// HTTPGinServer wraps configuration, Gin router, HTTP server, and validator dependencies.
type HTTPGinServer struct {
	Config    *config.Config
	Router    *gin.Engine
	Server    *http.Server
	Validator *validator.Instance
}

// NewHTTPGin initializes a new Gin HTTP server instance with configured middlewares and validators.
func NewHTTPGin(cfg *config.Config, logger *zaplog.Logger) (*HTTPGinServer, error) {
	switch strings.ToLower(cfg.App.Mode) {
	case "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	case "staging", "test":
		gin.SetMode(gin.TestMode)
	case "dev", "development", "local", "":
		gin.SetMode(gin.DebugMode)
	default:
		return nil, fmt.Errorf("gin setup error: unknown env %q", cfg.App.Mode)
	}

	router := gin.New()

	router.Use(cors.New(config.GetCORSConfig(&cfg.Security.CORS)))
	router.Use(middleware.GinLoggerMiddleware(logger))
	router.Use(middleware.GinRecoveryMiddleware(logger))

	vdt, err := validator.New()
	if err != nil {
		return nil, fmt.Errorf("validator error: %w", err)
	}

	// Limit multipart form request to 8 MB
	router.MaxMultipartMemory = 8 << 20

	srv := &HTTPGinServer{
		Config: cfg,
		Router: router,
		Server: &http.Server{
			Addr:              fmt.Sprintf("0.0.0.0:%d", cfg.Server.HTTP.Port),
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
			IdleTimeout:       60 * time.Second,
			Handler:           router.Handler(),
		},
		Validator: vdt,
	}

	return srv, nil
}

// Run starts listening for HTTP requests and blocks until context cancellation, performing graceful shutdown.
func (h *HTTPGinServer) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	shutdownTimeout := 15 * time.Second

	// 1. Start HTTP Server in background goroutine
	go func() {
		if err := h.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	// 2. Block until startup error occurs or OS signal cancels parent context
	select {
	case err := <-errCh:
		// Early exit (e.g., port already bound, invalid TLS config)
		return err
	case <-ctx.Done():
		// OS Signal received (SIGINT/SIGTERM)
	}

	// 3. Create isolated timeout context for graceful HTTP connection draining
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()

	// 4. Trigger HTTP server shutdown (stops accepting new connections, drains active ones)
	shutdownErr := h.Server.Shutdown(shutdownCtx)

	// 5. Block until the background ListenAndServe goroutine completely terminates
	serverErr := <-errCh

	// Combine errors if both occurred (e.g., shutdown timeout + socket errors)
	return errors.Join(shutdownErr, serverErr)
}
