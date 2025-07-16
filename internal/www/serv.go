package www

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"log/slog"
	"net/http"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/game"
)

//go:embed all:static
var static embed.FS

//go:embed *.html
var htmlFiles embed.FS
var templates = template.Must(template.New("main").Funcs(funcMap).ParseFS(htmlFiles, "*.html"))

func Render(board *chess.Board) []byte {
	buf := bytes.Buffer{}
	templates.ExecuteTemplate(&buf, "board.html", board)

	return buf.Bytes()
}

func Handler(state *game.State) http.Handler {
	staticServer := http.FileServerFS(static)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("serving request", "path", r.URL.Path)

		if r.URL.Path == "/" {
			file, err := htmlFiles.Open("index.html")

			if err != nil {
				panic("could not open index")
			}

			io.Copy(w, file)

			return
		}

		staticServer.ServeHTTP(w, r)
	})
}

var funcMap = template.FuncMap(map[string]any{
	"ranks":   ranks,
	"files":   files,
	"color":   color,
	"pieceAt": pieceAt,
})

// pieceAt returns the URL of the image that is needed for the piece at file/rank or an empty string
func pieceAt(board *chess.Board, file, rank int) string {
	piece := board.Square(chess.Square{
		File: chess.File(file),
		Rank: chess.Rank(rank),
	})

	if piece == chess.Empty {
		return ""
	}

	return pieceImgSrc[piece]
}

// ranks return the ranks of the board in the order they are needed in the HTML
func ranks() []string {
	return []string{"8", "7", "6", "5", "4", "3", "2", "1"}
}

// files returns the files of the board in the order they are needed in the HTML
func files() []string {
	return []string{"a", "b", "c", "d", "e", "f", "g", "h"}
}

// color returns "black" or "white" depending on the given indices of the ranks and file lists
func color(rank, file int) string {
	if (rank+file)%2 == 1 {
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
