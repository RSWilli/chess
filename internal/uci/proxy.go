package uci

import (
	"bufio"
	"fmt"
	"io"
)

func NewProxy(eng Engine) *Client {
	clientIn, serverOut := io.Pipe()
	serverIn, clientOut := io.Pipe()

	uciServer := NewServer(serverIn, serverOut, eng)

	close := func() {
		serverOut.Close()
		clientOut.Close()
		serverIn.Close()
		clientIn.Close()
	}

	go func() {
		defer close()
		err := uciServer.Run()

		fmt.Printf("error from uci server: %v\n", err)
	}()

	return &Client{
		close: close,
		scan:  bufio.NewScanner(clientIn),
		w:     clientOut,
	}
}
