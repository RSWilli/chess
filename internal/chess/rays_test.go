package chess

import (
	"testing"
)

func TestRays(t *testing.T) {
	for _, ray := range northRays.store {
		t.Log(ray)
	}
}
