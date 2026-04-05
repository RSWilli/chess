package chess

// squareLookup is a datastructure that allows to lookup a value for a bitBoard or a square.
//
// it is more efficient than using a map (see benchmark)
type squareLookup[V any] struct {
	store [64]V
}

func (s *squareLookup[V]) get(b bitBoard) V {
	return s.store[b.index()]
}

func (s *squareLookup[V]) set(b bitBoard, value V) {
	s.store[b.index()] = value
}
