package main

import (
	"flag"
	"log/slog"
	"net/http"

	"github.com/rswilli/chess/internal/api"
)

var addr = flag.String("addr", ":3000", "listen address of the server")

func main() {
	flag.Parse()

	h := api.NewHandler()

	serv := http.Server{
		Addr:    *addr,
		Handler: h,
	}

	slog.Info("starting server", "addr", *addr)
	err := serv.ListenAndServe()

	slog.Error("server ended", "error", err)
}
