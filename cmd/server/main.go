package main

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/rswilli/chess/internal/game"
	"github.com/rswilli/chess/internal/www"
	"golang.org/x/net/websocket"
)

func main() {
	state := game.New()

	mux := http.NewServeMux()

	ws := websocket.Server{
		Handshake: func(c *websocket.Config, r *http.Request) error {
			return nil
		},
		Handler: websocket.Handler(func(ws *websocket.Conn) {
			slog.Info("subscribing to events")
			events := state.Subscribe()
			defer state.Unsubscribe(events)

			done := make(chan struct{})

			go func() {
				// this blocks until the client closes the socket:
				io.Copy(io.Discard, ws)
				done <- struct{}{}
			}()

			for {
				boardData := www.Render(state.Board())
				// needs to happen as a single call to write otherwise we create multiple
				// messages
				ws.Write(boardData)

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

	mux.Handle("GET /websocket", ws)

	mux.Handle("GET /", www.Handler(state))

	serv := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	slog.Info("starting server", "addr", "localhost:3000")
	serv.ListenAndServe()
}
