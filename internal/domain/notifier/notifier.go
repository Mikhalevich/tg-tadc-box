package notifier

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/usecase/openbox"
)

var (
	_ openbox.Notifier = (*Notifier)(nil)
)

type Sender interface {
	SendMessage(ctx context.Context, msg msginfo.Message) error
}

type MarkdownEscaper interface {
	EscapeMarkdown(s string) string
}

type Notifier struct {
	sender  Sender
	escaper MarkdownEscaper
}

func New(
	sender Sender,
	escaper MarkdownEscaper,
) *Notifier {
	return &Notifier{
		sender:  sender,
		escaper: escaper,
	}
}
