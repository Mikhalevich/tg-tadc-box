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
	msg := fmt.Sprintf("Reward: %d - %s", receivedReward.ID.Int(), receivedReward.Type.String())
	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID: chatID,
			Text:   msg,
			Type:   msginfo.MessageTypePlain,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
