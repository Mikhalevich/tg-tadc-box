package tghandler

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/card"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/like"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/perror"
)

func (t *TGHandler) DefaultCallbackQuery(ctx context.Context, msg tgbot.BotMessage, sender tgbot.MessageSender) error {
	if msg.Data == "" {
		return nil
	}

	btn, err := t.buttonProvider.GetButton(ctx, button.IDFromString(msg.Data))
	if err != nil {
		if perror.IsType(err, perror.TypeNotFound) {
			sender.SendMessage(ctx, msg.ChatID, "Button expired")

			return nil
		}

		return fmt.Errorf("get button: %w", err)
	}

	hndlr, ok := t.cbHanlers[btn.Operation]
	if !ok {
		return fmt.Errorf("invalid operation %v", btn.Operation)
	}

	if err := hndlr(ctx, msg, btn); err != nil {
		return fmt.Errorf("process cb handler: %w", err)
	}

	if btn.IsDeleteAfterProcess {
		sender.DeleteMessage(ctx, msg.ChatID, msg.MessageID)
	}

	return nil
}

func (t *TGHandler) cbGetCommonBox(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	if err := t.boxProcessor.OpenCommonBox(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
	); err != nil {
		return fmt.Errorf("open common box: %w", err)
	}

	return nil
}

func (t *TGHandler) cbOpenBox(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[box.OpenBoxButtonPayload](*btn)
	if err != nil {
		return fmt.Errorf("get open box payload: %w", err)
	}

	if err := t.boxProcessor.OpenByID(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
		payload.ID,
	); err != nil {
		return fmt.Errorf("open box by id: %w", err)
	}

	return nil
}

func (t *TGHandler) cbCollectedCardPage(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[card.PageButtonPayload](*btn)
	if err != nil {
		return fmt.Errorf("get collected cards payload: %w", err)
	}

	if err := t.cardViewer.Page(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
		payload.Type,
		payload.Page,
	); err != nil {
		return fmt.Errorf("view page: %w", err)
	}

	return nil
}

func (t *TGHandler) cbCollectedCardTotalPage(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	if err := t.cardViewer.Total(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
	); err != nil {
		return fmt.Errorf("total page: %w", err)
	}

	return nil
}

func (t *TGHandler) cbAbstractDuplicatedAll(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	if err := t.cardAbstracter.All(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
	); err != nil {
		return fmt.Errorf("abstract all: %w", err)
	}

	return nil
}

func (t *TGHandler) cbAbstractCard(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[card.AbstractCardPayload](*btn)
	if err != nil {
		return fmt.Errorf("get abstract card payload: %w", err)
	}

	if err := t.cardAbstracter.Abstract(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
		payload.Type,
		payload.Pos,
		payload.Count,
	); err != nil {
		return fmt.Errorf("abstract card: %w", err)
	}

	return nil
}

func (t *TGHandler) cbBuyBox(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[gloink.BuyBoxButtonPayload](*btn)
	if err != nil {
		return fmt.Errorf("get buy box payload: %w", err)
	}

	if err := t.shop.BuyBox(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		payload.BoxType,
	); err != nil {
		return fmt.Errorf("buy box: %w", err)
	}

	return nil
}

func (t *TGHandler) cbShop(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	if err := t.shop.ViewBoxes(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
	); err != nil {
		return fmt.Errorf("view boxes: %w", err)
	}

	return nil
}

func (t *TGHandler) cbLikeReward(
	ctx context.Context,
	msg tgbot.BotMessage,
	btn *button.Button,
) error {
	payload, err := button.GetPayload[like.LikeButtonPayload](*btn)
	if err != nil {
		return fmt.Errorf("get like button payload: %w", err)
	}

	if err := t.likeProcessor.Like(
		ctx,
		msginfo.ChatIDFromInt64(msg.ChatID),
		msginfo.MessageIDFromInt(msg.MessageID),
		payload.BoxID,
		payload.RewardID,
		payload.Type,
	); err != nil {
		return fmt.Errorf("process like: %w", err)
	}

	return nil
}
