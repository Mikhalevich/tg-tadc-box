package imageprovider

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

const (
	chestPath = "assets/chest/1.jpg"
)

func (i *ImageProvider) Reward(
	ctx context.Context,
	rew reward.Reward,
) ([]byte, error) {
	payload, err := assetsFS.ReadFile(rew.URI)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return payload, nil
}

func (i *ImageProvider) Chest(ctx context.Context) ([]byte, error) {
	payload, err := assetsFS.ReadFile(chestPath)
	if err != nil {
		return nil, fmt.Errorf("read chest file: %w", err)
	}

	return payload, nil
}
