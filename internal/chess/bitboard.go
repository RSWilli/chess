package chess

import (
	"iter"
	"math/bits"
)

const fileHBitBoard bitBoard = 0x0101010101010101

const (
	fileGBitBoard = fileHBitBoard << (iota + 1)
	fileFBitBoard
	fileEBitBoard
	fileDBitBoard
	fileCBitBoard
	fileBBitBoard
	fileABitBoard
)

const rank8BitBoard bitBoard = 0x00000000000000ff

const (
	rank7BitBoard = rank8BitBoard << ((iota + 1) * 8)
	rank6BitBoard
	rank5BitBoard
	rank4BitBoard
	rank3BitBoard
	rank2BitBoard
	rank1BitBoard
)

const pos8DiagBitBoard bitBoard = 0x0102040810204080

const (
	pos7DiagBitBoard bitBoard = (pos8DiagBitBoard >> 1) &^ fileABitBoard
	pos6DiagBitBoard          = (pos7DiagBitBoard >> 1) &^ fileABitBoard
	pos5DiagBitBoard          = (pos6DiagBitBoard >> 1) &^ fileABitBoard
	pos4DiagBitBoard          = (pos5DiagBitBoard >> 1) &^ fileABitBoard
	pos3DiagBitBoard          = (pos4DiagBitBoard >> 1) &^ fileABitBoard
	pos2DiagBitBoard          = (pos3DiagBitBoard >> 1) &^ fileABitBoard
	pos1DiagBitBoard          = (pos2DiagBitBoard >> 1) &^ fileABitBoard

	pos9DiagBitBoard  bitBoard = (pos8DiagBitBoard << 1) &^ fileHBitBoard
	pos10DiagBitBoard          = (pos9DiagBitBoard << 1) &^ fileHBitBoard
	pos11DiagBitBoard          = (pos10DiagBitBoard << 1) &^ fileHBitBoard
	pos12DiagBitBoard          = (pos11DiagBitBoard << 1) &^ fileHBitBoard
	pos13DiagBitBoard          = (pos12DiagBitBoard << 1) &^ fileHBitBoard
	pos14DiagBitBoard          = (pos13DiagBitBoard << 1) &^ fileHBitBoard
	pos15DiagBitBoard          = (pos14DiagBitBoard << 1) &^ fileHBitBoard
)

const neg8DiagBitBoard bitBoard = 0x8040201008040201

const (
	neg7DiagBitBoard bitBoard = (neg8DiagBitBoard >> 1) &^ fileABitBoard
	neg6DiagBitBoard          = (neg7DiagBitBoard >> 1) &^ fileABitBoard
	neg5DiagBitBoard          = (neg6DiagBitBoard >> 1) &^ fileABitBoard
	neg4DiagBitBoard          = (neg5DiagBitBoard >> 1) &^ fileABitBoard
	neg3DiagBitBoard          = (neg4DiagBitBoard >> 1) &^ fileABitBoard
	neg2DiagBitBoard          = (neg3DiagBitBoard >> 1) &^ fileABitBoard
	neg1DiagBitBoard          = (neg2DiagBitBoard >> 1) &^ fileABitBoard

	neg9DiagBitBoard  bitBoard = 0x0080402010080402 // left shift would overflow
	neg10DiagBitBoard          = (neg9DiagBitBoard << 1) &^ fileHBitBoard
	neg11DiagBitBoard          = (neg10DiagBitBoard << 1) &^ fileHBitBoard
	neg12DiagBitBoard          = (neg11DiagBitBoard << 1) &^ fileHBitBoard
	neg13DiagBitBoard          = (neg12DiagBitBoard << 1) &^ fileHBitBoard
	neg14DiagBitBoard          = (neg13DiagBitBoard << 1) &^ fileHBitBoard
	neg15DiagBitBoard          = (neg14DiagBitBoard << 1) &^ fileHBitBoard
)

// bitBoard represents some kind of chess board state (either a move or a piece). Each bit
// in the uint64 represents a single square on the board. It is represented from a8-h8 through a1-h1
type bitBoard uint64

func (b bitBoard) String() string {
	s := ""
	var i bitBoard = 1

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

func (b bitBoard) each(f func(b bitBoard)) {
	for i := range bitBoard(64) {
		if b&(1<<i) == 0 {
			continue
		}

		f(1 << i)
	}
}

func (b bitBoard) ones() iter.Seq[bitBoard] {
	return func(yield func(bitBoard) bool) {
		for i := range bitBoard(64) {
			if b&(1<<i) == 0 {
				continue
			}

			if !yield(1 << i) {
				break
			}
		}
	}
}

func (b bitBoard) has(sq Square) bool {
	return b&bitBoard(sq) != 0
}

func (b bitBoard) set(sq Square) bitBoard {
	return b | bitBoard(sq)
}

func (b bitBoard) unset(sq Square) bitBoard {
	return b &^ bitBoard(sq)
}

func (b bitBoard) count() int {
	return bits.OnesCount64(uint64(b))
}

func (b bitBoard) left() bitBoard {
	s := b >> 1

	if s == 0 || !b.isSameRank(s) {
		return 0
	}

	return s
}

func (b bitBoard) right() bitBoard {
	s := b << 1

	if !b.isSameRank(s) {
		return 0
	}

	return s
}

func (b bitBoard) up() bitBoard {
	return b >> 8
}

func (b bitBoard) down() bitBoard {
	return b << 8
}

// isSameRank efficiently checks whether the two bitboard squares are on the same rank
//
// this is useful to check whether we accidentally wrapped around from a<->h file
func (b bitBoard) isSameRank(other bitBoard) bool {
	if b == 0 || other == 0 {
		return false // on or both is outside of the board
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
