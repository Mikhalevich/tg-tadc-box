//nolint:testpackage
package rewardgenerator

import (
	"testing"
)

func BenchmarkPercent(b *testing.B) {
	for b.Loop() {
		percent()
	}
}
