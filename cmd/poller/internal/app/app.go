package app

import (
	"context"
	"sync"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/cmd/poller/internal/config"
	"github.com/Mikhalevich/tg-tadc-box/internal/infra/logger"
)

type BoxReadyNotification interface {
	SendNotifications(ctx context.Context, limit int) error
}

type App struct {
	boxReadyNotification BoxReadyNotification
}

func New(
	boxReadyNotification BoxReadyNotification,
) *App {
	return &App{
		boxReadyNotification: boxReadyNotification,
	}
}

func (a *App) Run(
	ctx context.Context,
	boxReadyNotificationCfg config.Worker,
) {
	var wgr sync.WaitGroup

	runWorkers(
		ctx,
		"box ready to open notifications",
		boxReadyNotificationCfg.Count,
		boxReadyNotificationCfg.Interval,
		&wgr,
		func(ctx context.Context) error {
			return a.boxReadyNotification.SendNotifications(
				ctx,
				boxReadyNotificationCfg.BatchSize,
			)
		},
	)

	wgr.Wait()
}

func runWorkers(
	ctx context.Context,
	workerName string,
	workersCount int,
	pollerInterval time.Duration,
	wgr *sync.WaitGroup,
	processFn func(ctx context.Context) error,
) {
	for i := range workersCount {
		wgr.Go(func() {
			log := logger.FromContext(ctx).
				WithFields(logger.Fields{
					"worker_name":   workerName,
					"worker_number": i,
				})
			runPoller(
				logger.WithLogger(ctx, log),
				pollerInterval,
				processFn,
			)
		})
	}
}

func runPoller(
	ctx context.Context,
	interval time.Duration,
	processFn func(ctx context.Context) error,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := processFn(ctx); err != nil {
				logger.FromContext(ctx).
					WithError(err).
					Error("process error")
			}

		case <-ctx.Done():
			return
		}
	}
}
