package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/nscuro/dependency-track-client/internal/report"
)

var (
	reportCmd = &cobra.Command{
		Use:   "report",
		Short: "Generate reports",
		Run:   runReportCmd,
	}
	reportOpts reportOptions
)

type reportOptions struct {
	templateFilePath string
	outputFilePath   string
}

func init() {
	reportCmd.Flags().StringVarP(&reportOpts.templateFilePath, "template", "t", "", "Template path")
	reportCmd.Flags().StringVarP(&reportOpts.outputFilePath, "output", "o", "", "Output path")

	rootCmd.AddCommand(reportCmd)
}

func runReportCmd(_ *cobra.Command, _ []string) {
	dtrackClient := mustGetDTrackClient()
	project := mustResolveProject(dtrackClient)

	outputFile, err := os.Create(reportOpts.outputFilePath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	reportGenerator := report.NewGenerator(dtrackClient)
	if err := reportGenerator.GenerateProjectReport(project, reportOpts.templateFilePath, outputFile); err != nil {
		log.Fatalf("failed to generate report: %v", err)
	}
}
