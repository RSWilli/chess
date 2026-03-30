package uci

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rswilli/chess/internal/uci/search"
)

// Client talks to an engine via two streams in the uci protocol
type Client struct {
	close func()
	scan  *bufio.Scanner
	w     io.WriteCloser
}

var _ Engine = (*Client)(nil)

func (c *Client) Close() {
	c.close()
}

const StartPosition = "startpos"

func (c *Client) read() (string, error) {
	ok := c.scan.Scan()

	if !ok {
		return "", fmt.Errorf("could not scan")
	}

	line := c.scan.Text()

	// fmt.Println("read " + line)

	return line, nil
}

func (c *Client) writeCommand(cmd string) error {
	// fmt.Println("writing " + cmd)

	l := append([]byte(cmd), '\n')

	_, err := c.w.Write(l)

	return err
}

func (c *Client) skipUntil(returnCmd string) error {
	for {
		line, err := c.read()

		if err != nil {
			return err
		}

		if line == returnCmd {
			return nil
		}
	}
}

// readUntil reads until the given line appears, returning all but the given line in the list
func (c *Client) readUntil(targetline string) ([]string, error) {
	var lines []string
	for {
		line, err := c.read()

		if err != nil {
			return nil, err
		}

		if line == targetline {
			return lines, nil
		}

		lines = append(lines, line)
	}
}

// initUCI calls the uci commmand and waits for uciok
func (c *Client) initUCI() error {
	err := c.writeCommand("uci")

	if err != nil {
		return err
	}

	// TODO: uci responds with all options first

	return c.skipUntil("uciok")
}

// waitForReady calls the isready commmand and waits for readyok
func (c *Client) waitForReady() error {
	err := c.writeCommand("isready")

	if err != nil {
		return err
	}

	return c.skipUntil("readyok")
}

func (c *Client) NewGame() error {
	err := c.writeCommand("ucinewgame")

	if err != nil {
		return err
	}

	// https://official-stockfish.github.io/docs/stockfish-wiki/UCI-&-Commands.html#ucinewgame
	// > GUI should always send isready after ucinewgame

	return c.waitForReady()
}

func (c *Client) Perft(depth int) (total int, moves map[string]int, err error) {
	err = c.writeCommand(fmt.Sprintf("go perft %d", depth))

	if err != nil {
		return 0, nil, err
	}

	// first empty line signifies the end of the moves
	response, err := c.readUntil("")

	if err != nil {
		return 0, nil, err
	}

	moves = make(map[string]int)

	for _, l := range response {
		if strings.HasPrefix(l, "info") {
			// before our desired output stockfish outputs info strings
			// we skip these for perft
			continue
		}

		// e.g. "g2g4: 214048"

		parts := strings.Split(l, ": ")

		if len(parts) != 2 {
			return 0, nil, fmt.Errorf("move wrong format: %s", l)
		}

		count, err := parseInt(parts[1])

		if err != nil {
			return 0, nil, err
		}

		moves[parts[0]] = count
	}

	totalLine, err := c.read()

	if err != nil {
		return 0, nil, err
	}

	totalNum := strings.TrimPrefix(totalLine, "Nodes searched: ")

	total, err = parseInt(totalNum)

	if err != nil {
		return 0, nil, err
	}

	// ends with a new line
	_, err = c.read()

	if err != nil {
		return 0, nil, err
	}

	return total, moves, nil
}

func (c *Client) Go(opts search.Options) (bestmove string, ponder string) {
	var parts []string

	parts = append(parts, "go")

	if len(opts.SearchMoves) > 0 {
		parts = append(parts, "searchmoves")
		parts = append(parts, strings.Join(opts.SearchMoves, " "))
	}
	if opts.Ponder {
		parts = append(parts, "ponder")
	}
	if opts.WhiteTime != 0 {
		parts = append(parts, fmt.Sprintf("wtime %d", opts.WhiteTime/time.Millisecond))
	}
	if opts.BlackTime != 0 {
		parts = append(parts, fmt.Sprintf("btime %d", opts.BlackTime/time.Millisecond))
	}
	if opts.WhiteIncrement != 0 {
		parts = append(parts, fmt.Sprintf("winc %d", opts.WhiteIncrement/time.Millisecond))
	}
	if opts.BlackIncrement != 0 {
		parts = append(parts, fmt.Sprintf("binc %d", opts.BlackIncrement/time.Millisecond))
	}
	if opts.MovesToGo != 0 {
		parts = append(parts, fmt.Sprintf("movestogo %d", opts.MovesToGo))
	}
	if opts.Depth != 0 {
		parts = append(parts, fmt.Sprintf("depth %d", opts.Depth))
	}
	if opts.Nodes != 0 {
		parts = append(parts, fmt.Sprintf("nodes %d", opts.Nodes))
	}
	if opts.Mate != 0 {
		parts = append(parts, fmt.Sprintf("mate %d", opts.Mate))
	}
	if opts.MoveTime != 0 {
		parts = append(parts, fmt.Sprintf("movetime %d", opts.MoveTime/time.Millisecond))
	}
	if opts.Infinite {
		parts = append(parts, "infinite")
	}

	err := c.writeCommand(strings.Join(parts, " "))

	if err != nil {
		return "", ""
	}

	for {
		line, err := c.read()

		if err != nil {
			return "", ""
		}

		if strings.HasPrefix(line, "info") {
			continue // drop the info lines
		}

		if !strings.HasPrefix(line, "bestmove") {
			return "", "" // unexpected command
		}

		parts := strings.Split(line, " ")

		if len(parts) == 2 {
			// no ponder move, e.g. checkmate
			return parts[1], ""
		}

		if len(parts) != 4 || parts[2] != "ponder" {
			return "", "" // invalid format
		}

		return parts[1], parts[3]
	}
}

// Position calls the position command to intialize a position. as a special case fen can also be "startpos"
func (c *Client) Position(fen string, moves []string) error {
	var buf bytes.Buffer
	buf.WriteString("position ")

	if fen != StartPosition {
		buf.WriteString("fen ")
	}

	buf.WriteString(fen)

	if len(moves) != 0 {
		buf.WriteString(" moves ")

		for _, m := range moves {
			buf.WriteString(m)
			buf.WriteString(" ")
		}
	}

	err := c.writeCommand(buf.String())

	if err != nil {
		return err
	}

	return nil
}

// Ready implements [Engine].
func (c *Client) Ready() error {
	return c.waitForReady()
}

// Stop implements [Engine].
func (c *Client) Stop() {
	c.close()
}

// func (c *Client) readInfo() (Info, error) {
// 	line, err := c.read()

// 	if err != nil {
// 		return Info{}, err
// 	}

// 	var info Info

// 	err = info.UnmarshalText([]byte(line))

// 	if err != nil {
// 		return Info{}, err
// 	}

// 	return info, nil
// }
