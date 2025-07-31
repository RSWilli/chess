package chess

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
