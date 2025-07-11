package app

import (
	"context"
	"financing-aggregator/internal/banks"
	"financing-aggregator/internal/banks/fastbank"
	"financing-aggregator/internal/banks/solidbank"
	"financing-aggregator/internal/config"
	"financing-aggregator/internal/controllers"
	httpHandlers "financing-aggregator/internal/controllers/http"
	"financing-aggregator/internal/controllers/ws"
	"financing-aggregator/internal/repositories"
	"financing-aggregator/internal/services"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
	db     *gorm.DB
	cron   gocron.Scheduler
	srv    *http.Server
}

func New(cfg *config.Config, logger *zap.Logger, db *gorm.DB) (*App, error) {
	cron, err := gocron.NewScheduler()
	if err != nil {
		return nil, errors.Errorf("failed to create scheduler: %s", err.Error())
	}

	return &App{
		cfg:    cfg,
		logger: logger,
		db:     db,
		cron:   cron,
	}, nil
}

func (a *App) Run() error {
	defer func() {
		if r := recover(); r != nil {
			a.logger.Error("recovering panic", zap.Any("panic", r))
		}
	}()

	fastBank := fastbank.NewFastBank(a.cfg.Banks.FastBankURL)
	solidBank := solidbank.NewSolidBank(a.cfg.Banks.SolidBankURL)

	applicationRepository := repositories.NewApplicationRepository(a.db)
	offerRepository := repositories.NewOfferRepository(a.db)

	wsHandler := ws.NewWebSocketHandler(a.logger)
	defer wsHandler.CloseAll()

	applicationService := services.NewApplicationService(a.logger, []banks.Bank{fastBank, solidBank}, applicationRepository, offerRepository, wsHandler)
	applicationHandler := httpHandlers.NewApplicationHandler(applicationService)

	if err := a.registerCronJob("check offers", a.cfg.CronTabs.CheckOffersCronTab, applicationService.UpdateApplicationStatuses); err != nil {
		return fmt.Errorf("failed to register cron job: %v", err)
	}

	a.cron.Start()

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	if a.cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Use(controllers.AuthMiddleware())

	r.GET("/ws/applications/:id", wsHandler.SubscribeToApplicationUpdates)
	r.POST("/api/applications", applicationHandler.SubmitApplication)
	r.GET("/api/applications/:id", applicationHandler.GetApplication)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.Port),
		Handler: r,
	}

	return srv.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.cron.StopJobs(); err != nil {
		a.logger.Error("failed to stop cron jobs", zap.Error(err))
	}

	db, err := a.db.DB()
	if err != nil {
		a.logger.Error("failed to get db connection for closing", zap.Error(err))
	} else {
		if err := db.Close(); err != nil {
			a.logger.Error("failed to close db connection", zap.Error(err))
		}
	}

	return a.srv.Shutdown(ctx)
}
