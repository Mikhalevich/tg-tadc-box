package reward

import (
	"fmt"
	"time"
)

type RewardType string

const (
	RewardTypeCommon    RewardType = "common"
	RewardTypeRare      RewardType = "rare"
	RewardTypeEpic      RewardType = "epic"
	RewardTypeLegendary RewardType = "legendary"
)

func (rt RewardType) String() string {
	return string(rt)
}

func TypeFromString(raw string) (RewardType, error) {
	switch raw {
	case RewardTypeCommon.String(),
		RewardTypeRare.String(),
		RewardTypeEpic.String(),
		RewardTypeLegendary.String():
		return RewardType(raw), nil
	}

	return RewardTypeCommon, fmt.Errorf("invalid string reward type: %q", raw)
}

type ID int

func IDFromInt(id int) ID {
	return ID(id)
}

func (id ID) Int() int {
	return int(id)
}

type Reward struct {
	ID        ID
	Type      RewardType
	URI       string
	CreatedAt time.Time
}
