package player_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func TestAddReward(t *testing.T) {
	t.Parallel()

	var (
		now = time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)

		newCard = func(
			rewardID reward.ID,
			count int,
			createdAt time.Time,
			updatedAt time.Time,
		) player.Card {
			return player.Card{
				RewardID:  rewardID,
				Count:     count,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}
		}
	)

	tests := map[string]struct {
		RewardType reward.RewardType
		RewardID   reward.ID
		CreatedAt  time.Time
		Cards      map[reward.RewardType][]player.Card
		CardsAfter map[reward.RewardType][]player.Card
	}{
		"nil cards": {
			RewardType: reward.RewardTypeCommon,
			RewardID:   reward.IDFromInt(1),
			CreatedAt:  now,
			CardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 1, now, now),
				},
			},
		},
		"empty cards": {
			RewardType: reward.RewardTypeCommon,
			RewardID:   reward.IDFromInt(2),
			CreatedAt:  now,
			Cards:      map[reward.RewardType][]player.Card{},
			CardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(2), 1, now, now),
				},
			},
		},
		"add a new common card": {
			RewardType: reward.RewardTypeCommon,
			RewardID:   reward.IDFromInt(2),
			CreatedAt:  now,
			Cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 1, now, now),
				},
			},
			CardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 1, now, now),
					newCard(reward.IDFromInt(2), 1, now, now),
				},
			},
		},
		"add existing common card": {
			RewardType: reward.RewardTypeCommon,
			RewardID:   reward.IDFromInt(1),
			CreatedAt:  now,
			Cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 1, now.Add(-1*time.Hour), now.Add(-1*time.Hour)),
				},
			},
			CardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 2, now.Add(-1*time.Hour), now),
				},
			},
		},
		"add a new rare card in existing common cards": {
			RewardType: reward.RewardTypeRare,
			RewardID:   reward.IDFromInt(2),
			CreatedAt:  now,
			Cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 1, now.Add(-1*time.Hour), now.Add(-1*time.Hour)),
				},
			},
			CardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 1, now.Add(-1*time.Hour), now.Add(-1*time.Hour)),
				},
				reward.RewardTypeRare: {
					newCard(reward.IDFromInt(2), 1, now, now),
				},
			},
		},
		"add existing rare card in existing common cards": {
			RewardType: reward.RewardTypeRare,
			RewardID:   reward.IDFromInt(2),
			CreatedAt:  now,
			Cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 2, now.Add(-1*time.Hour), now.Add(-1*time.Hour)),
				},
				reward.RewardTypeRare: {
					newCard(reward.IDFromInt(2), 2, now.Add(-1*time.Hour), now.Add(-1*time.Hour)),
				},
			},
			CardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 2, now.Add(-1*time.Hour), now.Add(-1*time.Hour)),
				},
				reward.RewardTypeRare: {
					newCard(reward.IDFromInt(2), 3, now.Add(-1*time.Hour), now),
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cards := player.CardsCollected{
				CardsByType: tt.Cards,
			}

			cards.AddReward(tt.RewardType, tt.RewardID, tt.CreatedAt)

			require.Equal(t, tt.CardsAfter, cards.CardsByType)
		})
	}
}
