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
)

var (
	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Audit for vulnerabilities",
		Run:   runAuditCmd,
	}

	auditOpts AuditOptions
)

type AuditOptions struct {
	BOMPath    string
	AutoCreate bool
}

func init() {
	auditCmd.Flags().StringVarP(&auditOpts.BOMPath, "bom", "b", "", "BOM path")
	auditCmd.Flags().BoolVar(&auditOpts.AutoCreate, "autocreate", false, "Automatically create project")

	auditCmd.MarkFlagRequired("bom")
	auditCmd.MarkFlagFilename("bom", "xml", "json")

	rootCmd.AddCommand(auditCmd)
}

func runAuditCmd(cmd *cobra.Command, _ []string) {
	log.Println("resolving project")
	project, err := dtrackClient.ResolveProject(projectUUID, projectName, projectVersion)
	if err != nil {
		log.Fatal("failed to resolve project: ", err)
		return
	}

	log.Println("reading bom")
	bomContent, err := ioutil.ReadFile(auditOpts.BOMPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("uploading bom")
	uploadToken, err := dtrackClient.UploadBOM(dtrack.BOMSubmitRequest{
		ProjectUUID:    projectUUID,
		ProjectName:    projectName,
		ProjectVersion: projectVersion,
		AutoCreate:     auditOpts.AutoCreate,
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
