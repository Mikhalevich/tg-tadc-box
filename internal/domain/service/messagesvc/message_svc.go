package messagesvc

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/button"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

type Sender interface {
	SendMessage(
		ctx context.Context,
		msg msginfo.SenderMessage,
	) error
	DeleteMessage(
		ctx context.Context,
		chatID msginfo.ChatID,
		messageID msginfo.MessageID,
	) error
}

type MarkdownEscaper interface {
	EscapeMarkdown(s string) string
}

type ButtonRepository interface {
	GetButton(ctx context.Context, id button.ID) (*button.Button, error)
	SetButtonRows(ctx context.Context, rows ...button.ButtonRow) error
}

type Service struct {
	sender           Sender
	escaper          MarkdownEscaper
	buttonRepository ButtonRepository
}

func New(
	sender Sender,
	escaper MarkdownEscaper,
	buttonRepository ButtonRepository,
) *Service {
	return &Service{
		sender:           sender,
		escaper:          escaper,
		buttonRepository: buttonRepository,
	}
}
