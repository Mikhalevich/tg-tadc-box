package imageprovider

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

const (
	commonPath    = "assets/chest/common.jpg"
	rarePath      = "assets/chest/rare.jpg"
	epicPath      = "assets/chest/epic.jpg"
	legendaryPath = "assets/chest/legendary.jpg"
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

func (i *ImageProvider) Chest(
	ctx context.Context,
	boxType box.Type,
) ([]byte, error) {
	payload, err := assetsFS.ReadFile(chestPathByType(boxType))
	if err != nil {
		return nil, fmt.Errorf("read chest file: %w", err)
	}

	return payload, nil
}

func chestPathByType(boxType box.Type) string {
	switch boxType {
	case box.TypeCommon:
		return commonPath

	case box.TypeRare:
		return rarePath

	case box.TypeEpic:
		return epicPath

	case box.TypeLegendary:
		return legendaryPath
	}

	return commonPath
}
