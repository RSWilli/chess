package www

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/rswilli/chess/internal/chess"
)

//go:embed static/*
var static embed.FS
var StaticServer = http.FileServerFS(static)

//go:embed *.tpl.html
var rawtemplates embed.FS
var templates = template.Must(template.New("main").Funcs(funcMap).ParseFS(rawtemplates, "*"))

func RenderIndex(w io.Writer, data Data) error {
	return templates.ExecuteTemplate(w, "index.tpl.html", data)
}

func RenderBoard(w io.Writer, data BoardData) error {
	return templates.ExecuteTemplate(w, "board.tpl.html", data)
}

type Data struct {
	Board    BoardData
	Controls ControlData
}

type ControlData struct {
	Engines []string
}

type BoardData struct {
	Position *chess.Position
	Selected chess.Square
	// MoveSources contains all squares from which a move can be performed
	MoveSources map[chess.Square]struct{}
	// MoveTargets contains target squares and their moves from a selected start square
	MoveTargets map[chess.Square][]chess.Move
}

// ClassesFor returns the HTML classes for the file and rank. Intended to be called from the template
func (d BoardData) ClassesFor(sq string) string {
	var classes []string
	square, err := chess.ParseSquare(sq)

	if err != nil {
		panic(err)
	}

	if (square.Rank()+square.File())%2 == 1 {
		classes = append(classes, "black")
	} else {
		classes = append(classes, "white")
	}

	if d.Position == nil {
		// uninitialized game:
		return strings.Join(classes, " ")
	}

	piece := d.Position.Square(square)

	if d.Selected == square {
		classes = append(classes, "highlighted")
	}

	if _, ok := d.MoveTargets[square]; ok {
		classes = append(classes, "target")
	}

	if (piece == (chess.King | d.Position.PlayerInTurn)) && d.Position.IsCheck() {
		classes = append(classes, "check")
	}

	if piece == (chess.King|d.Position.PlayerInTurn) && d.Position.IsCheckMate() {
		classes = append(classes, "loose")
	}

	if (piece.Type() == chess.King) && piece != (chess.King|d.Position.PlayerInTurn) && d.Position.IsCheckMate() {
		classes = append(classes, "win")
	}

	if (piece.Type() == chess.King) && d.Position.IsDraw() {
		classes = append(classes, "draw")
	}

	return strings.Join(classes, " ")
}

func (d BoardData) CanMoveFrom(sq string) bool {
	square, err := chess.ParseSquare(sq)

	if err != nil {
		panic(err)
	}

	_, ok := d.MoveSources[square]

	return ok
}

func (d BoardData) RegularMoveTo(sq string) string {
	square, err := chess.ParseSquare(sq)

	if err != nil {
		panic(err)
	}

	moves := d.MoveTargets[square]

	if len(moves) != 1 {
		// no moves or promotions
		return ""
	}

	return moves[0].String()
}

func (d BoardData) PromotionMovesTo(sq string) []chess.Move {
	square, err := chess.ParseSquare(sq)

	if err != nil {
		panic(err)
	}

	moves := d.MoveTargets[square]

	if len(moves) <= 1 {
		// needs multiple moves to the same square for promotion move
		return nil
	}

	return moves
}

// PieceAt returns the URL of the image that is needed for the piece at file/rank or an empty string
func (d BoardData) PieceAt(sq string) string {
	if d.Position == nil {
		return ""
	}

	square, err := chess.ParseSquare(sq)

	if err != nil {
		panic(err)
	}

	piece := d.Position.Square(square)

	if piece == chess.Empty {
		return ""
	}

	return pieceImgSrc[piece]
}

// PromotionPiece returns the URL of the image that corresponds to the given promotion move
func (d BoardData) PromotionPiece(m chess.Move) string {
	switch m.Special & chess.PromoteAny {
	case chess.PromoteBishop:
		return pieceImgSrc[d.Position.PlayerInTurn|chess.Bishop]
	case chess.PromoteKnight:
		return pieceImgSrc[d.Position.PlayerInTurn|chess.Knight]
	case chess.PromoteQueen:
		return pieceImgSrc[d.Position.PlayerInTurn|chess.Queen]
	case chess.PromoteRook:
		return pieceImgSrc[d.Position.PlayerInTurn|chess.Rook]
	default:
		panic(fmt.Sprintf("unexpected chess.MoveSpecial: %#v", m.Special))
	}
}

var funcMap = template.FuncMap(map[string]any{
	"ranks":          ranks,
	"files":          files,
	"startpositions": startpositions,
})

// ranks return the ranks of the board in the order they are needed in the HTML
func ranks() []string {
	return []string{"8", "7", "6", "5", "4", "3", "2", "1"}
}

// files returns the files of the board in the order they are needed in the HTML
func files() []string {
	return []string{"a", "b", "c", "d", "e", "f", "g", "h"}
}

func startpositions() []startPosition {
	return startPositions
}

var pieceImgSrc = map[chess.Piece]string{
	chess.WhitePawn:   "/static/pw.svg",
	chess.WhiteKnight: "/static/nw.svg",
	chess.WhiteBishop: "/static/bw.svg",
	chess.WhiteRook:   "/static/rw.svg",
	chess.WhiteQueen:  "/static/qw.svg",
	chess.WhiteKing:   "/static/kw.svg",

	chess.BlackPawn:   "/static/pb.svg",
	chess.BlackKnight: "/static/nb.svg",
	chess.BlackBishop: "/static/bb.svg",
	chess.BlackRook:   "/static/rb.svg",
	chess.BlackQueen:  "/static/qb.svg",
	chess.BlackKing:   "/static/kb.svg",
}
