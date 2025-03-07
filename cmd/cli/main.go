//go:build !js

package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/rosvit/mw-fitfix/internal/fitfix"
)

const extension = ".fit"

func main() {
	laps := flag.Bool("laps", true, "update laps to single lap for whole activity")
	device := flag.Bool("device", true, "remove default Garmin Edge device")
	replace := false
	target := ""
	flag.BoolVar(&replace, "replace", false, "replace original FIT file with fixed one")
	flag.StringVar(&target, "target", "", "target FIT file")
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		log.Fatalf("Invalid number of arguments: got %d, expected 1", len(args))
	}
	src := args[0]

	if target == "" {
		if replace {
			target = src
		} else if orig, found := strings.CutSuffix(src, extension); found {
			target = fmt.Sprintf("%s_fixed%s", orig, extension)
		} else {
			log.Fatal("Target file is required")
		}
	}

	err := fitfix.ProcessAndWrite(src, target, *laps, *device)
	if err != nil {
		log.Fatal("Failed to fix MyWhoosh FIT file")
	}

	fmt.Printf("Fixed MyWhoosh FIT file written: %s\n", target)
}
