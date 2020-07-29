package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use: "dtrack",
	}

	pProjectUUID    string
	pProjectName    string
	pProjectVersion string
)

func init() {
	rootCmd.PersistentFlags().StringP("url", "u", "", "dependency-track base url")
	rootCmd.PersistentFlags().StringP("api-key", "k", "", "dependency-track api key")

	rootCmd.PersistentFlags().StringVar(&pProjectUUID, "project-uuid", "", "project uuid")
	rootCmd.PersistentFlags().StringVar(&pProjectName, "project-name", "", "project name")
	rootCmd.PersistentFlags().StringVar(&pProjectVersion, "project-version", "", "project version")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))

	viper.SetEnvPrefix("DTRACK")
	viper.BindEnv("url", "URL")
	viper.BindEnv("api-key", "API_KEY")
}

func Execute() error {
	return rootCmd.Execute()
}
