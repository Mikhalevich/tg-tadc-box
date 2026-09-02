package like

import (
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type LikeButtonPayload struct {
	BoxID    box.ID
	RewardID reward.ID
	Type     Type
}

func LikeButton(
	caption string,
	boxID box.ID,
	rewardID reward.ID,
	likeType Type,
) (button.Button, error) {
	btn, err := button.CreateButton(
		caption,
		button.OperationLike,
		button.WithPayload(
			LikeButtonPayload{
				BoxID:    boxID,
				RewardID: rewardID,
				Type:     likeType,
			},
		),
	)
	if err != nil {
		return button.Button{}, fmt.Errorf("create button: %w", err)
	}

	return btn, nil
}
