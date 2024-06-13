package wasm

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WasmEngine struct {
	runtime wazero.Runtime
	modules map[string]wazero.CompiledModule
}

func NewWasmEngine() *WasmEngine {
	runtime := wazero.NewRuntime(context.Background())
	wasi_snapshot_preview1.MustInstantiate(context.Background(), runtime)

	return &WasmEngine{
		runtime: runtime,
		modules: make(map[string]wazero.CompiledModule),
	}
}

func (e *WasmEngine) Load(ctx context.Context, name string, data []byte) error {
	module, err := e.runtime.CompileModule(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to compile module %q: %w", name, err)
	}
	e.modules[name] = module
	return nil
}

func (e *WasmEngine) Run(ctx context.Context, module string, body io.ReadCloser) (string, error) {
	compiled, ok := e.modules[module]
	if !ok {
		return "", fmt.Errorf("module %q not found", module)
	}

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	done := make(chan struct{}, 1)

	go func() {
		defer stdinWriter.Close()

		_, err := io.Copy(stdinWriter, body)
		if err != nil {
			log.Error().Err(err).Msg("failed to copy body to stdin")
		}
	}()

	output := bytes.Buffer{}
	out := bufio.NewWriter(&output)

	go func() {
		defer stdoutReader.Close()
		done <- struct{}{}

		_, err := io.Copy(out, stdoutReader)
		if err != nil {
			log.Error().Err(err).Msg("failed to copy stdout to discard")
		}
	}()

	cfg := wazero.NewModuleConfig().
		WithName(module).
		WithStdin(stdinReader).
		WithStdout(stdoutWriter)

	_, err := e.runtime.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to instantiate module %q: %w", module, err)
	}

	<-done

	return output.String(), nil
}
