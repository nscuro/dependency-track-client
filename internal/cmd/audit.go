package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/nscuro/dependency-track-client"
	"github.com/nscuro/dependency-track-client/internal/audit"
)

var (
	auditCmd = &cobra.Command{
		Use:   "audit",
		Short: "Audit for vulnerabilities",
		Run:   runAuditCmd,
	}
	auditOpts auditOptions
)

type auditOptions struct {
	bomFilePath string
	autoCreate  bool
}

func init() {
	auditCmd.Flags().StringVar(&auditOpts.bomFilePath, "bom", "", "BOM file path")
	auditCmd.Flags().BoolVar(&auditOpts.autoCreate, "autocreate", false, "Automatically create project")

	rootCmd.AddCommand(auditCmd)
}

func runAuditCmd(_ *cobra.Command, _ []string) {
	dtrackClient := mustGetDTrackClient()
	project := mustResolveProject(dtrackClient)

	log.Printf("reading BOM from %s\n", auditOpts.bomFilePath)
	bomContent, err := ioutil.ReadFile(auditOpts.bomFilePath)
	if err != nil {
		log.Fatalf("failed to read BOM file: %v", err)
	}

	log.Println("uploading BOM")
	token, err := dtrackClient.UploadBOM(dtrack.BOMSubmitRequest{
		ProjectUUID:    project.UUID,
		ProjectName:    project.Name,
		ProjectVersion: project.Version,
		AutoCreate:     auditOpts.autoCreate,
		BOM:            base64.StdEncoding.EncodeToString(bomContent),
	})
	if err != nil {
		log.Fatalf("failed to upload BOM: %v", err)
	}

	done := make(chan bool)
	go func(ticker *time.Ticker, timeout <-chan time.Time) {
	loop:
		for {
			select {
			case <-ticker.C:
				if processing, err := dtrackClient.IsTokenBeingProcessed(token); err == nil {
					if !processing {
						break loop
					}
				} else {
					log.Fatalf("failed to get processing status for token %s: %v", token, err)
				}
			case <-timeout:
				log.Fatalln("timeout exceeded")
			}
		}
		done <- true
	}(time.NewTicker(5*time.Second), time.After(30*time.Second))

	log.Println("waiting for BOM processing to complete")
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

	fmt.Println("retrieving findings")
	findings, err := dtrackClient.GetFindings(project.UUID)
	if err != nil {
		log.Fatalf("failed to retrieve findings: %v", err)
	}

	fmt.Println("retrieving project metrics")
	metrics, err := dtrackClient.GetCurrentProjectMetrics(project.UUID)
	if err != nil {
		log.Fatalf("failed to retrieve project metrics: %v", err)
	}

	log.Println("evaluating quality gate")
	gate := audit.QualityGate{
		MaxSeverity: dtrack.HighSeverity,
	}
	if err = gate.Evaluate(*metrics, findings); err != nil {
		log.Fatalf("quality gate failed: %v", err)
	}
}
