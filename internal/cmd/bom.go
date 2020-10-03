package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/nscuro/dependency-track-client"
)

var (
	bomCmd = &cobra.Command{
		Use:   "bom",
		Short: "Export and Upload BOMs",
	}

	bomExportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export a BOM",
		Run:   runBomExportCmd,
	}
	bomExportOpts bomExportOptions

	bomUploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload a BOM",
		Run:   runBomUploadCmd,
	}
	bomUploadOpts bomUploadOptions
)

type bomExportOptions struct {
	outputFilePath string
}

type bomUploadOptions struct {
	bomFilePath string
	autoCreate  bool
	wait        bool
}

func init() {
	bomExportCmd.Flags().StringVarP(&bomExportOpts.outputFilePath, "output", "o", "", "Output file path")
	bomCmd.AddCommand(bomExportCmd)

	bomUploadCmd.Flags().StringVar(&bomUploadOpts.bomFilePath, "bom", "", "BOM to upload")
	bomUploadCmd.Flags().BoolVar(&bomUploadOpts.autoCreate, "autocreate", false, "Automatically create project")
	bomUploadCmd.Flags().BoolVar(&bomUploadOpts.wait, "wait", false, "Wait for BOM processing to complete")
	bomCmd.AddCommand(bomUploadCmd)

	rootCmd.AddCommand(bomCmd)
}

func runBomExportCmd(_ *cobra.Command, _ []string) {
	dtrackClient := mustGetDTrackClient()

	projectUUID := globalOpts.projectUUID
	if projectUUID == "" {
		projectUUID = mustResolveProject(dtrackClient).UUID
	}

	bomXML, err := dtrackClient.ExportProjectAsCycloneDX(projectUUID)
	if err != nil {
		log.Fatalf("failed to export BOM: %v", err)
	}

	if bomExportOpts.outputFilePath == "" || bomExportOpts.outputFilePath == "-" {
		fmt.Println(bomXML)
		return
	}

	if err = ioutil.WriteFile(bomExportOpts.outputFilePath, []byte(bomXML), 0644); err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}
}

func runBomUploadCmd(_ *cobra.Command, _ []string) {
	var bomContent []byte
	var err error

	if bomUploadOpts.bomFilePath == "-" {
		bomContent, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read BOM from stdin: %v", err)
		}
	} else {
		bomContent, err = ioutil.ReadFile(bomUploadOpts.bomFilePath)
		if err != nil {
			log.Fatalf("failed to read BOM file: %v", err)
		}
	}

	token, err := mustGetDTrackClient().UploadBOM(dtrack.BOMSubmitRequest{
		ProjectUUID:    globalOpts.projectUUID,
		ProjectName:    globalOpts.projectName,
		ProjectVersion: globalOpts.projectVersion,
		AutoCreate:     bomUploadOpts.autoCreate,
		BOM:            base64.StdEncoding.EncodeToString(bomContent),
	})
	if err != nil {
		log.Fatalf("failed to upload BOM: %v", err)
	}

	fmt.Println(token)
}
