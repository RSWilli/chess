package chess

func kingMoves(king BitBoard) BitBoard {
	return king.Up() |
		king.Down() |
		king.Up().Left() |
		king.Down().Left() |
		king.Up().Right() |
		king.Down().Right() |
		king.Left() |
		king.Right()
}

func knightMoves(knight BitBoard, same, _ BitBoard) BitBoard {
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

func rookMoves(sq, same, opposing BitBoard) BitBoard {
	all := same | opposing

	mask := rookMoveMasks.get(sq)

	relevant := mask & all

	// moves contains friendly fire targets
	moves := rookMoveTargets.get(sq).get(relevant)

	return moves &^ same
}

func bishopMoves(sq BitBoard, same, opposing BitBoard) BitBoard {
	all := same | opposing

	mask := bishopMoveMasks.get(sq)

	relevant := mask & all

	// moves contains friendly fire targets
	moves := bishopMoveTargets.get(sq).get(relevant)

	return moves &^ same
}

func queenMoves(sq BitBoard, same, opposing BitBoard) BitBoard {
	return rookMoves(sq, same, opposing) | bishopMoves(sq, same, opposing)
}

func whitePawnAttacks(pawn BitBoard) BitBoard {
	return pawn.Up().Left() | pawn.Up().Right()
}

func blackPawnAttacks(pawn BitBoard) BitBoard {
	return pawn.Down().Left() | pawn.Down().Right()
}
