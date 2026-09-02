package notificationsvc

import (
	"context"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/reward"
)

type Sender interface {
	SendMessage(ctx context.Context, msg msginfo.Message) error
}

type MarkdownEscaper interface {
	EscapeMarkdown(s string) string
}

type ImageProvider interface {
	Chest(ctx context.Context, boxType box.Type) ([]byte, error)
	Reward(ctx context.Context, rew reward.Reward) ([]byte, error)
}

type Service struct {
	botName       string
	sender        Sender
	escaper       MarkdownEscaper
	imageProvider ImageProvider
}

func New(
	botName string,
	sender Sender,
	escaper MarkdownEscaper,
	imageProvider ImageProvider,
) *Service {
	return &Service{
		botName:       botName,
		sender:        sender,
		escaper:       escaper,
		imageProvider: imageProvider,
	}
}
