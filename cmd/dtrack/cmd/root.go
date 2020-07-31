package cmd

import (
	"fmt"

	"github.com/nscuro/dependency-track-client/pkg/dtrack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:               "dtrack",
		PersistentPreRunE: preRunRootCmd,
	}

	projectUUID    string
	projectName    string
	projectVersion string

	dtrackClient *dtrack.Client
)

func init() {
	rootCmd.PersistentFlags().StringP("url", "u", "", "Dependency-Track URL")
	rootCmd.PersistentFlags().StringP("api-key", "k", "", "Dependency-Track API key")

	rootCmd.PersistentFlags().StringVar(&projectUUID, "project", "", "Project UUID")
	rootCmd.PersistentFlags().StringVar(&projectName, "project-name", "", "Project name")
	rootCmd.PersistentFlags().StringVar(&projectVersion, "project-version", "", "Project version")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("api-key"))

	viper.SetEnvPrefix("DTRACK")
	viper.BindEnv("url")
	viper.BindEnv("apikey")
}

func preRunRootCmd(_ *cobra.Command, _ []string) error {
	client, err := dtrack.NewClient(viper.GetString("url"), viper.GetString("apikey"))
	if err != nil {
		return fmt.Errorf("failed to initialize dependency-track client: %w", err)
	}

	dtrackClient = client
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
