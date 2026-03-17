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
	for bb := range b.Ones() {
		f(bb)
	}
}

func (b BitBoard) Ones() iter.Seq[BitBoard] {
	return func(yield func(BitBoard) bool) {
		tmp := b
		for tmp != 0 {
			// find the least significant bit
			lsb := tmp & -tmp

			// remove the lsb from tmp for next iteration
			tmp &^= tmp & -tmp

			if !yield(lsb) {
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

const (
	aFile BitBoard = 0x0101010101010101
	hFile          = aFile << 7
)

func (b BitBoard) Left() BitBoard {
	// prevent wrap around to h file, going left can never reach it
	return (b >> 1) &^ hFile
}

func (b BitBoard) Right() BitBoard {
	// prevent wrap around to a file, going right can never reach it
	return (b << 1) &^ aFile
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
