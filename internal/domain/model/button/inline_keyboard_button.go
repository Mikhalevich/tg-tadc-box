package button

type InlineKeyboardButton struct {
	ID      ID
	Caption string
	Style   Style
}

type InlineKeyboardButtonRow []InlineKeyboardButton

func InlineRow(buttons ...InlineKeyboardButton) InlineKeyboardButtonRow {
	return buttons
}
