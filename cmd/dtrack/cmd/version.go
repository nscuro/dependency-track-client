package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE:  runVersionCmd,
	}

	versionOpts VersionOptions
)

type VersionOptions struct {
	OnlyClient bool
	OnlyServer bool
}

func init() {
	versionCmd.Flags().BoolVar(&versionOpts.OnlyClient, "client", false, "Show only client version")
	versionCmd.Flags().BoolVar(&versionOpts.OnlyServer, "server", false, "Show only server version")

	rootCmd.AddCommand(versionCmd)
}

func runVersionCmd(_ *cobra.Command, _ []string) error {
	showBoth := versionOpts.OnlyClient && versionOpts.OnlyServer

	if showBoth || !versionOpts.OnlyServer {
		fmt.Printf("Client: %s\n", "0.1.0")
	}

	if showBoth || !versionOpts.OnlyClient {
		about, err := dtrackClient.GetAbout()
		if err != nil {
			return fmt.Errorf("failed to retrieve server version: %w", err)
		}
		fmt.Printf("Server: %s %s (%s %s)\n", about.Application, about.Version, about.Framework.Name, about.Framework.Version)
	}

	return nil
}
