package chess

import (
	"testing"
)

func TestSlides(t *testing.T) {

	bb := bitBoard(NewSquare(5, 7))

	t.Log(bb)
	t.Log(bb.Right())
	t.Log(bb.Up())
	t.Log(bb.Up().Right())
}

func TestRotations(t *testing.T) {
	cases := []bitBoard{
		0xffff,
		0b1000000001000000010000001000001000010001001010,
	}

	for _, b := range cases {
		t.Logf("\n%s\n", b.String())
		t.Logf("\n%s\n", b.rotate180().String())
	}
}
