package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use: "dtrack",
	}

	pBaseURL string
	pAPIKey  string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&pBaseURL, "base-url", "u", "", "dependency-track base url")
	rootCmd.PersistentFlags().StringVarP(&pAPIKey, "api-key", "k", "", "dependency-track api key")

	rootCmd.MarkPersistentFlagRequired("base-url")
	rootCmd.MarkPersistentFlagRequired("api-key")

	viper.BindPFlag("base-url", rootCmd.PersistentFlags().Lookup("base-url"))
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))

	viper.SetEnvPrefix("DTRACK")
	viper.BindEnv("base-url", "BASE_URL")
	viper.BindEnv("api-key", "API_KEY")
}

func Execute() error {
	return rootCmd.Execute()
}
