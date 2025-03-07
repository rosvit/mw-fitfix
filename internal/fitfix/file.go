//go:build !js

package fitfix

import (
	"log/slog"
	"os"
)

func ProcessAndWrite(src, target string, laps, device bool) error {
	closeFile := func(f *os.File) {
		if err := f.Close(); err != nil {
			slog.Error("Failed to close file:", KeyFile, f.Name(), KeyError, err)
		}
	}

	srcF, err := os.Open(src)
	if err != nil {
		slog.Error("Failed to open source file:", KeyFile, src, KeyError, err)
		return err
	}
	defer closeFile(srcF)

	targetF, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		slog.Error("Failed to open target file:", KeyFile, target, KeyError, err)
		return err
	}
	defer closeFile(targetF)

	return Process(srcF, targetF, Options{Device: device, Laps: laps})
}
