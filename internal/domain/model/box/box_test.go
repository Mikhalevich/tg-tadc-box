package box_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/box"
	"github.com/Mikhalevich/tg-tadc-box/internal/domain/model/msginfo"
)

func TestToInProgressBoxMapByType(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 17, 12, 0, 0, 0, time.UTC)

	newBox := func(id int, bType box.Type, createdAt, availableAt time.Time) box.Box {
		return box.Box{
			ID:                  box.IDFromInt(id),
			ChatID:              msginfo.ChatIDFromInt64(int64(id)),
			Status:              box.StatusInProgress,
			Type:                bType,
			CreatedAt:           createdAt,
			AvailableAt:         availableAt,
			ReadyNotificationAt: createdAt.Add(time.Hour),
			CompletedAt:         availableAt.Add(time.Hour),
			Meta: box.Meta{
				BonusBox: box.BonusBox{
					IsValid:  true,
					Type:     bType,
					Attempts: 5,
				},
			},
		}
	}

	tests := map[string]struct {
		boxes []box.Box
		now   time.Time
		want  map[box.Type]box.InProgressBox
	}{
		"empty boxes returns nil": {
			boxes: []box.Box{},
			now:   now,
			want:  nil,
		},
		"nil boxes returns nil": {
			boxes: nil,
			now:   now,
			want:  nil,
		},
		"single box available in the future": {
			boxes: []box.Box{
				newBox(1, box.TypeCommon, now.Add(-2*time.Hour), now.Add(3*time.Hour)),
			},
			now: now,
			want: map[box.Type]box.InProgressBox{
				box.TypeCommon: {
					Box:            newBox(1, box.TypeCommon, now.Add(-2*time.Hour), now.Add(3*time.Hour)),
					AvailableAfter: 3 * time.Hour,
				},
			},
		},
		"multiple boxes with different types": {
			boxes: []box.Box{
				newBox(1, box.TypeCommon, now.Add(-2*time.Hour), now.Add(3*time.Hour)),
				newBox(2, box.TypeRare, now.Add(-1*time.Hour), now.Add(5*time.Hour)),
				newBox(3, box.TypeEpic, now.Add(-3*time.Hour), now.Add(1*time.Hour)),
				newBox(4, box.TypeLegendary, now.Add(-4*time.Hour), now.Add(10*time.Hour)),
			},
			now: now,
			want: map[box.Type]box.InProgressBox{
				box.TypeCommon: {
					Box:            newBox(1, box.TypeCommon, now.Add(-2*time.Hour), now.Add(3*time.Hour)),
					AvailableAfter: 3 * time.Hour,
				},
				box.TypeRare: {
					Box:            newBox(2, box.TypeRare, now.Add(-1*time.Hour), now.Add(5*time.Hour)),
					AvailableAfter: 5 * time.Hour,
				},
				box.TypeEpic: {
					Box:            newBox(3, box.TypeEpic, now.Add(-3*time.Hour), now.Add(1*time.Hour)),
					AvailableAfter: 1 * time.Hour,
				},
				box.TypeLegendary: {
					Box:            newBox(4, box.TypeLegendary, now.Add(-4*time.Hour), now.Add(10*time.Hour)),
					AvailableAfter: 10 * time.Hour,
				},
			},
		},
		"multiple boxes with same type keeps the last one": {
			boxes: []box.Box{
				newBox(1, box.TypeCommon, now.Add(-2*time.Hour), now.Add(3*time.Hour)),
				newBox(2, box.TypeCommon, now.Add(-1*time.Hour), now.Add(7*time.Hour)),
			},
			now: now,
			want: map[box.Type]box.InProgressBox{
				box.TypeCommon: {
					Box:            newBox(2, box.TypeCommon, now.Add(-1*time.Hour), now.Add(7*time.Hour)),
					AvailableAfter: 7 * time.Hour,
				},
			},
		},
		"box already available gives negative AvailableAfter": {
			boxes: []box.Box{
				newBox(1, box.TypeCommon, now.Add(-5*time.Hour), now.Add(-2*time.Hour)),
			},
			now: now,
			want: map[box.Type]box.InProgressBox{
				box.TypeCommon: {
					Box:            newBox(1, box.TypeCommon, now.Add(-5*time.Hour), now.Add(-2*time.Hour)),
					AvailableAfter: -2 * time.Hour,
				},
			},
		},
		"box available exactly now gives zero AvailableAfter": {
			boxes: []box.Box{
				newBox(1, box.TypeCommon, now.Add(-5*time.Hour), now),
			},
			now: now,
			want: map[box.Type]box.InProgressBox{
				box.TypeCommon: {
					Box:            newBox(1, box.TypeCommon, now.Add(-5*time.Hour), now),
					AvailableAfter: 0,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := box.ToInProgressBoxMapByType(tt.boxes, tt.now)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInProgressBoxMapByType() = %v, want %v", got, tt.want)
			}
		})
	}
}
