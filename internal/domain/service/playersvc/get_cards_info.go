package playersvc

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type CardsInfo struct {
	Cards                        map[reward.RewardType]int
	AbstractionDuplicatesAllCost gloink.Amount
}

func (s *Service) GetCardsInfo(
	ctx context.Context,
	chatID msginfo.ChatID,
) (CardsInfo, error) {
	profile, err := s.getOrCreatePlayer(ctx, chatID)
	if err != nil {
		return CardsInfo{}, fmt.Errorf("get or create player: %w", err)
	}

	return CardsInfo{
		Cards:                        profile.Profile.Cards.CardsCount(),
		AbstractionDuplicatesAllCost: profile.Profile.Cards.ViewCostOfAbstractionDuplicatesAll(s.abstractionCosts),
	}, nil
}
