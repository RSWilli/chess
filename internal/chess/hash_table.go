package chess

// hashTable is a datastructure utilizing perfect hashing for handling the magic bitboards
// of the sliding piece moves
//
// it works by hashing a given occupancy bitboard to return the correct sliding moves.
//
// see [hashTable.get] for the hashing formula
type hashTable struct {
	magic bitBoard
	bits  bitBoard

	// data must be initialized with the values at their correct location
	// as defined by rotate and shift
	data []bitBoard
}

// get returns the possible target squares when given the occupancy bitBoard
func (h hashTable) get(key bitBoard) bitBoard {
	i := (key * h.magic) >> (64 - h.bits)

	return h.data[i]
}
