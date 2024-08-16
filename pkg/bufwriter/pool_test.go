package bufwriter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcLevel(t *testing.T) {
	levels := [8]int{0, 10, 16, 31, 64, 128, 256, 257}
	for i := 0; i < levels[i]; i++ {
		assert.Equal(t, levels[i], calcLevel(calcSize(levels[i])))
	}
}
