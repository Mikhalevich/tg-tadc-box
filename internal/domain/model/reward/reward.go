package reward

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

type ID int

func IDFromInt(id int) ID {
	return ID(id)
}

func (id ID) Int() int {
	return int(id)
}

type Reward struct {
	ID   ID
	Type RewardType
}
