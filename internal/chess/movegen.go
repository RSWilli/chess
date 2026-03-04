package chess

import "fmt"

func (p *Position) GenerateMoves() []Move {
	// TODO: 50 move rule, draw by material, draw by repetition

	legalMoves := make([]Move, 0, maxMoveCount)

	for sq := range p.ours().Ones() {
		piece := p.Square(Square(sq))
		p.generateMovesForPiece(piece, sq, &legalMoves)
	}

	return legalMoves
}

func (p *Position) generateMovesForPiece(piece Piece, at BitBoard, legalMoves *[]Move) {
	switch piece {
	case WhitePawn:
		p.generateWhitePawnMoves(at, legalMoves)
		return
	case BlackPawn:
		p.generateBlackPawnMoves(at, legalMoves)
		return
	}

	switch piece &^ (White | Black) {
	case Bishop:
		p.generateBishopMoves(at, legalMoves)
	case King:
		p.generateKingMoves(at, legalMoves)
	case Knight:
		p.generateKnightMoves(at, legalMoves)
	case Queen:
		p.generateQueenMoves(at, legalMoves)
	case Rook:
		p.generateRookMoves(at, legalMoves)
	default:
		panic(fmt.Sprintf("unexpected chess.Piece to generate moves for: %#v", piece))
	}
}

func (p *Position) generatePawnMoves(from, pushed, doublePushed, doublePushRank, promoteRank, opposing BitBoard, legalMoves *[]Move) {
	all := p.all()

	if from&doublePushRank != 0 && all&pushed == 0 && all&doublePushed == 0 {
		// can double push

		m := Move{
			From:    Square(from),
			To:      Square(doublePushed),
			Special: DoublePawnPush,
		}

		if p.isLegalMove(m) {
			*legalMoves = append(*legalMoves, m)
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
		*legalMoves = append(*legalMoves, moves...)
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
			*legalMoves = append(*legalMoves, m)
		}
	}
}

const rank7BitBoard BitBoard = 0xff00
const rank2BitBoard BitBoard = 0xff000000000000

func (p *Position) generateWhitePawnMoves(bb BitBoard, legalMoves *[]Move) {
	p.generatePawnMoves(bb, bb.Up(), bb.Up().Up(), rank2BitBoard, rank7BitBoard, p.blackPieces(), legalMoves)
}

func (p *Position) generateBlackPawnMoves(bb BitBoard, legalMoves *[]Move) {
	p.generatePawnMoves(bb, bb.Down(), bb.Down().Down(), rank7BitBoard, rank2BitBoard, p.whitePieces(), legalMoves)
}

// bitboard to check when trying to castle
var (
	whiteCastleKing            = BitBoard(MustParseSquare("f1") | MustParseSquare("g1"))
	whiteCastleKingKingTarget  = MustParseSquare("g1")
	whiteCastleKingRookTarget  = MustParseSquare("f1")
	whiteCastleQueenNoCheck    = BitBoard(MustParseSquare("c1") | MustParseSquare("d1"))
	whiteCastleQueenEmpty      = BitBoard(MustParseSquare("b1") | MustParseSquare("c1") | MustParseSquare("d1"))
	whiteCastleQueenKingTarget = MustParseSquare("c1")
	whiteCastleQueenRookTarget = MustParseSquare("d1")

	blackCastleQueenNoCheck    = BitBoard(MustParseSquare("c8") | MustParseSquare("d8"))
	blackCastleQueenEmpty      = BitBoard(MustParseSquare("b8") | MustParseSquare("c8") | MustParseSquare("d8"))
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
	return p.whiteKing == BitBoard(e1) && p.castling.Has(CastleWhiteKing) && p.all()&whiteCastleKing == 0 && p.notAttacked(whiteCastleKing)
}

func (p *Position) canCastleWhiteQueen() bool {
	return p.whiteKing == BitBoard(e1) && p.castling.Has(CastleWhiteQueen) && p.all()&whiteCastleQueenEmpty == 0 && p.notAttacked(whiteCastleQueenNoCheck)
}
func (p *Position) canCastleBlackKing() bool {
	return p.blackKing == BitBoard(e8) && p.castling.Has(CastleBlackKing) && p.all()&blackCastleKing == 0 && p.notAttacked(blackCastleKing)
}

func (p *Position) canCastleBlackQueen() bool {
	return p.blackKing == BitBoard(e8) && p.castling.Has(CastleBlackQueen) && p.all()&blackCastleQueenEmpty == 0 && p.notAttacked(blackCastleQueenNoCheck)
}

