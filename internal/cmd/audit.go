package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nscuro/dependency-track-client/internal/audit"
	"github.com/nscuro/dependency-track-client/pkg/dtrack"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit for vulnerabilities",
	Run:   runAuditCmd,
}

func init() {
	auditCmd.Flags().StringP("project", "p", "", "project name")
	auditCmd.Flags().StringP("version", "v", "", "project version")
	auditCmd.Flags().String("uuid", "", "project uuid")
	auditCmd.Flags().StringP("bom", "b", "", "bom path")
	auditCmd.Flags().Bool("autocreate", false, "automatically create project")

	auditCmd.MarkFlagRequired("bom")
	auditCmd.MarkFlagFilename("bom", "xml", "json")

	rootCmd.AddCommand(auditCmd)
}

func runAuditCmd(cmd *cobra.Command, _ []string) {
	projectName, _ := cmd.Flags().GetString("project")
	projectVersion, _ := cmd.Flags().GetString("version")
	projectUUID, _ := cmd.Flags().GetString("uuid")
	bomPath, _ := cmd.Flags().GetString("bom")
	autoCreate, _ := cmd.Flags().GetBool("autocreate")

	if (projectName == "" || projectVersion == "") && projectUUID == "" {
		log.Fatal("either project name and version OR project uuid must be provided")
		return
	}

	dtrackClient := dtrack.NewClient(pBaseURL, pAPIKey)

	log.Println("retrieving project info")
	project, err := dtrackClient.GetProject(projectUUID)
	if err != nil {
		log.Fatal("failed to retrieve project: ", err)
		return
	}

	log.Println("reading bom")
	bomFile, err := os.Open(bomPath)
	if err != nil {
		log.Fatal("failed to read bom: ", err)
		return
	}

	bomContent, err := ioutil.ReadAll(bomFile)
	bomFile.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("uploading bom")
	uploadToken, err := dtrackClient.UploadBOM(dtrack.BOMSubmitRequest{
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

	done := make(chan bool)
	go func(ticker *time.Ticker, timeout <-chan time.Time) {
	loop:
		for {
			select {
			case <-ticker.C:
				processing, err := dtrackClient.IsTokenBeingProcessed(uploadToken)
				if err != nil {
					log.Fatal(err)
					return
				}
				if !processing {
					break loop
				}
				log.Println("still processing")
			case <-timeout:
				log.Println("timeout exceeded")
				break loop
			default:
			}
		}
		done <- true
	}(time.NewTicker(5*time.Second), time.After(time.Duration(30)*time.Second))

	log.Printf("waiting for bom processing to complete")
	ticker := time.NewTicker(time.Second)
loop:
	for {
		select {
		case <-ticker.C:
			fmt.Printf(".")
		case <-done:
			break loop
		}
	}
	fmt.Println()
	log.Println("processing completed")

	log.Println("retrieving findings")
	findings, err := dtrackClient.GetFindings(project.UUID)
	if err != nil {
		log.Fatal(err)
		return
	}

	// TODO: load quality gate from file
	log.Println("evaluating quality gate")
	qualityGate := audit.QualityGate{MaxRiskScore: 10}
	if err = qualityGate.Evaluate(findings); err != nil {
		log.Fatal("quality gate failed: ", err)
	}
}
