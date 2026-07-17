package imageprovider

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

const (
	imageTemplate = "assets/%s/%d.jpg"
	chestPath     = "assets/chest/1.jpg"

	commonFolder    = "common"
	rareFolder      = "rare"
	epicFolder      = "epic"
	legendaryFolder = "legendary"
)

func (i *ImageProvider) Reward(
	ctx context.Context,
	rew reward.Reward,
) ([]byte, error) {
	payload, err := assetsFS.ReadFile(imageAbsPath(rew))
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

func imageAbsPath(rew reward.Reward) string {
	var (
		rewardFolder = folderByRewardType(rew.Type)
		rewardImage  = fmt.Sprintf(imageTemplate, rewardFolder, rew.ID.Int())
	)

	return rewardImage
}

func folderByRewardType(rewardType reward.RewardType) string {
	switch rewardType {
	case reward.RewardTypeCommon:
		return commonFolder

	case reward.RewardTypeRare:
		return rareFolder

	case reward.RewardTypeEpic:
		return epicFolder

	case reward.RewardTypeLegendary:
		return legendaryFolder
	}

	return ""
}
