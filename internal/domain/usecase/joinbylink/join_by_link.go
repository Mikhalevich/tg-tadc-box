package joinbylink

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/player"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/referral"
	"github.com/Mikhalevich/tg-tadc-box/internal/infra/logger"
)

const (
	referralJoinGloinksReward = 100
)

type Transactor interface {
	Transaction(ctx context.Context, trxFn func(ctx context.Context) error) error
}

type PlayerService interface {
	GetPlayerByChatID(
		ctx context.Context,
		chatID msginfo.ChatID,
	) (player.Player, bool, error)
	AddGloinks(
		ctx context.Context,
		chatID msginfo.ChatID,
		amount gloink.Amount,
	) error
}

type ReferralService interface {
	GetReferralByCode(ctx context.Context, code referral.Code) (referral.Referral, error)
	InsertReferral(ctx context.Context, ref referral.Referral) error
}

type Notifier interface {
	Welcome(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
	ReferralLinkActivated(
		ctx context.Context,
		chatID msginfo.ChatID,
		gloinksReceived gloink.Amount,
	) error
}

type TimeProvider interface {
	Now() time.Time
}

type JoinByLink struct {
	transactor      Transactor
	playerService   PlayerService
	referralService ReferralService
	notifier        Notifier
	timeProvider    TimeProvider
}

func New(
	transactor Transactor,
	playerService PlayerService,
	referralService ReferralService,
	notifier Notifier,
	timeProvider TimeProvider,
) *JoinByLink {
	return &JoinByLink{
		transactor:      transactor,
		playerService:   playerService,
		referralService: referralService,
		notifier:        notifier,
		timeProvider:    timeProvider,
	}
}

func (j *JoinByLink) Join(
	ctx context.Context,
	refereeChatID msginfo.ChatID,
	code referral.Code,
) error {
	if err := j.processReferral(ctx, refereeChatID, code); err != nil {
		logger.FromContext(ctx).WithError(err).Error("process referral")
		// skip error
	}

	if err := j.notifier.Welcome(ctx, refereeChatID); err != nil {
		return fmt.Errorf("welcome message: %w", err)
	}

	return nil
}

func (j *JoinByLink) processReferral(
	ctx context.Context,
	refereeChatID msginfo.ChatID,
	code referral.Code,
) error {
	if code.IsValid() {
		return perror.InvalidParam("referral code is not valid")
	}
	_, isNew, err := j.playerService.GetPlayerByChatID(ctx, refereeChatID)
	if err != nil {
		return fmt.Errorf("get player by chat id: %w", err)
	}

	if !isNew {
		return nil
	}

	ref, err := j.referralService.GetReferralByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("get referral by code: %w", err)
	}

	if ref.ChatID == refereeChatID {
		return nil
	}

	if err := j.transactor.Transaction(ctx, func(ctx context.Context) error {
		if err := j.referralService.InsertReferral(
			ctx,
			referral.Referral{
				ChatID:    refereeChatID,
				InvitedBy: ref.ChatID,
				Code:      referral.GenerateCode(),
				CreatedAt: j.timeProvider.Now(),
			},
		); err != nil {
			return fmt.Errorf("insert referral: %w", err)
		}

		gloinkReward := gloink.AmountFromInt(referralJoinGloinksReward)

		if err := j.playerService.AddGloinks(
			ctx,
			ref.ChatID,
			gloinkReward,
		); err != nil {
			return fmt.Errorf("add referral gloinks: %w", err)
		}

		if err := j.notifier.ReferralLinkActivated(ctx, ref.ChatID, gloinkReward); err != nil {
			return fmt.Errorf("referral link activated: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
