package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nscuro/dependency-track-client/pkg/dtrack"
	"github.com/spf13/cobra"
)

var (
	bomCmd = &cobra.Command{
		Use:   "bom",
		Short: "Upload and export BOMs",
	}

	bomUploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload a BOM",
		Run:   runBomUploadCmd,
	}

	bomExportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export a BOM",
		Run:   runBomExportCmd,
	}
)

func init() {
	initBomUploadCmd()
	initBomGetCmd()

	bomCmd.AddCommand(bomUploadCmd)
	bomCmd.AddCommand(bomExportCmd)

	rootCmd.AddCommand(bomCmd)
}

func initBomUploadCmd() {
	bomUploadCmd.Flags().StringP("bom", "b", "", "BOM path")
	bomUploadCmd.Flags().Bool("autocreate", false, "Automatically create project")

	bomUploadCmd.MarkFlagRequired("bom")
	bomUploadCmd.MarkFlagFilename("bom", "xml", "json")
}

func initBomGetCmd() {
	bomExportCmd.Flags().StringP("output", "o", "", "")
}

func runBomUploadCmd(cmd *cobra.Command, _ []string) {
	bomPath, _ := cmd.Flags().GetString("bom")
	autoCreate, _ := cmd.Flags().GetBool("autocreate")

	log.Println("reading bom")
	bomContent, err := ioutil.ReadFile(bomPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("uploading bom")
	_, err = dtrackClient.UploadBOM(dtrack.BOMSubmitRequest{
		ProjectUUID:    projectUUID,
		ProjectName:    projectName,
		ProjectVersion: projectVersion,
		AutoCreate:     autoCreate,
		BOM:            base64.StdEncoding.EncodeToString(bomContent),
	})
	if err != nil {
		log.Fatal("uploading bom failed: ", err)
		return
	}
	log.Println("bom was successfully uploaded")
}

func runBomExportCmd(cmd *cobra.Command, _ []string) {
	log.Println("resolving project")
	project, err := dtrackClient.ResolveProject(projectUUID, projectName, projectVersion)
	if err != nil {
		log.Fatal("failed to resolve project: ", err)
		return
	}

	log.Println("retrieving bom")
	bom, err := dtrackClient.ExportProjectAsCycloneDX(project.UUID)
	if err != nil {
		log.Fatal("retrieving bom failed: ", err)
		return
	}

	output, _ := cmd.Flags().GetString("output")
	if output == "" {
		fmt.Println(bom)
		return
	}

	log.Printf("writing bom to %s\n", output)
	if err = ioutil.WriteFile(output, []byte(bom), 0644); err != nil {
		log.Fatal("failed to write output file: ", err)
	}
}
