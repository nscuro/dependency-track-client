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
	Licenses   []dtrack.License
}

type Generator struct {
	dtrackClient *dtrack.Client
}

func NewGenerator(dtrackClient *dtrack.Client) *Generator {
	return &Generator{dtrackClient: dtrackClient}
}

func (g Generator) GenerateProjectReport(project *dtrack.Project, templatePath string, writer io.Writer) error {
	log.Println("parsing template file")
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	log.Println("retrieving project dependencies")
	dependencies, err := g.dtrackClient.GetDependenciesForProject(project.UUID)
	if err != nil {
		return err
	}

	components := make([]dtrack.Component, len(dependencies))
	for i := 0; i < len(dependencies); i++ {
		components[i] = dependencies[i].Component
	}

	log.Println("retrieving findings for project")
	findings, err := g.dtrackClient.GetFindings(project.UUID)
	if err != nil {
		return err
	}

	log.Println("collecting licenses")
	licenses := make([]dtrack.License, 0)
	for _, component := range components {
		if component.ResolvedLicense.UUID == "" {
			continue
		}

		alreadyAdded := false
		for _, license := range licenses {
			if component.ResolvedLicense.UUID == license.UUID {
				alreadyAdded = true
				break
			}
		}
		if !alreadyAdded {
			licenses = append(licenses, component.ResolvedLicense)
		}
	}

	reportContext := ProjectReportContext{
		Project:    *project,
		Components: components,
		Findings:   findings,
		Licenses:   licenses,
	}

	log.Println("writing report")
	return tpl.Execute(writer, reportContext)
}
