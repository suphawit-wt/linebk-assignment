package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/docs"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/internal/auth"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/server"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/zaplog"
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
	serverLogger := appLogger.Named("http")

	appLogger.Info("starting service initialization",
		zap.String("service", "banking_service"),
		zap.String("mode", string(cfg.App.Mode)),
	)

	srv, err := server.NewHTTPGin(cfg, serverLogger)
	if err != nil {
		return fmt.Errorf("init http server: %w", err)
	}

	authHandler := auth.NewHandler(srv.Validator, cfg)

	registerRoutes(
		srv.Router,
		&cfg.App,
		authHandler,
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
	authHandler auth.Handler,
) {
	switch cfg.Mode {
	case "dev", "development", "local":
		docs.SwaggerInfo.Title = cfg.Name
		docs.SwaggerInfo.Description = cfg.Description
		docs.SwaggerInfo.Version = cfg.Version

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.GET("/health", authHandler.HealthCheck)
}
