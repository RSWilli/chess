package chess

// ray is a name for a [BitBoard] that connects from and to
type ray struct {
	from BitBoard
	to   BitBoard
}

var rays = map[ray]BitBoard{}

func init() {
	// generate rays in all queen move directions for each square
	for i := range 64 {
		var r BitBoard
		from := BitBoard(1 << i)

		r = from
		for to := from.Right(); to != 0; to = to.Right() {
			r |= to
			rays[ray{from, to}] = r
		}
		r = from
		for to := from.Left(); to != 0; to = to.Left() {
			r |= to
			rays[ray{from, to}] = r
		}
		r = from
		for to := from.Up(); to != 0; to = to.Up() {
			r |= to
			rays[ray{from, to}] = r
		}
		r = from
		for to := from.Down(); to != 0; to = to.Down() {
			r |= to
			rays[ray{from, to}] = r
		}

		r = from
		for to := from.Up().Right(); to != 0; to = to.Up().Right() {
			r |= to
			rays[ray{from, to}] = r
		}
		r = from
		for to := from.Up().Left(); to != 0; to = to.Up().Left() {
			r |= to
			rays[ray{from, to}] = r
		}
		r = from
		for to := from.Down().Right(); to != 0; to = to.Down().Right() {
			r |= to
			rays[ray{from, to}] = r
		}
		r = from
		for to := from.Down().Left(); to != 0; to = to.Down().Left() {
			r |= to
			rays[ray{from, to}] = r
		}
	}
}
