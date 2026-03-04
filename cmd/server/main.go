package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/engine"
	"github.com/rswilli/chess/internal/mm"
	"github.com/rswilli/chess/internal/pkg/notify"
	"github.com/rswilli/chess/internal/player"
	"github.com/rswilli/chess/internal/sse"
	"github.com/rswilli/chess/internal/uci"
	"github.com/rswilli/chess/internal/www"
)

func main() {
	mux := http.NewServeMux()

	human := player.NewHuman()

	engine := engine.NewEngine()

	game := mm.NewGame(uci.NewProxy(human), uci.NewProxy(engine))

	b := notify.NewBroker[sse.Event]()

	h := Handler{game: game, b: b}

	mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		c := make(chan sse.Event)
		b.Register(c)
		defer b.Deregister(c)

		flusher.Flush()

		for {
			select {
			case event := <-c:
				w.Write(event.Bytes())
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	})

	mux.HandleFunc("PUT /move/{move}", func(w http.ResponseWriter, r *http.Request) {
		m := r.PathValue("move")

		err := human.DoMove(m)

		if err != nil {
			http.Error(w, fmt.Sprintf("could not do move: %v", err), http.StatusBadRequest)
			return
		}

		h.BroadcastBoard(chess.InvalidSquare)

		w.Write([]byte("ok"))
	})

	mux.HandleFunc("PUT /square/{square}", func(w http.ResponseWriter, r *http.Request) {
		sq := r.PathValue("square")

		square, err := chess.ParseSquare(sq)

		if err != nil {
			http.Error(w, "invalid square", http.StatusBadRequest)
			return
		}

		human.DoSquare(square)

		h.BroadcastBoard(human.CurrentSquare())

		w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		position := game.Position()

		err := www.RenderIndex(w, www.Data{
			Board: position,
		})

		if err != nil {
			slog.Error("error rendering template", "error", err)
			http.Error(w, "could not render template", http.StatusInternalServerError)
		}
	})

	mux.Handle("GET /", www.StaticServer)

	go func() {
		for {
			err := game.Move()

			h.BroadcastBoard(chess.InvalidSquare)

			if err != nil {
				slog.Info("error from game", "error", err)
				return
			}

		}
	}()

	serv := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	slog.Info("starting server", "addr", "localhost:3000")
	err := serv.ListenAndServe()

	slog.Error("server ended", "error", err)
}
