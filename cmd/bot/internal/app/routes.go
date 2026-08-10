package app

import (
	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tgbot"
	"github.com/Mikhalevich/tg-tadc-box/cmd/bot/internal/app/tghandler"
)

func makeRoutes(tbot *tgbot.TGBot, handler *tghandler.TGHandler) {
	tbot.AddMiddleware(handler.ErrorMiddleware)

	tbot.AddTextCommand("start", handler.Start)

	tbot.AddMenuCommand("get_box", "🎁 view boxes for purchase", handler.Shop)
	tbot.AddMenuCommand("my_cards", "👀 view my cards", handler.MyCards)

	tbot.AddDefaultHandler(handler.DefaultHandler)
	tbot.AddDefaultCallbackQueryHander(handler.DefaultCallbackQuery)
}
