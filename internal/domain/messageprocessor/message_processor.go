package messageprocessor

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

type ImageProvider interface {
	Image(ctx context.Context) ([]byte, error)
}

type MessageProcessor struct {
	sender           Sender
	escaper          MarkdownEscaper
	buttonRepository ButtonRepository
	imageProvider    ImageProvider
}

func New(
	sender Sender,
	escaper MarkdownEscaper,
	buttonRepository ButtonRepository,
	imageProvider ImageProvider,
) *MessageProcessor {
	return &MessageProcessor{
		sender:           sender,
		escaper:          escaper,
		buttonRepository: buttonRepository,
		imageProvider:    imageProvider,
	}
}
