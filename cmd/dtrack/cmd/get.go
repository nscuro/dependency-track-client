package cmd

import (
	"fmt"
	"github.com/nscuro/dependency-track-client/pkg/dtrack"
	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Access API resources",
	}

	getProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "Show project information",
		RunE:  runGetProjectCmd,
	}

	getProjectsCmd = &cobra.Command{
		Use:   "projects",
		Short: "Show all projects",
		RunE:  runGetProjectsCmd,
	}

	getDependenciesCmd = &cobra.Command{
		Use:   "deps",
		Short: "Show dependencies of a project",
		RunE:  runGetDependenciesCmd,
	}

	getFindingsCmd = &cobra.Command{
		Use:   "findings",
		Short: "Show findings of a project",
		RunE:  runGetFindingsCmd,
	}

	getVulnerabilityCmd = &cobra.Command{
		Use:   "vuln",
		Short: "Show vulnerability information",
		RunE:  runGetVulnerabilityCmd,
	}
)

type GetVulnOptions struct {
	UUID   string
	VulnID string
}

var (
	getVulnOpts GetVulnOptions
)

func init() {
	getCmd.AddCommand(getProjectCmd)
	getCmd.AddCommand(getProjectsCmd)
	getCmd.AddCommand(getDependenciesCmd)
	getCmd.AddCommand(getFindingsCmd)

	getVulnerabilityCmd.Flags().StringVar(&getVulnOpts.UUID, "uuid", "", "Vulnerability UUID")
	getVulnerabilityCmd.Flags().StringVar(&getVulnOpts.VulnID, "id", "", "Vulnerability ID (e.g. as CVE)")
	getCmd.AddCommand(getVulnerabilityCmd)

	rootCmd.AddCommand(getCmd)
}

func runGetProjectCmd(_ *cobra.Command, _ []string) error {
	project, err := dtrackClient.ResolveProject(projectUUID, projectName, projectVersion)
	if err != nil {
		return fmt.Errorf("failed to resolve project: %w", err)
	}

	fmt.Printf("%s %s %s\n", project.UUID, project.Name, projectVersion)
	return nil
}

func runGetProjectsCmd(_ *cobra.Command, _ []string) error {
	projects, err := dtrackClient.GetProjects()
	if err != nil {
		return fmt.Errorf("failed to retrieve projects: %w", err)
	}

	for _, project := range projects {
		fmt.Printf("%s %s %s\n", project.UUID, project.Name, project.Version)
	}
	return nil
}

func runGetDependenciesCmd(_ *cobra.Command, _ []string) error {
	var uuid string
	if projectUUID != "" {
		uuid = projectUUID
	} else {
		project, err := dtrackClient.ResolveProject(projectUUID, projectName, projectVersion)
		if err != nil {
			return fmt.Errorf("failed to resolve project: %w", err)
		}
		uuid = project.UUID
	}

	dependencies, err := dtrackClient.GetDependenciesForProject(uuid)
	if err != nil {
		return fmt.Errorf("failed to retrieve dependencies: %w", err)
	}

	for _, dep := range dependencies {
		fmt.Printf("%s %s %s\n", dep.Component.Name, dep.Component.Group, dep.Component.Version)
	}

	return nil
}

func runGetFindingsCmd(_ *cobra.Command, _ []string) error {
	var uuid string
	if projectUUID != "" {
		uuid = projectUUID
	} else {
		project, err := dtrackClient.ResolveProject(projectUUID, projectName, projectVersion)
		if err != nil {
			return fmt.Errorf("failed to resolve project: %w", err)
		}
		uuid = project.UUID
	}

	findings, err := dtrackClient.GetFindings(uuid)
	if err != nil {
		return fmt.Errorf("failed to retrieve findings: %w", err)
	}

	for _, finding := range findings {
		fmt.Printf("%s %s", finding.Vulnerability.VulnID, finding.Vulnerability.Severity)
	}

	return nil
}

func runGetVulnerabilityCmd(_ *cobra.Command, _ []string) error {
	var err error
	var vuln *dtrack.Vulnerability
	if getVulnOpts.UUID != "" {
		vuln, err = dtrackClient.GetVulnerability(getVulnOpts.UUID)
		if err != nil {
			return fmt.Errorf("failed to retrieve vulnerability: %w", err)
		}
	} else if getVulnOpts.VulnID != "" {
		source := dtrackClient.GuessVulnerabilitySource(getVulnOpts.VulnID)
		vuln, err = dtrackClient.GetVulnerabilityByVulnID(getVulnOpts.VulnID, source)
		if err != nil {
			return fmt.Errorf("failed to retrieve vulnerability: %w", err)
		}
	} else {
		return fmt.Errorf("either vulnerability ID or UUID must be provided")
	}

	fmt.Printf("%s %s %s\n", vuln.UUID, vuln.VulnID, vuln.Severity)
	return nil
}
