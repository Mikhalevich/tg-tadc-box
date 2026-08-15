package notifier

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/like"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (n *Notifier) ShowReward(
	ctx context.Context,
	chatID msginfo.ChatID,
	messageID msginfo.MessageID,
	receivedReward reward.Reward,
	openingBox box.Box,
	withLikeButtons bool,
) error {
	payload, err := n.imageProvider.Reward(ctx, receivedReward)
	if err != nil {
		return fmt.Errorf("receive image paylod: %w", err)
	}

	buttons, err := makeShowRewardButtons(withLikeButtons, openingBox.ID, receivedReward.ID)
	if err != nil {
		return fmt.Errorf("make buttons: %w", err)
	}

	if err := n.sender.SendMessage(
		ctx,
		msginfo.Message{
			ChatID:     chatID,
			ReplyMsgID: messageID,
			Type:       msginfo.MessageTypePNG,
			Text:       makeBonusRewardDescription(openingBox.Meta),
			Payload:    payload,
			Buttons:    buttons,
		},
	); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func makeBonusRewardDescription(meta box.Meta) string {
	if !meta.BonusBox.IsValid {
		return ""
	}

	return fmt.Sprintf(
		"Reward for opening x%d *%s* boxes",
		meta.BonusBox.Attempts,
		meta.BonusBox.Type.String(),
	)
}

func makeShowRewardButtons(
	withLikeButtons bool,
	boxID box.ID,
	rewardID reward.ID,
) ([]button.ButtonRow, error) {
	shopBtn := box.ShopButton("Get next box")

	if !withLikeButtons {
		return []button.ButtonRow{
			button.Row(shopBtn),
		}, nil
	}

	likeBtn, err := like.LikeButton("👍", boxID, rewardID, like.TypeLike)
	if err != nil {
		return nil, fmt.Errorf("create like button: %w", err)
	}

	dislikeBtn, err := like.LikeButton("👎", boxID, rewardID, like.TypeDislike)
	if err != nil {
		return nil, fmt.Errorf("crate dislike button: %w", err)
	}

	return []button.ButtonRow{
		button.Row(likeBtn, dislikeBtn),
		button.Row(shopBtn),
	}, nil
}
