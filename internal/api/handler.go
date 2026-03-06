package api

import (
	"bytes"
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

	broker *notify.Broker[sse.Event]

	// gamecond is used to notify the game loop about game changes
	gamecond *sync.Cond
	// game is nil when no game is running currently
	game *mm.Game
}

func NewHandler() *Handler {
	mux := http.NewServeMux()

	h := &Handler{
		Handler: mux,
		broker:  notify.NewBroker[sse.Event](),
	}
	h.gamecond = sync.NewCond(&h.lock)

	mux.HandleFunc("POST /game/new", h.handleNewGame)

	mux.HandleFunc("GET /events", h.handleSubscriber)
	mux.HandleFunc("PUT /move/{move}", h.handleMove)
	mux.HandleFunc("PUT /square/{square}", h.handleSquare)
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

	h.broadcastBoard()

	w.Write([]byte("ok"))
}

func (h *Handler) handleSquare(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	human := h.getHuman()

	if human == nil {
		http.Error(w, "human is not in turn", http.StatusBadRequest)
		return
	}

	sq := r.PathValue("square")

	square, err := chess.ParseSquare(sq)

	if err != nil {
		http.Error(w, "invalid square", http.StatusBadRequest)
		return
	}

	human.DoSquare(square)

	h.broadcastBoard()

	w.Write([]byte("ok"))
}

// getPosition returns the current position or nil if no game is active. Must be called with the lock held
func (h *Handler) getPosition() *chess.Position {
	if h.game == nil {
		return nil
	}

	return h.game.Position()
}

// renderData returns the data object needed to render the current board
func (h *Handler) renderBoardData() www.BoardData {
	position := h.getPosition()
	targets := make(map[chess.Square][]chess.Move)
	promotion := false

	var currentSquare chess.Square

	if h := h.getHuman(); h != nil {
		currentSquare = h.CurrentSquare()
	}

	if position != nil {
		moves := position.GenerateMoves()

		if currentSquare != chess.InvalidSquare {
			for _, m := range moves {
				if m.From == currentSquare {
					promotion = m.Special.Has(chess.PromoteAny)
					targets[m.To] = append(targets[m.To], m)
				}
			}
		}
	}

	return www.BoardData{
		Position:    position,
		Selected:    currentSquare,
		MoveTargets: targets,
		Promotion:   promotion,
	}
}

func (h *Handler) handleNewGame(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	white := r.FormValue("white")
	black := r.FormValue("black")

	whitePlayer, err := player.NewEngine(white)

	if err != nil {
		http.Error(w, "invalid white player", http.StatusBadRequest)
		return
	}

	blackPlayer, err := player.NewEngine(black)

	if err != nil {
		http.Error(w, "invalid black player", http.StatusBadRequest)
		return
	}

	h.game = mm.NewGame(whitePlayer, blackPlayer)
	h.gamecond.Signal()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	engines, err := player.AvailableEngines()

	if err != nil {
		slog.Error("could not get engines", "error", err)
		http.Error(w, "could not get engines", http.StatusInternalServerError)
		return
	}

	err = www.RenderIndex(w, www.Data{
		Board: h.renderBoardData(),
		Controls: www.ControlData{
			Engines: engines,
		},
	})

	if err != nil {
		slog.Error("error rendering template", "error", err)
		http.Error(w, "could not render template", http.StatusInternalServerError)
	}
}

func (h *Handler) broadcastBoard() {
	var buf bytes.Buffer

	err := www.RenderBoard(&buf, h.renderBoardData())

	if err != nil {
		panic(fmt.Sprintf("error rendering board: %v", err))
	}

	ev := sse.Event{
		ID:    fmt.Sprintf("id-%d", h.eventid),
		Event: EventMarkup,
		Data: Markup{
			Selector: "#board",
			Markup:   buf.String(),
		},
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

		h.broadcastBoard()
		h.lock.Unlock()
	}
}
