package uci

import (
	"bufio"
	"context"
	"io"
	"os/exec"
)

func NewStockfish() (*Client, error) {
	return NewClient("stockfish")
}

func NewClient(name string, arg ...string) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, name, arg...)

	var stdout io.ReadCloser
	var stdin io.WriteCloser

	stdout, cmd.Stdout = io.Pipe()
	cmd.Stdin, stdin = io.Pipe()

	err := cmd.Start()

	if err != nil {
		cancel()
		return nil, err
	}

	c := &Client{
		close: func() {
			cancel()
			stdin.Close()
		},
		scan: bufio.NewScanner(stdout),
		w:    stdin,
	}

	go func() {
		defer stdout.Close()
		defer stdin.Close()
		defer cancel()
		cmd.Wait()
	}()

	err = c.initUCI()

	if err != nil {
		cancel()
		stdin.Close()
		return nil, err
	}

	return c, nil
}
