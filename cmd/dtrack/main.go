package main

import (
	"log"

	"github.com/nscuro/dependency-track-client/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
