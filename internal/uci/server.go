package uci

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rswilli/chess/internal/chess"
	"github.com/rswilli/chess/internal/uci/search"
)

// Engine is a chess engine that responds to the UCI commands.
// Any returned error will make [Server.Run] return and thus stop the server
type Engine interface {
	NewGame() error
	Position(fen string, moves []string) error
	Perft(depth int) (total int, moves map[string]int, err error)
	Ready() error

	// Go starts the search with the given options and returns the best move.
	//
	// FIXME: stockfish outputs a lot of "info" messages here that get dropped by this
	// interface
	Go(options search.Options) (bestmove string, ponder string)
	Stop()
}

// Server exposes an [Engine] to the given streams, typically [os.Stdout] and [os.Stdin]
type Server struct {
	r io.Reader
	w io.Writer

	h Engine
}

func NewServer(r io.Reader, w io.Writer, h Engine) *Server {
	return &Server{
		r: r,
		w: w,
		h: h,
	}
}

func (s *Server) respond(cmd string) error {
	_, err := s.w.Write([]byte(cmd + "\n"))

	return err
}

func (s *Server) Run() error {
	scan := bufio.NewScanner(s.r)

	for scan.Scan() {
		err := s.handleCommand(scan.Text())

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) handleCommand(line string) error {
	if line == "uci" {
		return s.respond("uciok")
	}
	if line == "ucinewgame" {
		s.h.NewGame()
		return nil
	}
	if line == "quit" {
		return fmt.Errorf("server closed by user")
	}
	if line == "stop" {
		s.h.Stop()
		return nil
	}
	if line == "isready" {
		err := s.h.Ready()

		if err != nil {
			return err
		}

		return s.respond("readyok")
	}

	if strings.HasPrefix(line, "position") {
		return s.handlePositionCommand(line)
	}

	if strings.HasPrefix(line, "go") {
		return s.handleGoCommand(line)
	}

	return fmt.Errorf("unexpected command received: %s", line)
}

func (s *Server) handlePositionCommand(line string) error {
	// cases:
	// position startpos
	// position startpos moves
	// position startpos moves a1a3 a4g6
	// position fen rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
	// position fen rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 moves
	// position fen rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 moves a1a3 a4g6

	parts := strings.Split(line, " moves")

	if len(parts) > 2 {
		return fmt.Errorf("malformed position command")
	}

	var moves []string

	if len(parts) > 1 {
		moves = strings.Fields(parts[1])
	}

	var fen string

	if cut, ok := strings.CutPrefix(parts[0], "position fen "); ok {
		fen = cut
	} else if cut, ok := strings.CutPrefix(parts[0], "position "); ok {
		fen = cut
	} else {
		return fmt.Errorf("malformed position command")
	}

	if fen == StartPosition {
		fen = chess.DefaultFen
	}

	return s.h.Position(fen, moves)
}

func (s *Server) handleGoCommand(line string) error {
	parts := strings.Split(line, " ")

	if parts[0] != "go" {
		panic("must be called with a go command")
	}

	parts = parts[1:]

	opts := search.Options{}

	for len(parts) > 0 {
		switch parts[0] {
		case "ponder":
			opts.Ponder = true
			parts = parts[1:]
			continue
		case "infinite":
			opts.Infinite = true
			parts = parts[1:]
			continue
		case "wtime":
			if len(parts) == 1 {
				return fmt.Errorf("missing argument")
			}

			ms, err := parseInt(parts[1])

			if err != nil {
				return err
			}

			opts.WhiteTime = time.Duration(ms) * time.Millisecond

			parts = parts[2:]
			continue
		case "btime":
			if len(parts) == 1 {
				return fmt.Errorf("missing argument")
			}

			ms, err := parseInt(parts[1])

			if err != nil {
				return err
			}

			opts.BlackTime = time.Duration(ms) * time.Millisecond

			parts = parts[2:]
			continue
		case "winc":
			if len(parts) == 1 {
				return fmt.Errorf("missing argument")
			}

			ms, err := parseInt(parts[1])

			if err != nil {
				return err
			}

			opts.WhiteIncrement = time.Duration(ms) * time.Millisecond

			parts = parts[2:]
			continue
		case "binc":
			if len(parts) == 1 {
				return fmt.Errorf("missing argument")
			}

			ms, err := parseInt(parts[1])

			if err != nil {
				return err
			}

			opts.BlackIncrement = time.Duration(ms) * time.Millisecond

			parts = parts[2:]
			continue
		case "movestogo":
			n, err := parseInt(parts[1])

			if err != nil {
				return err
			}

			opts.MovesToGo = n

			parts = parts[2:]
			continue
		case "depth":
			n, err := parseInt(parts[1])

			if err != nil {
				return err
			}

			opts.Depth = n

			parts = parts[2:]
			continue
		case "nodes":
			n, err := parseInt(parts[1])

			if err != nil {
				return err
			}

			opts.Nodes = n

			parts = parts[2:]
			continue
		case "mate":
			n, err := parseInt(parts[1])

			if err != nil {
				return err
			}

			opts.Mate = n

			parts = parts[2:]
			continue
		case "perft":
			n, err := parseInt(parts[1])

			if err != nil {
				return err
			}

			// Special case: when any argument is perft call the perft instead
			var res perftResult

			res.total, res.moves, err = s.h.Perft(n)

			if err != nil {
				return err
			}

			fmt.Fprintf(s.w, "%s", res.String())
			return nil
		case "movetime":
			if len(parts) == 1 {
				return fmt.Errorf("missing argument")
			}

			ms, err := parseInt(parts[1])

			if err != nil {
				return err
			}

			opts.MoveTime = time.Duration(ms) * time.Millisecond

			parts = parts[2:]
			continue
		}

		return fmt.Errorf("unknown go command argument: %s", line)
	}

	go func() {
		var bm bestmove
		bm.bestMove, bm.ponder = s.h.Go(opts)

		s.respond(bm.String())
	}()

	return nil
}
