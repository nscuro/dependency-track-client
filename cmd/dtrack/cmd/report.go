package cmd

import (
	"log"
	"os"

	"github.com/nscuro/dependency-track-client/internal/report"
	"github.com/spf13/cobra"
)

var (
	reportCmd = &cobra.Command{
		Use:   "report",
		Short: "Generate reports",
		Run:   runReportCmd,
	}

	reportOpts ReportOptions
)

type ReportOptions struct {
	TemplatePath string
	OutputPath   string
}

func init() {
	reportCmd.Flags().StringVarP(&reportOpts.TemplatePath, "template", "t", "", "Template path")
	reportCmd.Flags().StringVarP(&reportOpts.OutputPath, "output", "o", "", "Output path")

	reportCmd.MarkFlagRequired("template")
	reportCmd.MarkFlagRequired("output")

	reportCmd.MarkFlagFilename("template", "html", "tpl")

	rootCmd.AddCommand(reportCmd)
}

func runReportCmd(cmd *cobra.Command, _ []string) {
	reportGenerator := report.NewGenerator(dtrackClient)

	project, err := dtrackClient.ResolveProject(projectUUID, projectName, projectVersion)
	if err != nil {
		log.Fatal("failed to resolve project: ", err)
		return
	}

	outputFile, err := os.Create(reportOpts.OutputPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer outputFile.Close()

	if err = reportGenerator.GenerateProjectReport(project, reportOpts.TemplatePath, outputFile); err != nil {
		log.Fatal("generating report failed: ", err)
	}
}
