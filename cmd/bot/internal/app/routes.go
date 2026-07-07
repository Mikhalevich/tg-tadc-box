package app

import (
	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tghandler"
)

func makeRoutes(tbot *tgbot.TGBot, handler *tghandler.TGHandler) {
	tbot.AddMiddleware(handler.ErrorMiddleware)

	tbot.AddTextCommand("start", handler.Start)

	tbot.AddDefaultHandler(handler.DefaultHandler)
	tbot.AddDefaultCallbackQueryHander(handler.DefaultCallbackQuery)
}
