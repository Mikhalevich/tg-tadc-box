package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (n *Notifier) ShowReward(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	receivedReward reward.Reward,
) error {
	payload, err := n.imageProvider.Reward(ctx, receivedReward)
	if err != nil {
		return fmt.Errorf("receive image paylod: %w", err)
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:     chatID,
			ReplyMsgID: messageID,
			Type:       msginfo.MessageTypeEditPNG,
			Payload:    payload,
			Buttons: []button.ButtonRow{
				button.Row(box.GetBoxButton("Schedule next box", false)),
			},
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
