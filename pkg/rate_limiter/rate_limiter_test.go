package rate_limiter

import (
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func BenchmarkXxx(b *testing.B) {
	now := time.Now()
	var l *limiterEntry
	for b.Loop() {
		l = &limiterEntry{
			l:        rate.NewLimiter(0, 0),
			lastSeen: now,
		}
	}
	_ = l
}
