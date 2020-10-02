package main

import (
	"os"

	"github.com/nscuro/dependency-track-client/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
