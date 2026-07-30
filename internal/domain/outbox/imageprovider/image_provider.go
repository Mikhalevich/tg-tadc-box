package imageprovider

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/outboxmsg/imagepayload"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type ImageProvider struct {
}

func New() *ImageProvider {
	return &ImageProvider{}
}

func (i *ImageProvider) Reward(
	ctx context.Context,
	rew reward.Reward,
) ([]byte, error) {
	payload := imagepayload.ImagePayload{
		Type:   imagepayload.PayloadTypeReward,
		Reward: rew,
	}

	buf, err := payload.GOBEncode()
	if err != nil {
		return nil, fmt.Errorf("encode payload: %w", err)
	}

	return buf, nil
}

func (i *ImageProvider) Chest(ctx context.Context) ([]byte, error) {
	payload := imagepayload.ImagePayload{
		Type: imagepayload.PayloadTypeCommonChest,
	}

	buf, err := payload.GOBEncode()
	if err != nil {
		return nil, fmt.Errorf("encode payload: %w", err)
	}

	return buf, nil
}
