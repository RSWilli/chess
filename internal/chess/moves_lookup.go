package chess

var rookMovesLookupTable = squareLookup[BitBoard]{}
var bishopMovesLookupTable = squareLookup[BitBoard]{}

func rookMoves(sq BitBoard) BitBoard {
	return rookMovesLookupTable.get(sq)
}

func bishopMoves(sq BitBoard) BitBoard {
	return bishopMovesLookupTable.get(sq)
}

func queenMoves(sq BitBoard) BitBoard {
	return rookMoves(sq) | bishopMoves(sq)
}
