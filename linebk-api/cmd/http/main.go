package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/docs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/auth"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/session"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/user"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/database"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/jwt"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/passwd"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/security/vault"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/server"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/shared/jwtkey"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/shared/middleware"
)

//	@title			Banking Service - LINE BK API
//	@version		1.0.0
//	@description	Banking Service handling secure authentication, authorization and core banking operations.
//	@termsOfService	https://linebk.com/terms/

//	@contact.name	Engineering Team at KASIKORN LINE Co., Ltd.
//	@contact.url	https://linebk.com/engineering
//	@contact.email	engineering@linebk.com

//	@license.name	Copyright © KASIKORN LINE All Rights reserved
//	@license.url	https://linebk.com/licenses

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Put **_ONLY_** your JWT Bearer token in textbox below

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	if err := run(); err != nil {
		log.Printf("fatal startup/runtime error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	baseLogger, err := zaplog.NewLogger(cfg.App.Mode, &cfg.Logger)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer func() {
		_ = baseLogger.Sync()
	}()

	appLogger := baseLogger.Named("linebk").Named("banking_service")
	infraLogger := appLogger.Named("infra")
	dbLogger := infraLogger.Named("mysql")
	redisLogger := infraLogger.Named("redis")
	serverLogger := appLogger.Named("http")

	appLogger.Info("starting service initialization",
		zap.String("service", "banking_service"),
		zap.String("mode", string(cfg.App.Mode)),
	)

	db, err := database.NewMySQLClient(cfg.App.Mode, &cfg.DB.MySQL)
	if err != nil {
		return fmt.Errorf("init mysql: %w", err)
	}
	defer func() {
		dbLogger.Info("closing mysql client pool")
		if err = db.Close(); err != nil {
			dbLogger.Error("mysql close error", zap.Error(err))
		}
	}()

	rdb, err := database.NewRedisClient(cfg.App.Mode, &cfg.DB.Redis)
	if err != nil {
		return fmt.Errorf("init redis: %w", err)
	}
	defer func() {
		redisLogger.Info("closing redis client connection")
		if err = rdb.Close(); err != nil {
			redisLogger.Error("redis close error", zap.Error(err))
		}
	}()

	vaultClient, err := vault.NewClient(cfg.App.Mode, &cfg.Security.Vault)
	if err != nil {
		return fmt.Errorf("init vault client: %w", err)
	}

	jwtKeys, err := jwtkey.LoadJwtKeysFromVault(vaultClient)
	if err != nil {
		return fmt.Errorf("load JWT keys: %w", err)
	}

	jwtSvc := jwt.NewService(cfg.Security.JWT)
	passwordEncoder := passwd.NewArgon2IDEncoder(passwd.DefaultArgon2IDParams())

	srv, err := server.NewHTTPGin(cfg, serverLogger)
	if err != nil {
		return fmt.Errorf("init http server: %w", err)
	}

	userRepo := user.NewRepositoryDB(db)
	userService := user.NewService(appLogger.Named("user"), userRepo)
	userHandler := user.NewHandler(srv.Validator, userService)

	sessionRepo := session.NewRepositoryDB(rdb)
	authService := auth.NewService(
		appLogger.Named("auth"),
		&cfg.Auth,
		jwtSvc,
		jwtKeys,
		passwordEncoder,
		userRepo,
		sessionRepo,
	)
	authHandler := auth.NewHandler(srv.Validator, cfg, authService)

	registerRoutes(
		srv.Router,
		&cfg.App,
		rdb,
		jwtSvc,
		jwtKeys,
		authHandler,
		userHandler,
	)

	serverLogger.Info("http server starting", zap.Int("port", int(cfg.Server.HTTP.Port)))

	if err := srv.Run(ctx); err != nil {
		serverLogger.Error("server exited with error", zap.Error(err))
		return err
	}

	serverLogger.Info("http server gracefully shut down, closing resources and exiting")
	return nil
}

func registerRoutes(
	router *gin.Engine,
	cfg *config.AppConfig,
	rdb *redis.Client,
	jwtSvc jwt.Service,
	jwtKeys *jwtkey.Keys,
	authHandler auth.Handler,
	userHandler user.Handler,
) {
	v1 := router.Group("v1")
	authV1 := v1.Group("auth")

	switch cfg.Mode {
	case "dev", "development", "local":
		docs.SwaggerInfo.Title = cfg.Name
		docs.SwaggerInfo.Description = cfg.Description
		docs.SwaggerInfo.Version = cfg.Version

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	signinLimiter := middleware.RateLimiter(rdb, 5, 1*time.Minute)

	authV1.POST("/signin/passcode", signinLimiter, authHandler.SignInWithPasscode)
	authV1.POST("/signin/password", signinLimiter, authHandler.SignInWithPassword)
	authV1.POST("/signout", authHandler.SignOut)
	authV1.POST("/refresh", authHandler.RefreshToken)

	authorizedUserV1 := v1.Group("user").
		Use(middleware.AuthorizedUser(jwtSvc, jwtKeys.User.Access.Public))
	authorizedUserV1.GET("", userHandler.GetUserInfo)

	router.GET("/health", authHandler.HealthCheck)
}
