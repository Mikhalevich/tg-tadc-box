package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/like"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type cbHandler func(ctx context.Context, msg tgbot.BotMessage, btn *button.Button) error

type ButtonProvider interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
}

type BoxOpenByID interface {
	OpenByID(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		id box.ID,
	) error
}

type BoxShowReadyToOpen interface {
	ShowReadyToOpenBox(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxID box.ID,
	) error
}

type CardViewer interface {
	Total(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
	Page(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		rewardType reward.RewardType,
		page int,
	) error
}

type CardAbstracter interface {
	All(ctx context.Context, chatID msginfo.ChatID) error
	Abstract(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		rewardType reward.RewardType,
		pos int,
		count int,
	) error
}

type Shop interface {
	ViewBoxes(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
	BuyBox(
		ctx context.Context,
		chatID msginfo.ChatID,
		boxType box.Type,
	) error
	BuyCommonBox(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
}

type LikeProcessor interface {
	Like(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		boxID box.ID,
		rewardID reward.ID,
		likeType like.Type,
	) error
}

type SendInviteLink interface {
	SendInvite(
		ctx context.Context,
		chatID msginfo.ChatID,
	) error
}

type Notifier interface {
	Welcome(ctx context.Context, chatID msginfo.ChatID) error
}

type ErrorNotifier interface {
	ParseError(ctx context.Context, chatID msginfo.ChatID, err error) error
}

type TGHandler struct {
	cbHanlers      map[button.Operation]cbHandler
	buttonProvider ButtonProvider
	boxOpenByID    BoxOpenByID
	boxShowReady   BoxShowReadyToOpen
	cardViewer     CardViewer
	cardAbstracter CardAbstracter
	shop           Shop
	likeProcessor  LikeProcessor
	inviteSender   SendInviteLink
	notifier       Notifier
	errorNotifier  ErrorNotifier
}

func New(
	buttonProvider ButtonProvider,
	boxOpenByID BoxOpenByID,
	boxShowReady BoxShowReadyToOpen,
	cardViewer CardViewer,
	cardAbstracter CardAbstracter,
	shop Shop,
	likeProcessor LikeProcessor,
	inviteSender SendInviteLink,
	notifier Notifier,
	errorNotifier ErrorNotifier,
) *TGHandler {
	tgh := &TGHandler{
		buttonProvider: buttonProvider,
		boxOpenByID:    boxOpenByID,
		boxShowReady:   boxShowReady,
		cardViewer:     cardViewer,
		cardAbstracter: cardAbstracter,
		shop:           shop,
		likeProcessor:  likeProcessor,
		inviteSender:   inviteSender,
		notifier:       notifier,
		errorNotifier:  errorNotifier,
	}

	tgh.registerCBHandlers()

	return tgh
}

func (t *TGHandler) registerCBHandlers() {
	t.cbHanlers = map[button.Operation]cbHandler{
		button.OperationGetCommonBox:          t.cbGetCommonBox,
		button.OperationOpenBox:               t.cbOpenBox,
		button.OperationBoxReadyToOpen:        t.cbReadyToOpenBox,
		button.OperationCardPage:              t.cbCollectedCardPage,
		button.OperationCardTotal:             t.cbCollectedCardTotalPage,
		button.OperationAbstractDuplicatesAll: t.cbAbstractDuplicatedAll,
		button.OperationAbstractCard:          t.cbAbstractCard,
		button.OperationBuyBox:                t.cbBuyBox,
		button.OperationShop:                  t.cbShop,
		button.OperationLike:                  t.cbLikeReward,
	}
}
