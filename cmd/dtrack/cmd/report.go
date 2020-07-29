package cmd

import (
	"log"
	"os"

	"github.com/nscuro/dependency-track-client/internal/report"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "generate reports",
	Run:   runReportCmd,
}

func init() {
	reportCmd.Flags().StringP("template", "t", "", "template path")
	reportCmd.Flags().StringP("output", "o", "", "output path")

	reportCmd.MarkFlagRequired("template")
	reportCmd.MarkFlagRequired("output")

	reportCmd.MarkFlagFilename("template", "html", "tpl")

	rootCmd.AddCommand(reportCmd)
}

func runReportCmd(cmd *cobra.Command, _ []string) {
	templatePath, _ := cmd.Flags().GetString("template")
	outputPath, _ := cmd.Flags().GetString("output")

	reportGenerator := report.NewGenerator(dtrackClient)

	project, err := dtrackClient.ResolveProject(pProjectUUID, pProjectName, pProjectVersion)
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
