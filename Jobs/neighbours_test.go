package jobs

import (
	"testing"
	"time"
)

func TestNeighbours(t *testing.T) {
	m := neighbours(1 * time.Second)
	t.Log(m)
}
