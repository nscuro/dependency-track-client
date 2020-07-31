package cmd

import (
	"log"
	"os"

	"github.com/nscuro/dependency-track-client/internal/report"
	"github.com/nscuro/dependency-track-client/pkg/dtrack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate reports",
	Run:   runReportCmd,
}

func init() {
	reportCmd.Flags().StringP("template", "t", "", "Template path")
	reportCmd.Flags().StringP("output", "o", "", "Output path")

	reportCmd.MarkFlagRequired("template")
	reportCmd.MarkFlagRequired("output")

	reportCmd.MarkFlagFilename("template", "html", "tpl")

	rootCmd.AddCommand(reportCmd)
}

func runReportCmd(cmd *cobra.Command, _ []string) {
	dtrackClient := dtrack.NewClient(viper.GetString("url"), viper.GetString("api-key"))

	templatePath, _ := cmd.Flags().GetString("template")
	outputPath, _ := cmd.Flags().GetString("output")

	reportGenerator := report.NewGenerator(dtrackClient)

	project, err := dtrackClient.ResolveProject(projectUUID, projectName, projectVersion)
	if err != nil {
		log.Fatal("failed to resolve project: ", err)
		return
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer outputFile.Close()

	if err = reportGenerator.GenerateProjectReport(project, templatePath, outputFile); err != nil {
		log.Fatal("generating report failed: ", err)
	}
}
