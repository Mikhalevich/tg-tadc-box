package main

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/cmd/poller/internal/config"
	"github.com/Mikhalevich/tg-tadc-box/cmd/poller/internal/setup"
	"github.com/Mikhalevich/tg-tadc-box/internal/infra/application"
)

func main() {
	var cfg config.Config
	application.Run(&cfg, func(ctx context.Context) error {
		if err := setup.StartWorker(ctx, cfg); err != nil {
			return fmt.Errorf("start bot: %w", err)
		}

		return nil
	})
}
