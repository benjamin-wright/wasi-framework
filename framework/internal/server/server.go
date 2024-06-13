package server

import (
	"io"
	"net/http"

	"github.com/benjamin-wright/wasi-framework/framework/internal/wasm"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Start(engine *wasm.WasmEngine) func() {
	r := gin.Default()

	r.POST("/wasm/:module", func(c *gin.Context) {
		module := c.Param("module")
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Info().Str("module", module).Int("bytes", len(data)).Msg("loading module")

		err = engine.Load(c.Request.Context(), module, data)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": "ok"})
	})

	r.PUT("/wasm/:module", func(c *gin.Context) {
		module := c.Param("module")

		result, err := engine.Run(c.Request.Context(), module, c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"result": result})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("server stopped with error")
		}
	}()

	return func() {
		err := srv.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close server")
		}
	}
}
