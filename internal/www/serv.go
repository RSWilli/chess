package www

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"

	"github.com/rswilli/chess/internal/chess"
)

//go:embed index.html static/*
var static embed.FS
var StaticServer = http.FileServerFS(static)

//go:embed board.tpl.html
var boardTpl string
var boardTemplate = template.Must(template.New("main").Funcs(funcMap).Parse(boardTpl))

func RenderBoard(data Data) ([]byte, error) {
	buf := bytes.Buffer{}
	err := boardTemplate.Execute(&buf, data)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Data struct {
	Board    *chess.Board
	Selected chess.Square
}

// IsSelected is called by the template for ease of use
func (d Data) IsSelected(fileIndex, rankIndex int) bool {
	return d.Selected == chess.Square{
		File: chess.File(fileIndex),
		Rank: chess.Rank(rankIndex),
	}
}

// pieceAt returns the URL of the image that is needed for the piece at file/rank or an empty string
func (d Data) PieceAt(fileIndex, rankIndex int) string {
	piece := d.Board.Square(chess.Square{
		File: chess.File(fileIndex),
		Rank: chess.Rank(rankIndex),
	})

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
func color(rankIndex, fileIndex int) string {
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
