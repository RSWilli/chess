package uci

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Client struct {
	close func()
	scan  *bufio.Scanner
	w     io.WriteCloser
}

func (c *Client) Close() {
	c.close()
}

const StartPositionFEN = "startpos"

func (c *Client) read() (string, error) {
	ok := c.scan.Scan()

	if !ok {
		return "", fmt.Errorf("could not scan")
	}

	line := c.scan.Text()

	fmt.Println("read " + line)

	return line, nil
}

func (c *Client) writeCommand(cmd string) error {
	fmt.Println("writing " + cmd)

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
	return c.writeCommand("ucinewgame")
}

func (c *Client) Perft(depth int) (int, map[string]int, error) {
	err := c.writeCommand(fmt.Sprintf("go perft %d", depth))

	if err != nil {
		return 0, nil, err
	}

	for range 4 { // FIXME: 4 is stockfish specific
		_, err := c.readInfo()

		if err != nil {
			return 0, nil, err
		}
	}

	// first empty line signifies the end of the moves
	moveLines, err := c.readUntil("")

	if err != nil {
		return 0, nil, err
	}

	moves := make(map[string]int)

	for _, l := range moveLines {
		// e.g. "g2g4: 214048"

		parts := strings.Split(l, ": ")

		if len(parts) != 2 {
			return 0, nil, fmt.Errorf("move wrong format: %s", l)
		}

		count, err := strconv.ParseInt(parts[1], 10, 64)

		if err != nil {
			return 0, nil, err
		}

		moves[parts[0]] = int(count)
	}

	totalLine, err := c.read()

	if err != nil {
		return 0, nil, err
	}

	totalNum := strings.TrimPrefix(totalLine, "Nodes searched: ")

	total, err := strconv.ParseInt(totalNum, 10, 64)

	if err != nil {
		return 0, nil, err
	}

	// ends with a new line
	_, err = c.read()

	if err != nil {
		return 0, nil, err
	}

	return int(total), moves, nil
}

// Position calls the position command to intialize a position. as a special case fen can also be "startpos"
func (c *Client) Position(fen string, moves []string) error {
	var buf bytes.Buffer
	buf.WriteString("position ")
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

	return c.waitForReady()
}

func (c *Client) readInfo() (Info, error) {
	line, err := c.read()

	if err != nil {
		return Info{}, err
	}

	var info Info

	err = info.UnmarshalText([]byte(line))

	if err != nil {
		return Info{}, err
	}

	return info, nil
}
