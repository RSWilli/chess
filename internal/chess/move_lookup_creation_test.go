package chess

import "testing"

func TestMoveLookupGen(t *testing.T) {
	magic := initMoveLookupTable(rookMoveMask, rookMoveTargetsSlow, rookBits)

	_ = magic
}
