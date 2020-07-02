package cmd

import (
	"html/template"
	"log"
	"os"

	"github.com/nscuro/dependency-track-client/pkg/dtrack"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a vulnerability report",
	Run:   runReportCmd,
}

func init() {
	reportCmd.Flags().StringP("project", "p", "", "project name")
	reportCmd.Flags().StringP("version", "v", "", "project version")
	reportCmd.Flags().String("uuid", "", "project uuid")
	reportCmd.Flags().StringP("template", "t", "", "template path")
	reportCmd.Flags().StringP("output", "o", "", "output path")

	rootCmd.AddCommand(reportCmd)
}

func runReportCmd(cmd *cobra.Command, _ []string) {
	projectUUID, _ := cmd.Flags().GetString("uuid")
	templatePath, _ := cmd.Flags().GetString("template")
	outputPath, _ := cmd.Flags().GetString("output")

	dtrackClient := dtrack.NewClient(pBaseURL, pAPIKey)

	findings, err := dtrackClient.GetFindings(projectUUID)
	if err != nil {
		log.Fatal("retrieving findings failed: ", err)
		return
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal("failed to parse template: ", err)
		return
	}

	outputFile, err := os.Open(outputPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer outputFile.Close()

	if err = tmpl.Execute(outputFile, findings); err != nil {
		log.Fatal("failed to execute template: ", err)
		return
	}
}
