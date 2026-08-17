package player_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

func TestCardsCollectedAddReward(t *testing.T) {
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

func TestCardsCollectedCardsCount(t *testing.T) {
	t.Parallel()

	var (
		now     = time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)
		newCard = func(id reward.ID, count int, createdAt time.Time) player.Card {
			return player.Card{
				RewardID:  id,
				Count:     count,
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			}
		}
	)

	tests := map[string]struct {
		Cards map[reward.RewardType][]player.Card
		Want  map[reward.RewardType]int
	}{
		"nil cards": {
			Cards: nil,
			Want: map[reward.RewardType]int{
				reward.RewardTypeCommon:    0,
				reward.RewardTypeRare:      0,
				reward.RewardTypeEpic:      0,
				reward.RewardTypeLegendary: 0,
			},
		},
		"empty cards": {
			Cards: map[reward.RewardType][]player.Card{},
			Want: map[reward.RewardType]int{
				reward.RewardTypeCommon:    0,
				reward.RewardTypeRare:      0,
				reward.RewardTypeEpic:      0,
				reward.RewardTypeLegendary: 0,
			},
		},
		"single common card": {
			Cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 1, now),
				},
			},
			Want: map[reward.RewardType]int{
				reward.RewardTypeCommon:    1,
				reward.RewardTypeRare:      0,
				reward.RewardTypeEpic:      0,
				reward.RewardTypeLegendary: 0,
			},
		},
		"duplicates counted as a single card entry": {
			Cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 2, now),
				},
			},
			Want: map[reward.RewardType]int{
				reward.RewardTypeCommon:    1,
				reward.RewardTypeRare:      0,
				reward.RewardTypeEpic:      0,
				reward.RewardTypeLegendary: 0,
			},
		},
		"common and rare cards": {
			Cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {
					newCard(reward.IDFromInt(1), 1, now),
					newCard(reward.IDFromInt(2), 1, now.Add(-1*time.Hour)),
				},
				reward.RewardTypeRare: {
					newCard(reward.IDFromInt(1), 1, now.Add(-2*time.Hour)),
				},
				reward.RewardTypeEpic: {
					newCard(reward.IDFromInt(1), 1, now.Add(-3*time.Hour)),
				},
				reward.RewardTypeLegendary: {
					newCard(reward.IDFromInt(1), 1, now.Add(-4*time.Hour)),
				},
			},
			Want: map[reward.RewardType]int{
				reward.RewardTypeCommon:    2,
				reward.RewardTypeRare:      1,
				reward.RewardTypeEpic:      1,
				reward.RewardTypeLegendary: 1,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cards := player.CardsCollected{
				CardsByType: tt.Cards,
			}

			got := cards.CardsCount()

			require.Equal(t, tt.Want, got)
		})
	}
}
