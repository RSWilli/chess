package chess

import (
	"iter"
	"math/bits"
	"strings"
)

// bitBoard represents some kind of chess board state (either a move or a piece). Each bit
// in the uint64 represents a single square on the board. It is represented from a8-h8 through a1-h1
type bitBoard uint64

func (b bitBoard) String() string {
	var s strings.Builder
	var i bitBoard = 1 << 63

	for range 8 {
		for range 8 {
			if i&b != 0 {
				s.WriteString("x")
			} else {
				s.WriteString(".")
			}

			i = i >> 1
		}

		s.WriteString("\n")
	}

	return s.String()
}

func (b bitBoard) Each(f func(b bitBoard)) {
	for bb := range b.Ones() {
		f(bb)
	}
}

func (b bitBoard) Ones() iter.Seq[bitBoard] {
	return func(yield func(bitBoard) bool) {
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

func (b bitBoard) index() int {
	return 64 - bits.Len64(uint64(b))
}

func (b bitBoard) Has(other bitBoard) bool {
	return b&other != 0
}

func (b bitBoard) Set(other bitBoard) bitBoard {
	return b | other
}

func (b bitBoard) Unset(other bitBoard) bitBoard {
	return b &^ other
}

func (b bitBoard) Count() int {
	return bits.OnesCount64(uint64(b))
}

func (b bitBoard) Left() bitBoard {
	// prevent wrap around to h file, going left can never reach it
	return (b << 1) &^ bitBoard(hFile)
}

func (b bitBoard) Right() bitBoard {
	// prevent wrap around to a file, going right can never reach it
	return (b >> 1) &^ bitBoard(aFile)
}

func (b bitBoard) Up() bitBoard {
	return b << 8
}

func (b bitBoard) Down() bitBoard {
	return b >> 8
}

func (b bitBoard) DiagUp() bitBoard {
	return b.Up().Right()
}

func (b bitBoard) DiagDown() bitBoard {
	return b.Down().Right()
}

func (b bitBoard) AntiDiagUp() bitBoard {
	return b.Up().Left()
}

func (b bitBoard) AntiDiagDown() bitBoard {
	return b.Down().Left()
}

// flipVertical flips the [bitBoard] vertically
func (b bitBoard) flipVertical() bitBoard {
	return bitBoard(bits.ReverseBytes64(uint64(b)))
}

// flipHorizontal flips the [bitBoard] horizontally
// taken from https://www.chessprogramming.org/Flipping_Mirroring_and_Rotating and adapted to Go
func (b bitBoard) flipHorizontal() bitBoard {
	const k1 = 0x5555555555555555
	const k2 = 0x3333333333333333
	const k4 = 0x0f0f0f0f0f0f0f0f
	b = ((b >> 1) & k1) | ((b & k1) << 1)
	b = ((b >> 2) & k2) | ((b & k2) << 2)
	b = ((b >> 4) & k4) | ((b & k4) << 4)
	return b
}

// rotate180 rotates the [bitBoard], useful for switching sides.
func (b bitBoard) rotate180() bitBoard {
	return b.flipHorizontal().flipVertical()
}
