package main

import (
	"bytes"
	"fmt"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/mm"
	"github.com/rswilli/chess/internal/pkg/notify"
	"github.com/rswilli/chess/internal/sse"
	"github.com/rswilli/chess/internal/www"
)

type Handler struct {
	id int

	game *mm.Game
	b    *notify.Broker[sse.Event]
}

func (h Handler) BroadcastBoard(currentSquare chess.Square) {
	position := h.game.Position()

	moves := position.GenerateMoves()

	targets := make(map[chess.Square][]chess.Move)

	promotion := false

	for _, m := range moves {
		if m.From == currentSquare {
			promotion = m.Special.Has(chess.PromoteAny)
			targets[m.To] = append(targets[m.To], m)
		}
	}

	var buf bytes.Buffer

	err := www.RenderBoard(&buf, www.Data{
		Board:       position,
		Selected:    currentSquare,
		MoveTargets: targets,
		Promotion:   promotion,
	})

	if err != nil {
		panic(fmt.Sprintf("error rendering board: %v", err))
	}

	ev := sse.Event{
		ID:    fmt.Sprintf("id-%d", h.id),
		Event: EventMarkup,
		Data: Markup{
			Selector: "#board",
			Markup:   buf.String(),
		},
	}

	h.id++

	h.b.Send(ev)
}
