package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/nscuro/dependency-track-client/internal/version"
)

var (
	versionCmd = &cobra.Command{
		Use: "version",
		Run: runVersionCmd,
	}
	versionOpts versionOptions
)

type versionOptions struct {
	clientOnly bool
	serverOnly bool
}

func init() {
	versionCmd.Flags().BoolVar(&versionOpts.clientOnly, "client", false, "Show only client version")
	versionCmd.Flags().BoolVar(&versionOpts.serverOnly, "server", false, "Show only server version")

	rootCmd.AddCommand(versionCmd)
}

func runVersionCmd(_ *cobra.Command, _ []string) {
	showBoth := versionOpts.clientOnly && versionOpts.serverOnly

	if !showBoth && versionOpts.clientOnly {
		fmt.Println(version.Version)
		return
	}

	var serverVersion string
	if about, err := mustGetDTrackClient().GetAbout(); err == nil {
		serverVersion = about.Version
		fmt.Printf("Server: %s\n", about.Version)
	} else {
		log.Fatalf("failed to retrieve server version: %v", err)
	}

	if !showBoth && versionOpts.serverOnly {
		log.Println(serverVersion)
		return
	}

	log.Printf("Client: %s\nServer: %s\n", version.Version, serverVersion)
}
