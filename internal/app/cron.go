package app

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (a *App) registerCronJob(name, schedule string, task func(context.Context)) error {
	if schedule == "" {
		a.logger.Warn("cron job disabled due missing schedule", zap.String("name", name))
		return nil
	}

	_, err := a.cron.NewJob(
		gocron.CronJob(schedule, true),
		gocron.NewTask(a.taskWrapper(name, task)),
		gocron.WithName(name),
	)
	if err != nil {
		return errors.Errorf("error register %s masterjob: %s", name, err.Error())
	}

	return nil
}

func (a *App) taskWrapper(name string, task func(context.Context)) func() {
	return func() {
		a.logger.Debug("cron job started", zap.String("name", name))
		task(context.Background())
	}
}
