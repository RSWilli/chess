package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/mm"
	"github.com/rswilli/chess/internal/pkg/notify"
	"github.com/rswilli/chess/internal/player"
	"github.com/rswilli/chess/internal/sse"
	"github.com/rswilli/chess/internal/www"
)

type Handler struct {
	http.Handler
	lock sync.Mutex

	eventid int

	engines player.EnginesRegistry

	broker *notify.Broker[sse.Event]

	// gamecond is used to notify the game loop about game changes
	gamecond *sync.Cond
	// game is nil when no game is running currently
	game *mm.Game
}

func NewHandler(r player.EnginesRegistry) *Handler {
	mux := http.NewServeMux()

	h := &Handler{
		Handler: mux,
		engines: r,
		broker:  notify.NewBroker[sse.Event](),
	}
	h.gamecond = sync.NewCond(&h.lock)

	mux.HandleFunc("POST /game/new", h.handleNewGame)

	mux.HandleFunc("GET /events", h.handleSubscriber)
	mux.HandleFunc("POST /move/{move}", h.handleMove)
	mux.HandleFunc("GET /{$}", h.handleIndex)
	mux.Handle("GET /", www.StaticServer)

	go h.runGameLoop()

	return h
}

func (h *Handler) handleSubscriber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	c := make(chan sse.Event)
	h.broker.Register(c)
	defer h.broker.Deregister(c)

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
}

// getHuman returns the current player if it is a Human. Must be called with the lock held
func (h *Handler) getHuman() *player.Human {
	if h.game == nil {
		return nil
	}

	current := h.game.Current()

	human, ok := current.(*player.Human)

	if !ok {
		return nil
	}

	return human
}

func (h *Handler) handleMove(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	human := h.getHuman()

	if human == nil {
		http.Error(w, "human is not in turn", http.StatusBadRequest)
		return
	}

	m := r.PathValue("move")

	err := human.DoMove(m)

	if err != nil {
		http.Error(w, fmt.Sprintf("could not do move: %v", err), http.StatusBadRequest)
		return
	}

	h.broadcastChange()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// getPosition returns the current position or nil if no game is active. Must be called with the lock held
func (h *Handler) getPosition() *chess.Position {
	if h.game == nil {
		return nil
	}

	return h.game.Position
}

// renderData returns the data object needed to render the current board
func (h *Handler) renderBoardData(selectedSquare chess.Square) www.BoardData {
	position := h.getPosition()
	sources := make(map[chess.Square]struct{})
	targets := make(map[chess.Square][]chess.Move)

	if position != nil {
		moves := position.GenerateMoves()

		for _, m := range moves {
			sources[m.From] = struct{}{}

			if m.From == selectedSquare {
				targets[m.To] = append(targets[m.To], m)
			}
		}
	}

	return www.BoardData{
		Position:    position,
		Selected:    selectedSquare,
		MoveSources: sources,
		MoveTargets: targets,
	}
}

func (h *Handler) handleNewGame(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	white := r.FormValue("white")
	black := r.FormValue("black")

	whitePlayer, err := h.engines.NewEngine(white)

	if err != nil {
		http.Error(w, "invalid white player", http.StatusBadRequest)
		return
	}

	blackPlayer, err := h.engines.NewEngine(black)

	if err != nil {
		http.Error(w, "invalid black player", http.StatusBadRequest)
		return
	}

	if h.game != nil {
		h.game.Stop()
	}

	h.game = mm.NewGame(whitePlayer, blackPlayer)
	h.gamecond.Signal()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	var square chess.Square
	sq := r.URL.Query().Get("square")

	if sq != "" {
		human := h.getHuman()
		if human == nil {
			slog.Warn("human is not in turn")
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		pos := h.getPosition()

		if pos == nil {
			slog.Warn("gamne not started")
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		parsed, err := chess.ParseSquare(sq)

		if err != nil {
			slog.Warn("invalid square given", "square", sq)
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		valid := false

		for _, m := range pos.GenerateMoves() {
			if m.From == parsed {
				valid = true
				break
			}
		}

		if !valid {
			slog.Warn("invalid square given", "square", sq)
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		square = parsed
	}

	h.renderIndex(w, h.renderBoardData(square))
}

func (h *Handler) renderIndex(w http.ResponseWriter, data www.BoardData) {
	engines, err := h.engines.AvailableEngines()

	if err != nil {
		slog.Error("could not get engines", "error", err)
		http.Error(w, "could not get engines", http.StatusInternalServerError)
		return
	}

	err = www.RenderIndex(w, www.Data{
		Board: data,
		Controls: www.ControlData{
			Engines: engines,
		},
	})

	if err != nil {
		slog.Error("error rendering template", "error", err)
		http.Error(w, "could not render template", http.StatusInternalServerError)
	}
}

func (h *Handler) broadcastChange() {
	ev := sse.Event{
		ID:    fmt.Sprintf("id-%d", h.eventid),
		Event: "change",
		Data:  struct{}{},
	}

	h.eventid++

	h.broker.Send(ev)
}

func (h *Handler) runGameLoop() {
	for {
		h.lock.Lock()

		for h.game == nil || h.game.State() != mm.Running {
			h.gamecond.Wait()
		}

		// keep a local copy so others can modify the game in the handler
		game := h.game
		h.lock.Unlock()

		// if a player is a human, then we need to be able to take the lock while the engine is thinking
		err := game.Move()

		h.lock.Lock()

		if err != nil {
			slog.Error("error while performing move in game", "error", err)
		}

		h.broadcastChange()
		h.lock.Unlock()
	}
}
