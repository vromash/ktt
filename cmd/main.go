package main

import (
	"context"
	"financing-aggregator/internal/app"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"

	"financing-aggregator/internal/config"
	"go.uber.org/zap"
)

// @title           			Financial Aggregator
// @version         			0.1.0
// @securityDefinitions.apikey 	BearerAuth
// @in 							header
// @name 						Authorization
func main() {
	cfg := config.ReadConfig()

	var logger *zap.Logger
	if cfg.Env == "prod" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	db, err := getDBClient(cfg.DB)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	app, err := app.New(cfg, logger, db)
	if err != nil {
		logger.Fatal("failed to create app", zap.Error(err))
	}

	go func() {
		if err := app.Run(); err != nil {
			logger.Fatal("failed to start app", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit
	logger.Info("gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		logger.Fatal("forcefully shutting down server", zap.Error(err))
	}

	logger.Info("service stopped")
}

func getDBClient(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection with database: %v", err)
	}

	return db, nil
}
