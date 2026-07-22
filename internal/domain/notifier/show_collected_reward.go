package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (n *Notifier) ShowCollectedReward(
	ctx context.Context,
	chatID msginfo.ChatID,
	rew reward.Reward,
	count int,
) error {
	payload, err := n.imageProvider.Reward(ctx, rew)
	if err != nil {
		return fmt.Errorf("receive image payload: %w", err)
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:  chatID,
			Type:    msginfo.MessageTypePNG,
			Payload: payload,
			Text:    fmt.Sprintf("x%d", count),
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
