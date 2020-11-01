package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/nscuro/dependency-track-client"
	"github.com/nscuro/dependency-track-client/internal/qualitygate"
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
	autoCreate   bool
	bomFilePath  string
	gateFilePath string
	timeout      time.Duration
}

func init() {
	flags := auditCmd.Flags()

	flags.BoolVar(&auditOpts.autoCreate, "autocreate", false, "Automatically create project")
	flags.StringVar(&auditOpts.bomFilePath, "bom", "", "BOM file path")
	flags.StringVarP(&auditOpts.gateFilePath, "gate", "g", "", "Quality gate file path")
	flags.DurationVar(&auditOpts.timeout, "timeout", 30, "Timeout in seconds")

	rootCmd.AddCommand(auditCmd)
}

func runAuditCmd(_ *cobra.Command, _ []string) {
	dtrackClient := mustGetDTrackClient()

	log.Printf("reading BOM from %s\n", auditOpts.bomFilePath)
	bomContent, err := ioutil.ReadFile(auditOpts.bomFilePath)
	if err != nil {
		log.Fatalf("failed to read BOM file: %v", err)
	}

	log.Printf("loading quality gate from %s\n", auditOpts.gateFilePath)
	gate, err := qualitygate.LoadGateFromFile(auditOpts.gateFilePath)
	if err != nil {
		log.Fatalf("failed to read quality gate file: %v", err)
	}

	log.Println("uploading BOM")
	token, err := dtrackClient.UploadBOM(dtrack.BOMUploadRequest{
		ProjectUUID:    globalOpts.projectUUID,
		ProjectName:    globalOpts.projectName,
		ProjectVersion: globalOpts.projectVersion,
		AutoCreate:     auditOpts.autoCreate,
		BOM:            base64.StdEncoding.EncodeToString(bomContent),
	})
	if err != nil {
		log.Fatalf("failed to upload BOM: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(ticker *time.Ticker, timeout <-chan time.Time) {
		defer wg.Done()
	loop:
		for {
			select {
			case <-ticker.C:
				if processing, err := dtrackClient.IsTokenBeingProcessed(token); err == nil {
					if !processing {
						break loop
					}
					fmt.Print(".")
				} else {
					log.Fatalf("failed to get processing status for token %s: %v", token, err)
				}
			case <-timeout:
				log.Fatalln("timeout exceeded")
			}
		}
	}(time.NewTicker(5*time.Second), time.After(auditOpts.timeout))

	log.Println("waiting for BOM processing to complete")
	wg.Wait()

	var projectUUID string
	if globalOpts.projectUUID != "" {
		projectUUID = globalOpts.projectUUID
	} else {
		if project, err := dtrackClient.LookupProject(globalOpts.projectName, globalOpts.projectVersion); err == nil {
			projectUUID = project.UUID
		} else {
			log.Fatalf("failed to lookup project: %v", err)
		}
	}

	log.Println("evaluating quality gate")
	if err = qualitygate.NewEvaluator(dtrackClient).Evaluate(projectUUID, gate); err != nil {
		log.Fatalf("quality gate failed: %v", err)
	}
}
