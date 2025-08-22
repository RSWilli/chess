package chess

func (p *Position) GenerateMoves() []Move {
	if p.possibleMoves != nil {
		return p.possibleMoves
	}

	p.computeAll()

	// TODO: 50 move rule, draw by material, draw by repetition
	p.possibleMoves = make([]Move, 0, maxMoveCount)

	if p.playerInTurn == White {

		p.generateKingMoves(p.whiteKing)
		p.whitePawns.Each(p.generateWhitePawnMoves)
		p.whiteKnights.Each(p.generateKnightMoves)

		p.whiteBishops.Each(p.generateBishopMoves)
		p.whiteRooks.Each(p.generateRookMoves)
		p.whiteQueens.Each(p.generateQueenMoves)
	} else {
		p.generateKingMoves(p.blackKing)
		p.blackPawns.Each(p.generateBlackPawnMoves)
		p.blackKnights.Each(p.generateKnightMoves)

		p.blackBishops.Each(p.generateBishopMoves)
		p.blackRooks.Each(p.generateRookMoves)
		p.blackQueens.Each(p.generateQueenMoves)
	}

	return p.possibleMoves
}

func (p *Position) generatePawnMoves(from, pushed, doublePushed, doublePushRank, promoteRank, opposing BitBoard) {
	all := p.allPieces()

	if from&doublePushRank != 0 && all&pushed == 0 && all&doublePushed == 0 {
		// can double push

		m := Move{
			From:    Square(from),
			To:      Square(doublePushed),
			Special: DoublePawnPush,
		}

		if p.isLegalMove(m) {
			p.possibleMoves = append(p.possibleMoves, m)
		}
	}

	// moves tracks the possible moves temporarily to defer
	// promotion decision
	moves := make([]Move, 0, 3)

	if all&pushed == 0 {
		// can push

		m := Move{
			From:    Square(from),
			To:      Square(pushed),
			Special: NoSpecial,
		}

		if p.isLegalMove(m) {
			moves = append(moves, m)
		}
	}

	takes := []BitBoard{
		pushed.Left(),
		pushed.Right(),
	}

	for _, t := range takes {
		if t == 0 {
			// edge of the board
			continue
		}

		if opposing&t == 0 && p.enPassantTarget != Square(t) {
			continue
		}

		taken := Empty
		special := Captures

		if p.enPassantTarget == Square(t) {
			taken = Pawn
			special |= EnPassant
		} else {
			taken = p.Square(Square(t))
		}

		m := Move{
			From:    Square(from),
			To:      Square(t),
			Special: special,
			Takes:   taken,
		}

		if p.isLegalMove(m) {
			moves = append(moves, m)
		}
	}

	if from&promoteRank == 0 {
		// no promotion
		p.possibleMoves = append(p.possibleMoves, moves...)
		return
	}

	promotions := []MoveSpecial{
		PromoteQueen,
		PromoteRook,
		PromoteBishop,
		PromoteKnight,
	}

	for _, prom := range promotions {
		for _, m := range moves {
			m.Special = m.Special | prom
			p.possibleMoves = append(p.possibleMoves, m)
		}
	}
}

const rank7BitBoard BitBoard = 0xff00
const rank2BitBoard BitBoard = 0xff000000000000

func (p *Position) generateWhitePawnMoves(bb BitBoard) {
	p.generatePawnMoves(bb, bb.Up(), bb.Up().Up(), rank2BitBoard, rank7BitBoard, p.blackPieces())
}

func (p *Position) generateBlackPawnMoves(bb BitBoard) {
	p.generatePawnMoves(bb, bb.Down(), bb.Down().Down(), rank7BitBoard, rank2BitBoard, p.whitePieces())
}

// bitboard to check when trying to castle
var (
	whiteCastleKing            = BitBoard(MustParseSquare("f1") | MustParseSquare("g1"))
	whiteCastleKingKingTarget  = MustParseSquare("g1")
	whiteCastleKingRookTarget  = MustParseSquare("f1")
	whiteCastleQueen           = BitBoard(MustParseSquare("c1") | MustParseSquare("d1"))
	whiteCastleQueenKingTarget = MustParseSquare("c1")
	whiteCastleQueenRookTarget = MustParseSquare("d1")

	blackCastleQueen           = BitBoard(MustParseSquare("c8") | MustParseSquare("d8"))
	blackCastleKingKingTarget  = MustParseSquare("g8")
	blackCastleKingRookTarget  = MustParseSquare("f8")
	blackCastleKing            = BitBoard(MustParseSquare("f8") | MustParseSquare("g8"))
	blackCastleQueenKingTarget = MustParseSquare("c8")
	blackCastleQueenRookTarget = MustParseSquare("d8")
)

// notAttacked returns true if none of the squares set in the given bitboard are attacked
func (p *Position) notAttacked(bb BitBoard) bool {
	for sq := range bb.Ones() {
		if p.attacksTo.get(sq) != 0 {
			return false
		}
	}
	return true
}

func (p *Position) canCastleWhiteKing() bool {
	return p.whiteKing == BitBoard(e1) && p.castling.Has(CastleWhiteKing) && p.all&whiteCastleKing == 0 && p.notAttacked(whiteCastleKing)
}

func (p *Position) canCastleWhiteQueen() bool {
	return p.whiteKing == BitBoard(e1) && p.castling.Has(CastleWhiteQueen) && p.all&whiteCastleQueen == 0 && p.notAttacked(whiteCastleQueen)
}
func (p *Position) canCastleBlackKing() bool {
	return p.blackKing == BitBoard(e8) && p.castling.Has(CastleBlackKing) && p.all&blackCastleKing == 0 && p.notAttacked(blackCastleKing)
}

