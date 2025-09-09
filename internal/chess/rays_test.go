package chess

import (
	"fmt"
	"testing"
)

func TestRays(t *testing.T) {
	for _, ray := range northRays.store {
		fmt.Println(ray)
	}
}
