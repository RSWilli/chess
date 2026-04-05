package chess

var rookMoveMasks = initMoveMasks(rookMoveMask)
var rookMoveLookup = initMoveLookupTable(rookMoveMask, rookMoveTargetsSlow, rookBits)

var bishopMoveMasks = initMoveMasks(bishopMoveMask)
var bishopMoveLookup = initMoveLookupTable(bishopMoveMask, bishopMoveTargetsSlow, bishopBits)

func kingMoves(king bitBoard) bitBoard {
	return king.Up() |
		king.Down() |
		king.Up().Left() |
		king.Down().Left() |
		king.Up().Right() |
		king.Down().Right() |
		king.Left() |
		king.Right()
}

func knightMoves(knight bitBoard, same, _ bitBoard) bitBoard {
	targets := knight.Up().Up().Left() |
		knight.Up().Up().Right() |

		knight.Down().Down().Left() |
		knight.Down().Down().Right() |

		knight.Up().Left().Left() |
		knight.Up().Right().Right() |

		knight.Down().Left().Left() |
		knight.Down().Right().Right()

	return targets &^ same
}

func rookMoves(sq, same, opposing bitBoard) bitBoard {
	all := same | opposing

	mask := rookMoveMasks.get(sq)

	relevant := mask & all

	// moves contains friendly fire targets
	moves := rookMoveLookup.get(sq).get(relevant)

	return moves &^ same
}

func bishopMoves(sq bitBoard, same, opposing bitBoard) bitBoard {
	all := same | opposing

	mask := bishopMoveMasks.get(sq)

	relevant := mask & all

	// moves contains friendly fire targets
	moves := bishopMoveLookup.get(sq).get(relevant)

	return moves &^ same
}

func queenMoves(sq bitBoard, same, opposing bitBoard) bitBoard {
	return rookMoves(sq, same, opposing) | bishopMoves(sq, same, opposing)
}

// opponentPawnAttacks returns a [bitBoard] of all squares an opponent pawn attacks, from the current player's
// perspective
func opponentPawnAttacks(pawn bitBoard) bitBoard {
	return pawn.Down().Left() | pawn.Down().Right()
}
