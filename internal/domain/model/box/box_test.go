package box_test

import (
	"testing"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
)

func TestToMapByType(t *testing.T) {
	t.Parallel()

	t.Run("empty boxes returns nil", func(t *testing.T) {
		t.Parallel()

		result := box.ToMapByType(nil)

		if result != nil {
			t.Fatalf("expected nil map, got %v", result)
		}
	})

	t.Run("single box", func(t *testing.T) {
		t.Parallel()

		originBox := box.Box{ID: 1, Type: box.TypeCommon}

		result := box.ToMapByType([]box.Box{originBox})

		if len(result) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(result))
		}

		if got, ok := result[box.TypeCommon]; !ok {
			t.Fatalf("expected type %q to be present", box.TypeCommon)
		} else if got != originBox {
			t.Fatalf("expected box %+v, got %+v", originBox, got)
		}
	})

	t.Run("multiple boxes", func(t *testing.T) {
		t.Parallel()

		boxes := []box.Box{
			{ID: 1, Type: box.TypeCommon},
			{ID: 2, Type: box.TypeRare},
			{ID: 3, Type: box.TypeEpic},
			{ID: 4, Type: box.TypeLegendary},
		}

		result := box.ToMapByType(boxes)

		if len(result) != len(boxes) {
			t.Fatalf("expected %d entries, got %d", len(boxes), len(result))
		}

		for _, b := range boxes {
			if got, ok := result[b.Type]; !ok {
				t.Fatalf("expected type %q to be present", b.Type)
			} else if got != b {
				t.Fatalf("expected box %+v, got %+v", b, got)
			}
		}
	})

	t.Run("duplicate types last wins", func(t *testing.T) {
		t.Parallel()

		first := box.Box{ID: 1, Type: box.TypeCommon}
		second := box.Box{ID: 2, Type: box.TypeCommon}

		result := box.ToMapByType([]box.Box{first, second})

		if len(result) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(result))
		}

		if got := result[box.TypeCommon]; got != second {
			t.Fatalf("expected last box %+v to win, got %+v", second, got)
		}
	})
}
