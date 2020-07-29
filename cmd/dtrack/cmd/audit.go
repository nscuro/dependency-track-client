package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/nscuro/dependency-track-client/internal/audit"
	"github.com/nscuro/dependency-track-client/pkg/dtrack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "audit for vulnerabilities",
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
	dtrackClient := dtrack.NewClient(viper.GetString("url"), viper.GetString("api-key"))
	bomPath, _ := cmd.Flags().GetString("bom")
	autoCreate, _ := cmd.Flags().GetBool("autocreate")

	log.Println("resolving project")
	project, err := dtrackClient.ResolveProject(pProjectUUID, pProjectName, pProjectVersion)
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
	uploadToken, err := dtrackClient.UploadBOM(dtrack.BOMSubmitRequest{
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
	qualityGate := audit.QualityGate{MaxSeverity: dtrack.HighSeverity}
	if err = qualityGate.Evaluate(findings); err != nil {
		log.Fatal("quality gate failed: ", err)
	}
}
