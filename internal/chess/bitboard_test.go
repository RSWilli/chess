package chess

import (
	"fmt"
	"testing"
)

func TestSlides(t *testing.T) {

	bb := BitBoard(NewSquare(5, 7))

	fmt.Println(bb)
	fmt.Println(bb.Right())
	fmt.Println(bb.Up())
	fmt.Println(bb.Up().Right())
}
