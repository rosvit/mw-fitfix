//go:build js

package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"syscall/js"

	"github.com/rosvit/mw-fitfix/internal/fitfix"
)

func main() {
	shutdown := make(chan struct{})
	js.Global().Set("fixFit", js.FuncOf(fixFIT))
	<-shutdown
}

func fixFIT(jsThis js.Value, args []js.Value) any {
	fileArg := args[0]
	optsArg := args[1].String()
	fitContent := make([]byte, fileArg.Get("byteLength").Int())
	js.CopyBytesToGo(fitContent, fileArg)
	opts := fitfix.Options{}
	err := json.Unmarshal([]byte(optsArg), &opts)
	if err != nil {
		slog.Warn("Can't unmarshal options, using defaults:", fitfix.KeyError, err)
	}

	inBuf := bytes.NewBuffer(fitContent)
	var outBuf bytes.Buffer
	if err = fitfix.Process(inBuf, &outBuf, opts); err != nil {
		slog.Error("Failed to fix FIT file:", fitfix.KeyError, err)
		return nil
	}
	output := outBuf.Bytes()
	result := js.Global().Get("Uint8Array").New(len(output))
	js.CopyBytesToJS(result, output)
	return result
}
