package chess

import "math/bits"

const fileHBitBoard bitBoard = 0x0101010101010101

const (
	fileGBitBoard = fileHBitBoard << iota
	fileFBitBoard
	fileEBitBoard
	fileDBitBoard
	fileCBitBoard
	fileBBitBoard
	fileABitBoard
)

const rank1BitBoard bitBoard = 0x00000000000000ff

const (
	rank2BitBoard = rank1BitBoard << iota * 8
	rank3BitBoard
	rank4BitBoard
	rank5BitBoard
	rank6BitBoard
	rank7BitBoard
	rank8BitBoard
)

// bitBoard represents some kind of chess board state (either a move or a piece). Each bit
// in the uint64 represents a single square on the board. It is represented from a8-h8 through a1-h1
type bitBoard uint64

func (b bitBoard) String() string {
	s := ""
	var i bitBoard = 1 << 63

	for range 8 {
		for range 8 {
			if i&b != 0 {
				s += "x"
			} else {
				s += "."
			}

			i = i >> 1
		}

		s += "\n"
	}

	return s
}

func (b bitBoard) hasSquare(sq Square) bool {
	return b.has(sq.Index())
}

func (b bitBoard) has(i int) bool {
	return b&(1<<i) != 0
}

func (b bitBoard) setSquare(sq Square) bitBoard {
	return b.set(sq.Index())
}

func (b bitBoard) set(i int) bitBoard {
	return b | 1<<i
}

func (b bitBoard) unsetSquare(sq Square) bitBoard {
	return b.unset(sq.Index())
}

func (b bitBoard) unset(i int) bitBoard {
	return b &^ (1 << i)
}

func (b bitBoard) count() int {
	return bits.OnesCount64(uint64(b))
}

func (b bitBoard) left() bitBoard {
	return b << 1
}

func (b bitBoard) right() bitBoard {
	return b >> 1
}

func (b bitBoard) up() bitBoard {
	return b << 8
}

func (b bitBoard) down() bitBoard {
	return b >> 8
}
