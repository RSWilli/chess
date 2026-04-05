package chess

import (
	"iter"
	"math/rand/v2"
)

// rookMoveMask returns a bitmask that has all bits set where an opposing color piece
// would restrict the possible moves of a rook
//
// this means that we skip the last square in each direction, because the rook can move there
// regardless whether there is an opposing piece or not.
func rookMoveMask(sq bitBoard) bitBoard {
	m := bitBoard(0)

	for c := sq.Right(); c != 0; c = c.Right() {
		m |= c &^ bitBoard(hFile)
	}
	for c := sq.Left(); c != 0; c = c.Left() {
		m |= c &^ bitBoard(aFile)
	}
	for c := sq.Up(); c != 0; c = c.Up() {
		m |= c &^ bitBoard(rank8)
	}
	for c := sq.Down(); c != 0; c = c.Down() {
		m |= c &^ bitBoard(rank1)
	}

	return m
}

// bishopMoveMask returns a bitmask that has all bits set where an opposing color piece
// would restrict the possible moves of a bishop
//
// this means that we skip the last square in each direction, because the bishop can move there
// regardless whether there is an opposing piece or not.
func bishopMoveMask(sq bitBoard) bitBoard {
	m := bitBoard(0)

	for c := sq.Up().Right(); c != 0; c = c.Up().Right() {
		m |= c &^ bitBoard(rank8|hFile)
	}
	for c := sq.Up().Left(); c != 0; c = c.Up().Left() {
		m |= c &^ bitBoard(rank8|aFile)
	}
	for c := sq.Down().Right(); c != 0; c = c.Down().Right() {
		m |= c &^ bitBoard(rank1|hFile)
	}
	for c := sq.Down().Left(); c != 0; c = c.Down().Left() {
		m |= c &^ bitBoard(rank1|aFile)
	}

	return m
}

// rookMoveTargetsSlow outputs the possible rook move targets given the opposite color occupies the given squares
//
// This means that the moves will end on top of the first blocking piece, signifying a capture move
func rookMoveTargetsSlow(sq, occupied bitBoard) bitBoard {
	m := bitBoard(0)

	for c := sq.Right(); c != 0; c = c.Right() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Left(); c != 0; c = c.Left() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Up(); c != 0; c = c.Up() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Down(); c != 0; c = c.Down() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}

	return m
}

// bishopMoveTargets outputs the possible rook move targets given the opposite color occupies the given squares
//
// This means that the moves will end on top of the first blocking piece, signifying a capture move
func bishopMoveTargetsSlow(sq, occupied bitBoard) bitBoard {
	m := bitBoard(0)

	for c := sq.Up().Right(); c != 0; c = c.Up().Right() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Up().Left(); c != 0; c = c.Up().Left() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Down().Right(); c != 0; c = c.Down().Right() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}
	for c := sq.Down().Left(); c != 0; c = c.Down().Left() {
		m |= c

		if c&occupied != 0 {
			break
		}
	}

	return m
}

func createOccupancyByIndex(mask bitBoard, index int) bitBoard {
	var res bitBoard

	for m := range mask.Ones() {
		if index&1 == 1 {
			res |= m
		}
		index = index >> 1
	}

	if index != 0 {
		panic("mask was not big enough")
	}

	return res
}

func possibleOccupationsCount(mask bitBoard) int {
	bits := mask.Count()

	return 1 << bits
}

func allPossibleOccupationsForSquare(mask bitBoard) []bitBoard {
	count := possibleOccupationsCount(mask)

	res := make([]bitBoard, count)

	for i := range count {
		res[i] = createOccupancyByIndex(mask, i)
	}

	return res
}

// rookBits is the amount of bits needed for the perfect hashing of the rook moves
// source https://www.chessprogramming.org/Looking_for_Magics
var rookBits = [64]bitBoard{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

// bishopBits is the amount of bits needed for the perfect hashing of the bishop moves
// source https://www.chessprogramming.org/Looking_for_Magics
var bishopBits = [64]bitBoard{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

// const maxShiftBits = BitBoard(12)

// findMagicNumber tries to find the best magic numbers for each of the squares
func findMagicNumber(occupations []bitBoard, bits bitBoard) bitBoard {
nextmagic:
	for {
		magic := newRandom()

		// track the taken indices, we don't want collisions
		taken := make([]bool, 1<<bits)

		for _, occupancy := range occupations {
			index := hash(occupancy, magic, bits)

			if taken[index] {
				continue nextmagic // hash collision
			}

			taken[index] = true
		}

		return magic
	}
}

func hash(b, magic, bits bitBoard) int {
	return int((b * magic) >> (64 - bits))
}

// newRandom returns a random uint64 for the multiplication part of the hash
func newRandom() bitBoard {
	// https://www.chessprogramming.org/Looking_for_Magics#Feeding_in_Randoms
	// a good approach is to have random numbers with very few 1 bits

	a, b, c := rand.Uint64(), rand.Uint64(), rand.Uint64()

	return bitBoard(a & b & c)
}

type entry struct {
	Key   int
	Value uint64
}

type magic struct {
	Magic  uint64
	Bits   uint64
	Values []entry
}

type pieceData struct {
	MoveMasks []uint64
	Magics    []magic
}

type data struct {
	Rook   pieceData
	Bishop pieceData
}

func (d data) Parts() iter.Seq2[string, pieceData] {
	return func(yield func(string, pieceData) bool) {
		if !yield("rook", d.Rook) {
			return
		}
		if !yield("bishop", d.Bishop) {
			return
		}
	}
}

type maskFunc = func(sq bitBoard) bitBoard
type moveTargetsFunc = func(sq, occupied bitBoard) bitBoard

// initMoveMasks returns the move mask [squareLookup] for the mask function
func initMoveMasks(maskFunc maskFunc) (res squareLookup[bitBoard]) {
	for i := range 64 {
		sq := bitBoard(NewSquareFromIndex(i))

		res.set(sq, maskFunc(sq))
	}

	return
}

// initMoveLookupTable returns a [squareLookup] that allows for fast retrieval of
// move targets for a piece with the given move function and mask
func initMoveLookupTable(maskFunc maskFunc, moveTargets moveTargetsFunc, hashbits [64]bitBoard) (res squareLookup[hashTable]) {
	for i := range 64 {
		sq := bitBoard(NewSquareFromIndex(i))

		mask := maskFunc(sq)

		bits := hashbits[i]

		allOccupations := allPossibleOccupationsForSquare(mask)

		magic := findMagicNumber(allOccupations, bits)

		table := hashTable{
			magic: magic,
			bits:  bits,
			data:  make([]bitBoard, 1<<bits),
		}

		for _, occ := range allOccupations {
			key := hash(occ, magic, bits)
			table.data[key] = moveTargets(sq, occ)
		}

		res.set(sq, table)
	}

	return
}
