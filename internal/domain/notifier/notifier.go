package notifier

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
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

type ImageProvider interface {
	Chest(ctx context.Context) ([]byte, error)
	Reward(ctx context.Context, rew reward.Reward) ([]byte, error)
}

type Notifier struct {
	sender        Sender
	escaper       MarkdownEscaper
	imageProvider ImageProvider
}

func New(
	sender Sender,
	escaper MarkdownEscaper,
	imageProvider ImageProvider,
) *Notifier {
	return &Notifier{
		sender:        sender,
		escaper:       escaper,
		imageProvider: imageProvider,
	}
}
