package main

import (
	"os"

	"github.com/nscuro/dependency-track-client/cmd/dtrack/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
