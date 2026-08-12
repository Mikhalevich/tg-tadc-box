//nolint:testpackage
package notifier

import (
	"testing"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/gloink"
)

func TestMessageForBoxAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		boxType  box.Type
		amount   gloink.Amount
		expected string
	}{
		{
			name:     "zero amount returns Free",
			boxType:  box.TypeCommon,
			amount:   gloink.AmountFromInt(0),
			expected: "Free",
		},
		{
			name:     "positive amount",
			boxType:  box.TypeRare,
			amount:   gloink.AmountFromInt(100),
			expected: "rare 100 gloinks",
		},
		{
			name:     "large amount",
			boxType:  box.TypeLegendary,
			amount:   gloink.AmountFromInt(1000),
			expected: "legendary 1000 gloinks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := messageForBoxAmount(tt.boxType, tt.amount)

			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestMessageBoxAlreadyInProgress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		msg          string
		isInProgress bool
		expected     string
	}{
		{
			name:         "not in progress returns message unchanged",
			msg:          "common 100 gloinks",
			isInProgress: false,
			expected:     "common 100 gloinks",
		},
		{
			name:         "in progress appends marker",
			msg:          "common 100 gloinks",
			isInProgress: true,
			expected:     "common 100 gloinks (⌛️)",
		},
		{
			name:         "in progress with free message",
			msg:          "Free",
			isInProgress: true,
			expected:     "Free (⌛️)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := messageBoxAlreadyInProgress(tt.msg, tt.isInProgress)

			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
