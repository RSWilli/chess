package chess

type attackRay struct {
	from BitBoard
	ray  BitBoard
}

// ray datastructures containing the respective ray until the end of the board for each square
var (
	northRays     = squareLookup[BitBoard]{}
	northEastRays = squareLookup[BitBoard]{}
	eastRays      = squareLookup[BitBoard]{}
	southEastRays = squareLookup[BitBoard]{}
	southRays     = squareLookup[BitBoard]{}
	northWestRays = squareLookup[BitBoard]{}
	westRays      = squareLookup[BitBoard]{}
	southWestRays = squareLookup[BitBoard]{}
)

func init() {
	type d struct {
		f     func(BitBoard) BitBoard
		store *squareLookup[BitBoard]
	}

	ds := []d{
		{BitBoard.Up, &northRays},
		{BitBoard.DiagUp, &northEastRays},
		{BitBoard.Right, &eastRays},
		{BitBoard.DiagDown, &southEastRays},
		{BitBoard.Down, &southRays},
		{BitBoard.AntiDiagUp, &northWestRays},
		{BitBoard.Left, &westRays},
		{BitBoard.AntiDiagDown, &southWestRays},
	}

	// generate rays in all queen move directions for each square
	for _, d := range ds {
		for i := range 64 {
			sq := BitBoard(1 << i)

			ray := BitBoard(sq)

			for sq := sq; sq != 0; sq = d.f(sq) {
				ray |= sq
			}

			d.store.set(sq, ray)
		}
	}
}
