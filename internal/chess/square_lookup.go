package chess

import "math/bits"

// squareLookup is a datastructure that allows to lookup a value for a bitBoard or a square.
//
// it is more efficient than using a map (TODO: bench)
type squareLookup[V any] struct {
	store [64]V
}

func (s *squareLookup[V]) get(b BitBoard) V {
	return s.store[bits.Len64(uint64(b))-1]
}

func (s *squareLookup[V]) set(b BitBoard, v V) {
	s.store[bits.Len64(uint64(b))-1] = v
}

// func (s *squareLookup[V]) getSquare(b Square) V {
// 	return s.store[bits.Len64(uint64(b))]
// }

// func (s *squareLookup[V]) setSquare(b Square, v V) {
// 	s.store[bits.Len64(uint64(b))] = v
// }
