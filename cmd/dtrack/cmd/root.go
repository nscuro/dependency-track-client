package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use: "dtrack",
	}

	projectUUID    string
	projectName    string
	projectVersion string
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

func Execute() error {
	return rootCmd.Execute()
}