func (p *Position) canCastleBlackQueen() bool {
	return p.blackKing == BitBoard(e8) && p.castling.Has(CastleBlackQueen) && p.all&blackCastleQueen == 0 && p.notAttacked(blackCastleQueen)
}

func (p *Position) generateKingMoves(bb BitBoard) {
	if bb.Count() != 1 {
		panic("expected 1 king")
	}

	for t := range kingMoves(bb).Ones() {
		if t == 0 {
			// wrapped around
			continue
		}

		if p.attacksTo.get(t) != 0 {
			continue
		}

		if t&p.ours != 0 {
			// occupied
			continue
		}

		s := NoSpecial
		taken := Empty

		if t&p.theirs != 0 {
			s |= Captures
			taken = p.Square(Square(t))
		}

		p.possibleMoves = append(p.possibleMoves, Move{
			From:    Square(bb),
			To:      Square(t),
			Special: s,
			Takes:   taken,
		})
	}

	if p.castling == NoCastling {
		return
	}

	if p.attacksTo.get(bb) != 0 {
		// king in check, no castling allowed
		return
	}

	// white O-O, aka e1g1
	if p.playerInTurn == White && p.canCastleWhiteKing() {
		p.possibleMoves = append(p.possibleMoves, Move{
			From:    e1,
			To:      whiteCastleKingKingTarget,
			Special: CastleShort,
		})
	}

	// white O-O-O, aka e1c1
	if p.playerInTurn == White && p.canCastleWhiteQueen() {
		p.possibleMoves = append(p.possibleMoves, Move{
			From:    e1,
			To:      whiteCastleQueenKingTarget,
			Special: CastleLong,
		})
	}

	// black O-O, aka e8g8
	if p.playerInTurn == Black && p.canCastleBlackKing() {
		p.possibleMoves = append(p.possibleMoves, Move{
			From:    e8,
			To:      blackCastleKingKingTarget,
			Special: CastleShort,
		})
	}

	// black O-O-O, aka e8c8
	if p.playerInTurn == Black && p.canCastleBlackQueen() {
		p.possibleMoves = append(p.possibleMoves, Move{
			From:    e8,
			To:      blackCastleQueenKingTarget,
			Special: CastleLong,
		})
	}
}

func (p *Position) generateKnightMoves(knight BitBoard) {
	for t := range knightMoves(knight).Ones() {
		if t == 0 {
			// wrapped around
			continue
		}

		if t&p.ours != 0 {
			// occupied
			continue
		}

		s := NoSpecial
		taken := Empty

		if t&p.theirs != 0 {
			s |= Captures
			taken = p.Square(Square(t))
		}

		m := Move{
			From:    Square(knight),
			To:      Square(t),
			Special: s,
			Takes:   taken,
		}

		if !p.isLegalMove(m) {
			continue
		}

		p.possibleMoves = append(p.possibleMoves, m)
	}
}

func (p *Position) generateRookMoves(rook BitBoard) {
	targets := rookMoves(rook, p.ours, p.theirs)

	for t := range targets.Ones() {
		s := NoSpecial
		taken := Empty

		if t&p.theirs != 0 {
			s |= Captures
			taken = p.Square(Square(t))
		}

		m := Move{
			From:    Square(rook),
			To:      Square(t),
			Special: s,
			Takes:   taken,
		}

		if !p.isLegalMove(m) {
			continue
		}

		p.possibleMoves = append(p.possibleMoves, m)
	}
}

func (p *Position) generateQueenMoves(queen BitBoard) {
	targets := queenMoves(queen, p.ours, p.theirs)

	for t := range targets.Ones() {
		s := NoSpecial
		taken := Empty

		if t&p.theirs != 0 {
			s |= Captures
			taken = p.Square(Square(t))
		}

		m := Move{
			From:    Square(queen),
			To:      Square(t),
			Special: s,
			Takes:   taken,
		}

		if !p.isLegalMove(m) {
			continue
		}

		p.possibleMoves = append(p.possibleMoves, m)
	}
}

func (p *Position) generateBishopMoves(bishop BitBoard) {
	targets := bishopMoves(bishop, p.ours, p.theirs)

	for t := range targets.Ones() {
		s := NoSpecial
		taken := Empty

		if t&p.theirs != 0 {
			s |= Captures
			taken = p.Square(Square(t))
		}

		m := Move{
			From:    Square(bishop),
			To:      Square(t),
			Special: s,
			Takes:   taken,
		}

		if !p.isLegalMove(m) {
			continue
		}

		p.possibleMoves = append(p.possibleMoves, m)
	}
}

// isLegalMove checks if the move would leave the king in check. This does not correctly
// check king moves or en passant captures.
func (p *Position) isLegalMove(m Move) bool {
	var ourKing BitBoard

	if p.playerInTurn == White {
		ourKing = p.whiteKing
	} else {
		ourKing = p.blackKing
	}

	// move the piece, needed for pin detection:
	moved := (p.all | BitBoard(m.To)) &^ BitBoard(m.From)

	for _, attack := range p.xRayKingAttacks {
		if BitBoard(m.To) != attack.from && moved&attack.ray == 0 {
			// no piece is blocking the sliding piece attack or taking it, piece could have moved
			// out of a pin, or did not block the check -> not legal
			return false
		}
	}

	kingAttacks := p.attacksTo.get(ourKing)

	if m.Special.Has(Captures) && (kingAttacks&^BitBoard(m.To)) == 0 {
		// the only attacker attacking the king is taken, so this is a legal move
		return true
	}

	if kingAttacks == 0 && len(p.xRayKingAttacks) == 0 {
		// king is not attacked and not xrayed in any way, so the move is legal.
		return true
	}

	// the move is only legal if the king is not attacked
	return kingAttacks == 0
}
