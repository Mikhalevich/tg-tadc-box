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
			expected: "Rare 100 gloinks",
		},
		{
			name:     "large amount",
			boxType:  box.TypeLegendary,
			amount:   gloink.AmountFromInt(1000),
			expected: "Legendary 1000 gloinks",
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
