package chess

// ray is a name for a [BitBoard] that connects from and to, but without containing to from and to
type ray struct {
	from BitBoard
	to   BitBoard
}

type attackRay struct {
	from BitBoard
	ray  BitBoard
}

// initialize rays to hold the correct amount
var rays = make(map[ray]BitBoard, 1036)

func init() {
	// generate rays in all queen move directions for each square
	for i := range 64 {
		from := BitBoard(1 << i)

		stepFuncs := []func(BitBoard) BitBoard{
			BitBoard.Right,
			BitBoard.Left,
			BitBoard.Up,
			BitBoard.Down,

			BitBoard.DiagUp,
			BitBoard.DiagDown,
			BitBoard.AntiDiagUp,
			BitBoard.AntiDiagDown,
		}

		for _, step := range stepFuncs {
			r := BitBoard(0)
			for to := step(from); to != 0; to = step(to) {
				if r != 0 {
					rays[ray{from, to}] = r
				}
				r |= to
			}
		}
	}
}
