package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

func (n *Notifier) ShowReadyToOpenBox(
	ctx context.Context,
	domBox box.Box,
) error {
	if err := n.sendBoxIsAvailable(ctx, domBox); err != nil {
		return fmt.Errorf("send box is available: %w", err)
	}

	return nil
}
