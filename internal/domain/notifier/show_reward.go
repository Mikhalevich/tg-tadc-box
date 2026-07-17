package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (n *Notifier) ShowReward(
	ctx context.Context,
	chatID msginfo.ChatID,
	receivedReward reward.Reward,
) error {
	payload, err := n.imageProvider.Reward(ctx, receivedReward)
	if err != nil {
		return fmt.Errorf("receive image paylod: %w", err)
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:  chatID,
			Type:    msginfo.MessageTypePNG,
			Payload: payload,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
