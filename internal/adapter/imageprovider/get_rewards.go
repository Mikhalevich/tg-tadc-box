package imageprovider

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

const (
	rewardFolderTemplate = "assets/%s"

	rewardTypes = 4
)

func (i *ImageProvider) GetRewards() (map[reward.RewardType][]reward.ID, error) {
	rewards := make(map[reward.RewardType][]reward.ID, rewardTypes)

	rewardTypes := [...]reward.RewardType{
		reward.RewardTypeCommon,
		reward.RewardTypeRare,
		reward.RewardTypeEpic,
		reward.RewardTypeLegendary,
	}

	for _, rewardType := range rewardTypes {
		rewardFolder := fmt.Sprintf(rewardFolderTemplate, folderByRewardType(rewardType))
		rewardIDs, err := collectRewards(rewardFolder)
		if err != nil {
			return nil, fmt.Errorf("collect rewards in %s: %w", rewardFolder, err)
		}

		rewards[rewardType] = rewardIDs
	}

	return rewards, nil
}

func collectRewards(folder string) ([]reward.ID, error) {
	entries, err := assetsFS.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	rewardIDs := make([]reward.ID, 0, len(entries))

	for _, entry := range entries {
		rawID := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

		rewardID, err := strconv.Atoi(rawID)
		if err != nil {
			return nil, fmt.Errorf("parse reward id: %w", err)
		}

		rewardIDs = append(rewardIDs, reward.IDFromInt(rewardID))
	}

	return rewardIDs, nil
}
