package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/benjamin-wright/wasi-framework/framework/internal/server"
	"github.com/benjamin-wright/wasi-framework/framework/internal/wasm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("starting server")

	engine := wasm.NewWasmEngine()
	closer := server.Start(engine)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	closer()
}
