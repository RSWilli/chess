package chess

import (
	"fmt"
	"testing"
)

func TestDiags(t *testing.T) {
	fmt.Println("Positive")
	fmt.Println(pos1DiagBitBoard)
	fmt.Println(pos2DiagBitBoard)
	fmt.Println(pos3DiagBitBoard)
	fmt.Println(pos4DiagBitBoard)
	fmt.Println(pos5DiagBitBoard)
	fmt.Println(pos6DiagBitBoard)
	fmt.Println(pos7DiagBitBoard)
	fmt.Println(pos8DiagBitBoard)
	fmt.Println(pos9DiagBitBoard)
	fmt.Println(pos10DiagBitBoard)
	fmt.Println(pos11DiagBitBoard)
	fmt.Println(pos12DiagBitBoard)
	fmt.Println(pos13DiagBitBoard)
	fmt.Println(pos14DiagBitBoard)
	fmt.Println(pos15DiagBitBoard)

	fmt.Println("Negative")
	fmt.Println(neg1DiagBitBoard)
	fmt.Println(neg2DiagBitBoard)
	fmt.Println(neg3DiagBitBoard)
	fmt.Println(neg4DiagBitBoard)
	fmt.Println(neg5DiagBitBoard)
	fmt.Println(neg6DiagBitBoard)
	fmt.Println(neg7DiagBitBoard)
	fmt.Println(neg8DiagBitBoard)
	fmt.Println(neg9DiagBitBoard)
	fmt.Println(neg10DiagBitBoard)
	fmt.Println(neg11DiagBitBoard)
	fmt.Println(neg12DiagBitBoard)
	fmt.Println(neg13DiagBitBoard)
	fmt.Println(neg14DiagBitBoard)
	fmt.Println(neg15DiagBitBoard)
}

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
		out := bitBoard(tc.a).isSameRank(bitBoard(tc.b))

		if out != tc.expect {
			t.Fatal("failed")
		}
	}
}

func TestSlides(t *testing.T) {

	bb := bitBoard(NewSquare(5, 7))

	fmt.Println(bb)
	fmt.Println(bb.right())
	fmt.Println(bb.up())
	fmt.Println(bb.up().right())
}
