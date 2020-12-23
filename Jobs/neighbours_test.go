package jobs

import (
	"testing"
	"time"
)

func TestNeighbours(t *testing.T) {
	m, res := Neighbours(1 * time.Second)
	t.Log(m, res)
}
