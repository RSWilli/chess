package chess

import (
	"iter"
	"math/bits"
)

// BitBoard represents some kind of chess board state (either a move or a piece). Each bit
// in the uint64 represents a single square on the board. It is represented from a8-h8 through a1-h1
type BitBoard uint64

func (b BitBoard) String() string {
	s := ""
	var i BitBoard = 1

	for range 8 {
		for range 8 {
			if i&b != 0 {
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

func (b BitBoard) Each(f func(b BitBoard)) {
	for i := range BitBoard(64) {
		if b&(1<<i) == 0 {
			continue
		}

		f(1 << i)
	}
}

func (b BitBoard) Ones() iter.Seq[BitBoard] {
	return func(yield func(BitBoard) bool) {
		for i := range BitBoard(64) {
			if b&(1<<i) == 0 {
				continue
			}

			if !yield(1 << i) {
				break
			}
		}
	}
}

func (b BitBoard) Has(sq Square) bool {
	return b&BitBoard(sq) != 0
}

func (b BitBoard) Set(sq Square) BitBoard {
	return b | BitBoard(sq)
}

func (b BitBoard) Unset(sq Square) BitBoard {
	return b &^ BitBoard(sq)
}

func (b BitBoard) Count() int {
	return bits.OnesCount64(uint64(b))
}

func (b BitBoard) Left() BitBoard {
	s := b >> 1

	if s == 0 || !b.isSameRank(s) {
		return 0
	}

	return s
}

func (b BitBoard) Right() BitBoard {
	s := b << 1

	if !b.isSameRank(s) {
		return 0
	}

	return s
}

func (b BitBoard) Up() BitBoard {
	return b >> 8
}

func (b BitBoard) Down() BitBoard {
	return b << 8
}

func (b BitBoard) DiagUp() BitBoard {
	return b.Up().Right()
}

func (b BitBoard) DiagDown() BitBoard {
	return b.Down().Right()
}

func (b BitBoard) AntiDiagUp() BitBoard {
	return b.Up().Left()
}

func (b BitBoard) AntiDiagDown() BitBoard {
	return b.Down().Left()
}

// isSameRank efficiently checks whether the two bitboard squares are on the same rank
//
// this is useful to check whether we accidentally wrapped around from a<->h file
func (b BitBoard) isSameRank(other BitBoard) bool {
	if b == 0 || other == 0 {
		return false // one or both is outside of the board
	}

	// continuously move both to the first rank.
	for range 8 {
		b = b >> 8
		other = other >> 8

		if (b > 0) != (other > 0) {
			return false // one of the two left the board earlier
		}
	}

	return true
}
