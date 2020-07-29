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
		Use: "bom",
	}

	bomUploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "upload a bom",
		Run:   runBomUploadCmd,
	}

	bomGetCmd = &cobra.Command{
		Use:   "get",
		Short: "retrieve a bom",
		Run:   runBomGetCmd,
	}
)

func init() {
	initBomUploadCmd()
	initBomGetCmd()

	bomCmd.AddCommand(bomUploadCmd)
	bomCmd.AddCommand(bomGetCmd)

	rootCmd.AddCommand(bomCmd)
}

func initBomUploadCmd() {
	bomUploadCmd.Flags().StringP("bom", "b", "", "bom path")
	bomUploadCmd.Flags().Bool("autocreate", false, "automatically create project")

	bomUploadCmd.MarkFlagRequired("bom")
	bomUploadCmd.MarkFlagFilename("bom", "xml", "json")
}

func initBomGetCmd() {
	bomGetCmd.Flags().StringP("output", "o", "", "")
}

func runBomUploadCmd(cmd *cobra.Command, _ []string) {
	bomPath, _ := cmd.Flags().GetString("bom")
	autoCreate, _ := cmd.Flags().GetBool("autocreate")

	log.Println("resolving project")
	_, err := dtrackClient.ResolveProject(pProjectUUID, pProjectName, pProjectVersion)
	if err != nil {
		log.Fatal("failed to resolve project: ", err)
		return
	}

	log.Println("reading bom")
	bomContent, err := ioutil.ReadFile(bomPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("uploading bom")
	_, err = dtrackClient.UploadBOM(dtrack.BOMSubmitRequest{
		ProjectUUID:    pProjectUUID,
		ProjectName:    pProjectName,
		ProjectVersion: pProjectVersion,
		AutoCreate:     autoCreate,
		BOM:            base64.StdEncoding.EncodeToString(bomContent),
	})
	if err != nil {
		log.Fatal("uploading bom failed: ", err)
		return
	}
	log.Println("bom was successfully uploaded")
}

func runBomGetCmd(cmd *cobra.Command, _ []string) {
	log.Println("resolving project")
	project, err := dtrackClient.ResolveProject(pProjectUUID, pProjectName, pProjectVersion)
	if err != nil {
		log.Fatal("failed to resolve project: ", err)
		return
	}

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

	if err = ioutil.WriteFile(output, []byte(bom), 0); err != nil {
		log.Fatal("failed to write output file: ", err)
	}
}
