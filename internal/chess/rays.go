package chess

type attackRay struct {
	from bitBoard
	ray  bitBoard
}

// ray datastructures containing the respective ray until the end of the board for each square
var (
	northRays     = squareLookup[bitBoard]{}
	northEastRays = squareLookup[bitBoard]{}
	eastRays      = squareLookup[bitBoard]{}
	southEastRays = squareLookup[bitBoard]{}
	southRays     = squareLookup[bitBoard]{}
	northWestRays = squareLookup[bitBoard]{}
	westRays      = squareLookup[bitBoard]{}
	southWestRays = squareLookup[bitBoard]{}
)

func init() {
	type d struct {
		f     func(bitBoard) bitBoard
		store *squareLookup[bitBoard]
	}

	ds := []d{
		{bitBoard.Up, &northRays},
		{bitBoard.DiagUp, &northEastRays},
		{bitBoard.Right, &eastRays},
		{bitBoard.DiagDown, &southEastRays},
		{bitBoard.Down, &southRays},
		{bitBoard.AntiDiagUp, &northWestRays},
		{bitBoard.Left, &westRays},
		{bitBoard.AntiDiagDown, &southWestRays},
	}

	// generate rays in all queen move directions for each square
	for _, d := range ds {
		for i := range 64 {
			sq := bitBoard(1 << i)

			ray := bitBoard(sq)

			for sq := sq; sq != 0; sq = d.f(sq) {
				ray |= sq
			}

			d.store.set(sq, ray)
		}
	}
}
