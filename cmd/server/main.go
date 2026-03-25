package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/rswilli/chess/internal/api"
	"github.com/rswilli/chess/internal/player"
)

var addr = flag.String("addr", ":3000", "listen address of the server")
var enginesFolder = flag.String("engines", "./engines", "the directory where to look for engines")

func main() {
	flag.Parse()

	reg, err := player.NewEnginesRegistry(*enginesFolder)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating engines registry: %v", err)
		flag.Usage()
		os.Exit(2)
	}

	h := api.NewHandler(reg)

	serv := http.Server{
		Addr:    *addr,
		Handler: h,
	}

	slog.Info("starting server", "addr", *addr)
	err = serv.ListenAndServe()

	slog.Error("server ended", "error", err)
}
