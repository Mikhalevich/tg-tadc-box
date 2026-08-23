package boxsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func (s *Service) OpenBox(
	ctx context.Context,
	chatID msginfo.ChatID,
	boxID box.ID,
	receivedRewardID reward.ID,
	completedAt time.Time,
) error {
	if err := s.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := s.repo.CompleteBox(ctx, boxID, box.StatusOpened, completedAt); err != nil {
			return fmt.Errorf("complete box: %w", err)
		}

		if err := s.repo.InsertReceivedReward(ctx, reward.ReceivedReward{
			ChatID:    chatID,
			RewardID:  receivedRewardID,
			BoxID:     boxID,
			CreatedAt: completedAt,
		}); err != nil {
			return fmt.Errorf("insert received reward: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
