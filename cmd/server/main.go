package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/gui"
	"github.com/rswilli/chess/internal/www"
	"golang.org/x/net/websocket"
)

func main() {
	gui := gui.New()

	mux := http.NewServeMux()

	ws := websocket.Server{
		Handshake: func(c *websocket.Config, r *http.Request) error {
			return nil
		},
		Handler: websocket.Handler(func(ws *websocket.Conn) {
			slog.Info("subscribing to events")
			events := gui.Subscribe()
			defer gui.Unsubscribe(events)

			done := make(chan struct{})

			go func() {
				// this blocks until the client closes the socket:
				io.Copy(io.Discard, ws)
				done <- struct{}{}
			}()

			// buf to write to the websocket in a single write operation
			var buf bytes.Buffer

			for {
				buf.Reset()

				err := gui.Render(&buf)

				if err != nil {
					slog.Error("error rendering board", "error", err)
					return
				}

				n, err := ws.Write(buf.Bytes())

				if err != nil || n < len(buf.Bytes()) {
					slog.Error("write error of short write", "error", err)
					return
				}

				select {
				case <-ws.Request().Context().Done():
					return
				case <-done:
					return
				case <-events:
					// render again
				}
			}
		}),
	}

	mux.HandleFunc("PUT /square/{square}/{promotion}", func(w http.ResponseWriter, r *http.Request) {
		sq := r.PathValue("square")

		square, err := chess.ParseSquare(sq)

		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid square given: %v", err), http.StatusBadRequest)
			return
		}

		var special chess.MoveSpecial

		switch r.PathValue("promotion") {
		case "x":
			special = chess.NoSpecial
		case "q":
			special = chess.PromoteQueen
		case "r":
			special = chess.PromoteRook
		case "n":
			special = chess.PromoteKnight
		case "b":
			special = chess.PromoteBishop
		default:
			http.Error(w, fmt.Sprintf("Invalid promotion given: %s", r.PathValue("promotion")), http.StatusBadRequest)
		}

		gui.DoSquare(square, special)
	})

	mux.Handle("GET /websocket", ws)

	mux.Handle("GET /", www.StaticServer)

	serv := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	slog.Info("starting server", "addr", "localhost:3000")
	err := serv.ListenAndServe()

	slog.Error("server ended", "error", err)
}
