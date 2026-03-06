package main

import (
	"log/slog"
	"net/http"

	"github.com/rswilli/chess/internal/api"
)

func main() {
	h := api.NewHandler()

	serv := http.Server{
		Addr:    ":3000",
		Handler: h,
	}

	slog.Info("starting server", "addr", "localhost:3000")
	err := serv.ListenAndServe()

	slog.Error("server ended", "error", err)
}
