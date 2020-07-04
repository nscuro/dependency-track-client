package report

import (
	"html/template"
	"io"
	"log"

	"github.com/nscuro/dependency-track-client/pkg/dtrack"
)

type ProjectReportContext struct {
	Project    dtrack.Project
	Components []dtrack.Component
	Findings   []dtrack.Finding
}

type Generator struct {
	dtrackClient *dtrack.Client
}

func NewGenerator(dtrackClient *dtrack.Client) *Generator {
	return &Generator{dtrackClient: dtrackClient}
}

func (g Generator) GenerateProjectReport(projectUUID string, templatePath string, writer io.Writer) error {
	log.Println("parsing template file")
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	log.Println("retrieving project info")
	project, err := g.dtrackClient.GetProject(projectUUID)
	if err != nil {
		return err
	}

	log.Println("retrieving project dependencies")
	dependencies, err := g.dtrackClient.GetDependenciesForProject(projectUUID)
	if err != nil {
		return err
	}

	components := make([]dtrack.Component, len(dependencies))
	for i := 0; i < len(dependencies); i++ {
		components[i] = dependencies[i].Component
	}

	log.Println("retrieving findings for project")
	findings, err := g.dtrackClient.GetFindings(projectUUID)
	if err != nil {
		return err
	}

	reportContext := ProjectReportContext{
		Project:    *project,
		Components: components,
		Findings:   findings,
	}

	log.Println("writing report")
	return tpl.Execute(writer, reportContext)
}
