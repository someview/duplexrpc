package breaker

import (
	"testing"
	"time"
)

func TestBreaker(t *testing.T) {
	b := NewBreaker(3, 4*time.Second)
	for b.Allow() {

		b.Fail()
	}
}
