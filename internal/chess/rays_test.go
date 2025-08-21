package chess

import (
	"fmt"
	"testing"
)

func printRay(t *testing.T, r ray, b BitBoard) string {
	t.Helper()

	s := ""
	var i BitBoard = 1

	for range 8 {
		for range 8 {
			if i&r.to != 0 {
				s += "e"
			} else if i&r.from != 0 {
				s += "s"
			} else if i&b != 0 {
				s += "x"
			} else {
				s += "."
			}

			i = i << 1
		}

		s += "\n"
	}

	return s
}

func TestRays(t *testing.T) {
	for r, ray := range rays {
		// ray must contain from but not to:
		if r.from&ray != 0 {
			t.Log("ray does contain from:")
			t.Log(r.from.String())
			t.Log(ray.String())
			t.Fatal(printRay(t, r, ray))
		}

		if r.to&ray != 0 {
			t.Log("ray does contain to:")
			t.Log(r.to.String())
			t.Log(ray.String())
			t.Fatal(printRay(t, r, ray))
		}
		fmt.Println(printRay(t, r, ray))
	}

	fmt.Println(len(rays))
}
