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

func RenderBoard(w io.Writer, data Data) error {
	return templates.ExecuteTemplate(w, "board.tpl.html", data)
}

type Data struct {
	Board       *chess.Position
	Selected    chess.Square
	MoveTargets map[chess.Square][]chess.Move
	Promotion   bool
}

// ClassesFor returns the HTML classes for the file and rank. Intended to be called from the template
func (d Data) ClassesFor(fileIndex, rankIndex int) string {
	var classes []string

	if (rankIndex+fileIndex)%2 == 1 {
		classes = append(classes, "black")
	} else {
		classes = append(classes, "white")
	}

	square := chess.NewSquare(rankIndex, fileIndex)
	piece := d.Board.Square(square)

	if d.Selected == square {
		classes = append(classes, "highlighted")
	}

	if _, ok := d.MoveTargets[square]; ok {
		classes = append(classes, "target")

		if d.Promotion {
			classes = append(classes, "promotion")
		}
	}

	if (piece == (chess.King | d.Board.PlayerInTurn)) && d.Board.IsCheck() {
		classes = append(classes, "check")
	}

	return strings.Join(classes, " ")
}

func (d Data) MoveTo(fileIndex, rankIndex int) string {
	sq := chess.NewSquare(rankIndex, fileIndex)

	moves, ok := d.MoveTargets[sq]

	if !ok {
		return ""
	}

	var ms []string

	for _, m := range moves {
		ms = append(ms, m.String())
	}

	return strings.Join(ms, ",")
}

// PieceAt returns the URL of the image that is needed for the piece at file/rank or an empty string
func (d Data) PieceAt(fileIndex, rankIndex int) string {
	piece := d.Board.Square(chess.NewSquare(rankIndex, fileIndex))

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
