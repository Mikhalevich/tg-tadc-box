package tghandler

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type cbHandler func(ctx context.Context, msg tgbot.BotMessage, btn *button.Button) error

type ButtonProvider interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
}

type BoxProcessor interface {
	Open(ctx context.Context, chatID msginfo.ChatID) error
	OpenByID(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
		id int,
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
	boxProcessor   BoxProcessor
	cardViewer     CardViewer
	cardAbstracter CardAbstracter
	shop           Shop
	notifier       Notifier
	errorNotifier  ErrorNotifier
}

func New(
	buttonProvider ButtonProvider,
	boxProcessor BoxProcessor,
	cardViewer CardViewer,
	cardAbstracter CardAbstracter,
	shop Shop,
	notifier Notifier,
	errorNotifier ErrorNotifier,
) *TGHandler {
	tgh := &TGHandler{
		buttonProvider: buttonProvider,
		boxProcessor:   boxProcessor,
		cardViewer:     cardViewer,
		cardAbstracter: cardAbstracter,
		shop:           shop,
		notifier:       notifier,
		errorNotifier:  errorNotifier,
	}

	tgh.registerCBHandlers()

	return tgh
}

func (t *TGHandler) registerCBHandlers() {
	t.cbHanlers = map[button.Operation]cbHandler{
		button.OperationGetBox:                t.cbGetBox,
		button.OperationOpenBox:               t.cbOpenBox,
		button.OperationCardPage:              t.cbCollectedCardPage,
		button.OperationCardTotal:             t.cbCollectedCardTotalPage,
		button.OperationAbstractDuplicatesAll: t.cbAbstractDuplicatedAll,
		button.OperationBuyBox:                t.cbBuyBox,
	}
}
