package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nscuro/dependency-track-client/pkg/dtrack"
)

var (
	rootCmd = &cobra.Command{
		Use: "dtrack",
	}
	globalOpts globalOptions
)

type globalOptions struct {
	projectUUID    string
	projectName    string
	projectVersion string
}

func init() {
	viper.SetDefault("dtrack-url", "")
	viper.SetDefault("dtrack-apikey", "")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	rootCmd.PersistentFlags().StringP("url", "u", "", "Dependency-Track URL")
	rootCmd.PersistentFlags().StringP("apikey", "k", "", "Dependency-Track API Key")

	rootCmd.PersistentFlags().StringVar(&globalOpts.projectUUID, "project", "", "Project UUID")
	rootCmd.PersistentFlags().StringVar(&globalOpts.projectName, "project-name", "", "Project Name")
	rootCmd.PersistentFlags().StringVar(&globalOpts.projectVersion, "project-version", "", "Project Version")

	viper.BindPFlag("dtrack-url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("dtrack-apikey", rootCmd.PersistentFlags().Lookup("apikey"))
}

func Execute() error {
	return rootCmd.Execute()
}

func mustGetDTrackClient() *dtrack.Client {
	client, err := dtrack.NewClient(viper.GetString("dtrack-url"), viper.GetString("dtrack-apikey"))
	if err != nil {
		log.Fatalf("failed to initialize dtrack client: %v", err)
	}
	return client
}

func mustResolveProject(dtrackClient *dtrack.Client) *dtrack.Project {
	project, err := dtrackClient.ResolveProject(globalOpts.projectUUID, globalOpts.projectName, globalOpts.projectVersion)
	if err != nil {
		log.Fatalf("failed to resolve project: %v", err)
	}
	return project
}
