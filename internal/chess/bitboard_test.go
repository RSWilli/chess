package chess

import (
	"fmt"
	"testing"
)

func TestSameRank(t *testing.T) {
	type test struct {
		a Square
		b Square

		expect bool
	}
	tests := []test{
		{NewSquare(7, 5), NewSquare(7, 2), true},
		{NewSquare(3, 5), NewSquare(3, 2), true},
		{NewSquare(7, 7), NewSquare(7, 7), true},
		{NewSquare(2, 5), NewSquare(3, 2), false},
		{NewSquare(7, 5), NewSquare(1, 2), false},
		{NewSquare(7, 5), NewSquare(0, 6), false},
		{NewSquare(5, 0), NewSquare(4, 7), false},
		{0, NewSquare(3, 2), false},
	}

	for _, tc := range tests {
		fmt.Println(tc.a, tc.b)
		out := BitBoard(tc.a).isSameRank(BitBoard(tc.b))

		if out != tc.expect {
			t.Fatal("failed")
		}
	}
}

func TestSlides(t *testing.T) {

	bb := BitBoard(NewSquare(5, 7))

	fmt.Println(bb)
	fmt.Println(bb.Right())
	fmt.Println(bb.Up())
	fmt.Println(bb.Up().Right())
}
