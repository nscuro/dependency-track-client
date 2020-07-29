package cmd

import (
	"github.com/nscuro/dependency-track-client/pkg/dtrack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use: "dtrack",
	}

	pBaseURL string
	pAPIKey  string

	pProjectUUID    string
	pProjectName    string
	pProjectVersion string

	dtrackClient *dtrack.Client
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&pBaseURL, "url", "u", "", "dependency-track base url")
	rootCmd.PersistentFlags().StringVarP(&pAPIKey, "api-key", "k", "", "dependency-track api key")

	rootCmd.PersistentFlags().StringVar(&pProjectUUID, "project-uuid", "", "project uuid")
	rootCmd.PersistentFlags().StringVar(&pProjectName, "project-name", "", "project name")
	rootCmd.PersistentFlags().StringVar(&pProjectVersion, "project-version", "", "project version")

	rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.MarkPersistentFlagRequired("api-key")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))

	viper.SetEnvPrefix("DTRACK")
	viper.BindEnv("url", "URL")
	viper.BindEnv("api-key", "API_KEY")

	dtrackClient = dtrack.NewClient(viper.GetString("url"), viper.GetString("api-key"))
}

func Execute() error {
	return rootCmd.Execute()
}