func (p *Position) generateKingMoves(bb BitBoard, legalMoves *[]Move) {
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

		if t&p.ours() != 0 {
			// occupied
			continue
		}

		s := NoSpecial
		taken := Empty

		if t&p.theirs() != 0 {
			s |= Captures
			taken = p.Square(Square(t))
		}

		*legalMoves = append(*legalMoves, Move{
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
	if p.PlayerInTurn == White && p.canCastleWhiteKing() {
		*legalMoves = append(*legalMoves, Move{
			From:    e1,
			To:      whiteCastleKingKingTarget,
			Special: CastleShort,
		})
	}

	// white O-O-O, aka e1c1
	if p.PlayerInTurn == White && p.canCastleWhiteQueen() {
		*legalMoves = append(*legalMoves, Move{
			From:    e1,
			To:      whiteCastleQueenKingTarget,
			Special: CastleLong,
		})
	}

	// black O-O, aka e8g8
	if p.PlayerInTurn == Black && p.canCastleBlackKing() {
		*legalMoves = append(*legalMoves, Move{
			From:    e8,
			To:      blackCastleKingKingTarget,
			Special: CastleShort,
		})
	}

	// black O-O-O, aka e8c8
	if p.PlayerInTurn == Black && p.canCastleBlackQueen() {
		*legalMoves = append(*legalMoves, Move{
			From:    e8,
			To:      blackCastleQueenKingTarget,
			Special: CastleLong,
		})
	}
}

func (p *Position) generateKnightMoves(knight BitBoard, legalMoves *[]Move) {
	for t := range knightMoves(knight).Ones() {
		if t == 0 {
			// wrapped around
			continue
		}

		if t&p.ours() != 0 {
			// occupied
			continue
		}

		s := NoSpecial
		taken := Empty

		if t&p.theirs() != 0 {
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

		*legalMoves = append(*legalMoves, m)
	}
}

func (p *Position) generateRookMoves(rook BitBoard, legalMoves *[]Move) {
	targets := rookMoves(rook, p.ours(), p.theirs())

	for t := range targets.Ones() {
		s := NoSpecial
		taken := Empty

		if t&p.theirs() != 0 {
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

		*legalMoves = append(*legalMoves, m)
	}
}

func (p *Position) generateQueenMoves(queen BitBoard, legalMoves *[]Move) {
	targets := queenMoves(queen, p.ours(), p.theirs())

	for t := range targets.Ones() {
		s := NoSpecial
		taken := Empty

		if t&p.theirs() != 0 {
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

		*legalMoves = append(*legalMoves, m)
	}
}

func (p *Position) generateBishopMoves(bishop BitBoard, legalMoves *[]Move) {
	targets := bishopMoves(bishop, p.ours(), p.theirs())

	for t := range targets.Ones() {
		s := NoSpecial
		taken := Empty

		if t&p.theirs() != 0 {
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

		*legalMoves = append(*legalMoves, m)
	}
}

// isLegalMove checks if the move would leave the king in check. This does not correctly
// check king moves or en passant captures.
func (p *Position) isLegalMove(m Move) bool {
	var ourKing BitBoard

	if p.PlayerInTurn == White {
		ourKing = p.whiteKing
	} else {
		ourKing = p.blackKing
	}

	// enpassantTaken is the taken pawn while doing en passant
	var enpassantTaken BitBoard
	if m.Special.Has(EnPassant) {
		if p.PlayerInTurn == White {
			enpassantTaken = BitBoard(p.enPassantTarget).Down()
		} else {
			enpassantTaken = BitBoard(p.enPassantTarget).Up()
		}
	}

	// move the piece, needed for pin detection:
	moved := (p.all() | BitBoard(m.To)) &^ (BitBoard(m.From) | enpassantTaken)

	var blockedSlidingAttackers BitBoard

	for _, attack := range p.xRayKingAttacks {
		if attack.from != 0 && BitBoard(m.To) != attack.from && moved&attack.ray == 0 {
			// no piece is blocking the sliding piece attack or taking it, piece could have moved
			// out of a pin, or did not block the check -> not legal
			return false
		}

		blockedSlidingAttackers |= attack.from
	}

	// omit the attacks we blocked earlier
	kingAttacks := p.attacksTo.get(ourKing) &^ blockedSlidingAttackers

	if kingAttacks == 0 {
		return true
	}

	// move is only legal if we capture the attacker, either directly or via en passant
	return m.Special.Has(Captures) && (kingAttacks&^(BitBoard(m.To))) == 0 ||
		m.Special.Has(EnPassant) && (kingAttacks&^enpassantTaken) == 0
}
