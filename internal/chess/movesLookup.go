package chess

// TODO: use something more sophisticated for faster lookup
var rookMovesLookupTable = map[bitBoard]bitBoard{}
var bishopMovesLookupTable = map[bitBoard]bitBoard{}

func init() {
	ranks := []bitBoard{
		rank1BitBoard,
		rank2BitBoard,
		rank3BitBoard,
		rank4BitBoard,
		rank5BitBoard,
		rank6BitBoard,
		rank7BitBoard,
		rank8BitBoard,
	}

	files := []bitBoard{
		fileABitBoard,
		fileBBitBoard,
		fileCBitBoard,
		fileDBitBoard,
		fileEBitBoard,
		fileFBitBoard,
		fileGBitBoard,
		fileHBitBoard,
	}

	for _, rank := range ranks {
		for _, file := range files {
			sq := rank & file

			moves := (rank | file) &^ sq

			rookMovesLookupTable[sq] = moves
		}
	}

	if len(rookMovesLookupTable) != 64 {
		panic("missing moves")
	}
}

func init() {
	pos := []bitBoard{
		pos1DiagBitBoard,
		pos2DiagBitBoard,
		pos3DiagBitBoard,
		pos4DiagBitBoard,
		pos5DiagBitBoard,
		pos6DiagBitBoard,
		pos7DiagBitBoard,
		pos8DiagBitBoard,
		pos9DiagBitBoard,
		pos10DiagBitBoard,
		pos11DiagBitBoard,
		pos12DiagBitBoard,
		pos13DiagBitBoard,
		pos14DiagBitBoard,
		pos15DiagBitBoard,
	}

	neg := []bitBoard{
		neg1DiagBitBoard,
		neg2DiagBitBoard,
		neg3DiagBitBoard,
		neg4DiagBitBoard,
		neg5DiagBitBoard,
		neg6DiagBitBoard,
		neg7DiagBitBoard,
		neg8DiagBitBoard,
		neg9DiagBitBoard,
		neg10DiagBitBoard,
		neg11DiagBitBoard,
		neg12DiagBitBoard,
		neg13DiagBitBoard,
		neg14DiagBitBoard,
		neg15DiagBitBoard,
	}

	for _, p := range pos {
		for _, n := range neg {
			sq := p & n

			if sq == 0 {
				continue
			}

			moves := (p | n) &^ sq

			bishopMovesLookupTable[sq] = moves
		}
	}

	if len(bishopMovesLookupTable) != 64 {
		panic("missing moves")
	}
}

func rookMoves(sq bitBoard) bitBoard {
	return rookMovesLookupTable[sq]
}

func bishopMoves(sq bitBoard) bitBoard {
	return bishopMovesLookupTable[sq]
}

func queenMoves(sq bitBoard) bitBoard {
	return rookMoves(sq) | bishopMoves(sq)
}
