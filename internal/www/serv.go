package www

import (
	"embed"
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
	// PromotionMove is the pawn move that will lead to a promotion, but without the choses piece type
	PromotionMove string
}

// ClassesFor returns the HTML classes for the file and rank. Intended to be called from the template
func (d BoardData) ClassesFor(fileIndex, rankIndex int) string {
	var classes []string

	if (rankIndex+fileIndex)%2 == 1 {
		classes = append(classes, "black")
	} else {
		classes = append(classes, "white")
	}

	if d.Position == nil {
		// uninitialized game:
		return strings.Join(classes, " ")
	}

	square := chess.NewSquare(rankIndex, fileIndex)
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

	if (piece&^(chess.White|chess.Black) == chess.King) && piece != (chess.King|d.Position.PlayerInTurn) && d.Position.IsCheckMate() {
		classes = append(classes, "win")
	}

	if (piece&^(chess.White|chess.Black) == chess.King) && d.Position.IsDraw() {
		classes = append(classes, "draw")
	}

	return strings.Join(classes, " ")
}

func (d BoardData) CanMoveFrom(fileIndex, rankIndex int) bool {
	sq := chess.NewSquare(rankIndex, fileIndex)

	_, ok := d.MoveSources[sq]

	return ok
}

func (d BoardData) MoveTo(fileIndex, rankIndex int) string {
	sq := chess.NewSquare(rankIndex, fileIndex)

	moves := d.MoveTargets[sq]

	if len(moves) == 0 {
		return ""
	}

	// remove the promotion flag:
	return moves[0].From.String() + moves[0].To.String()
}

func (d BoardData) PromotionTo(fileIndex, rankIndex int) bool {
	sq := chess.NewSquare(rankIndex, fileIndex)

	moves := d.MoveTargets[sq]

	// a promotion happens when the same target square has multiple possible moves:
	return len(moves) > 1
}

// PieceAt returns the URL of the image that is needed for the piece at file/rank or an empty string
func (d BoardData) PieceAt(fileIndex, rankIndex int) string {
	if d.Position == nil {
		return ""
	}

	piece := d.Position.Square(chess.NewSquare(rankIndex, fileIndex))

	if piece == chess.Empty {
		return ""
	}

	return pieceImgSrc[piece]
}

var funcMap = template.FuncMap(map[string]any{
	"ranks": ranks,
	"files": files,
	"color": color,
})

// ranks return the ranks of the board in the order they are needed in the HTML
func ranks() []string {
	return []string{"8", "7", "6", "5", "4", "3", "2", "1"}
}

// files returns the files of the board in the order they are needed in the HTML
func files() []string {
	return []string{"a", "b", "c", "d", "e", "f", "g", "h"}
}

// color returns "black" or "white" depending on the given indices of the ranks and file lists
func color(fileIndex, rankIndex int) string {
	if (rankIndex+fileIndex)%2 == 1 {
		return "black"
	}

	return "white"
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
