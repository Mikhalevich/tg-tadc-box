package player_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

//nolint:maintidx
func TestProfileAbstractByPos(t *testing.T) {
	t.Parallel()

	var (
		now = time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)

		makeCard = func(id, count int) player.Card {
			return player.Card{
				RewardID:  reward.IDFromInt(id),
				Count:     count,
				CreatedAt: now,
				UpdatedAt: now,
			}
		}

		abstractionCosts = map[reward.RewardType]gloink.Amount{
			reward.RewardTypeCommon:    gloink.AmountFromInt(10),
			reward.RewardTypeRare:      gloink.AmountFromInt(20),
			reward.RewardTypeEpic:      gloink.AmountFromInt(30),
			reward.RewardTypeLegendary: gloink.AmountFromInt(40),
		}
	)

	tests := map[string]struct {
		cards       map[reward.RewardType][]player.Card
		wallet      gloink.Amount
		rewardType  reward.RewardType
		pos         int
		count       int
		wantAmount  gloink.Amount
		wantWallet  gloink.Amount
		cardsAfter  map[reward.RewardType][]player.Card
		wantErr     bool
		wantErrType perror.Type
	}{
		"abstract single card": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:     gloink.AmountFromInt(100),
			rewardType: reward.RewardTypeCommon,
			pos:        1,
			count:      1,
			wantAmount: gloink.AmountFromInt(10),
			wantWallet: gloink.AmountFromInt(110),
			cardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 4)},
			},
		},
		"abstract multiple cards": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:     gloink.AmountFromInt(100),
			rewardType: reward.RewardTypeCommon,
			pos:        1,
			count:      3,
			wantAmount: gloink.AmountFromInt(30),
			wantWallet: gloink.AmountFromInt(130),
			cardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 2)},
			},
		},
		"abstract maximum allowed count keeps one card": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:     gloink.AmountFromInt(100),
			rewardType: reward.RewardTypeCommon,
			pos:        1,
			count:      4,
			wantAmount: gloink.AmountFromInt(40),
			wantWallet: gloink.AmountFromInt(140),
			cardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 1)},
			},
		},
		"abstract from second position": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 2), makeCard(2, 5)},
			},
			wallet:     gloink.AmountFromInt(100),
			rewardType: reward.RewardTypeCommon,
			pos:        2,
			count:      2,
			wantAmount: gloink.AmountFromInt(20),
			wantWallet: gloink.AmountFromInt(120),
			cardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 2), makeCard(2, 3)},
			},
		},
		"abstract with zero initial wallet": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:     gloink.AmountFromInt(0),
			rewardType: reward.RewardTypeCommon,
			pos:        1,
			count:      1,
			wantAmount: gloink.AmountFromInt(10),
			wantWallet: gloink.AmountFromInt(10),
			cardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 4)},
			},
		},
		"abstract rare card uses rare cost": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeRare: {makeCard(3, 5)},
			},
			wallet:     gloink.AmountFromInt(100),
			rewardType: reward.RewardTypeRare,
			pos:        1,
			count:      1,
			wantAmount: gloink.AmountFromInt(20),
			wantWallet: gloink.AmountFromInt(120),
			cardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeRare: {makeCard(3, 4)},
			},
		},
		"abstract epic card uses epic cost": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeEpic: {makeCard(4, 5)},
			},
			wallet:     gloink.AmountFromInt(100),
			rewardType: reward.RewardTypeEpic,
			pos:        1,
			count:      1,
			wantAmount: gloink.AmountFromInt(30),
			wantWallet: gloink.AmountFromInt(130),
			cardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeEpic: {makeCard(4, 4)},
			},
		},
		"abstract legendary card uses legendary cost": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeLegendary: {makeCard(5, 5)},
			},
			wallet:     gloink.AmountFromInt(100),
			rewardType: reward.RewardTypeLegendary,
			pos:        1,
			count:      1,
			wantAmount: gloink.AmountFromInt(40),
			wantWallet: gloink.AmountFromInt(140),
			cardsAfter: map[reward.RewardType][]player.Card{
				reward.RewardTypeLegendary: {makeCard(5, 4)},
			},
		},
		"zero count returns error": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:      gloink.AmountFromInt(100),
			rewardType:  reward.RewardTypeCommon,
			pos:         1,
			count:       0,
			wantErr:     true,
			wantErrType: perror.TypeInvalidParam,
		},
		"negative count returns error": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:      gloink.AmountFromInt(100),
			rewardType:  reward.RewardTypeCommon,
			pos:         1,
			count:       -1,
			wantErr:     true,
			wantErrType: perror.TypeInvalidParam,
		},
		"count equal to card count returns error": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:      gloink.AmountFromInt(100),
			rewardType:  reward.RewardTypeCommon,
			pos:         1,
			count:       5,
			wantErr:     true,
			wantErrType: perror.TypeInvalidParam,
		},
		"count greater than card count returns error": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:      gloink.AmountFromInt(100),
			rewardType:  reward.RewardTypeCommon,
			pos:         1,
			count:       6,
			wantErr:     true,
			wantErrType: perror.TypeInvalidParam,
		},
		"zero position returns error": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:      gloink.AmountFromInt(100),
			rewardType:  reward.RewardTypeCommon,
			pos:         0,
			count:       1,
			wantErr:     true,
			wantErrType: perror.TypeInvalidParam,
		},
		"negative position returns error": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:      gloink.AmountFromInt(100),
			rewardType:  reward.RewardTypeCommon,
			pos:         -1,
			count:       1,
			wantErr:     true,
			wantErrType: perror.TypeInvalidParam,
		},
		"position greater than cards count returns error": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5), makeCard(2, 5)},
			},
			wallet:      gloink.AmountFromInt(100),
			rewardType:  reward.RewardTypeCommon,
			pos:         3,
			count:       1,
			wantErr:     true,
			wantErrType: perror.TypeInvalidParam,
		},
		"no cards for reward type returns error": {
			cards: map[reward.RewardType][]player.Card{
				reward.RewardTypeCommon: {makeCard(1, 5)},
			},
			wallet:      gloink.AmountFromInt(100),
			rewardType:  reward.RewardTypeEpic,
			pos:         1,
			count:       1,
			wantErr:     true,
			wantErrType: perror.TypeInvalidParam,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			profile := &player.Profile{
				Cards: player.CardsCollected{
					CardsByType: tt.cards,
				},
				Wallet: player.Wallet{
					GloinksAmount: tt.wallet,
				},
			}

			gotAmount, err := profile.AbstractByPos(tt.rewardType, tt.pos, tt.count, abstractionCosts)

			if tt.wantErr {
				require.Error(t, err)
				require.True(t, perror.IsType(err, tt.wantErrType),
					"expected error type %v, got: %v", tt.wantErrType, err)
				require.Equal(t, gloink.Amount(0), gotAmount)
				require.Equal(t, tt.wallet, profile.Wallet.GloinksAmount,
					"wallet should be unchanged on error")
				require.Equal(t, tt.cards, profile.Cards.CardsByType,
					"cards should be unchanged on error")

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantAmount, gotAmount)
			require.Equal(t, tt.wantWallet, profile.Wallet.GloinksAmount)
			require.Equal(t, tt.cardsAfter, profile.Cards.CardsByType)
		})
	}
}
