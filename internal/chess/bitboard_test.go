package chess

import (
	"testing"
)

func TestSlides(t *testing.T) {

	bb := BitBoard(NewSquare(5, 7))

	t.Log(bb)
	t.Log(bb.Right())
	t.Log(bb.Up())
	t.Log(bb.Up().Right())
}
