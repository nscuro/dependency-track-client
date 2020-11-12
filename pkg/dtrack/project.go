package dtrack

import (
	"fmt"
)

type Project struct {
	Description            string          `json:"description"`
	LastBOMImport          int64           `json:"lastBomImport"`
	LastBOMImportFormat    string          `json:"lastBomImportFormat"`
	LastInheritedRiskScore float32         `json:"lastInheritedRiskScore"`
	Metrics                *ProjectMetrics `json:"metrics"`
	Name                   string          `json:"name"`
	UUID                   string          `json:"uuid"`
	Version                string          `json:"version"`
}

// GetProject gets a Project by its UUID
func (c Client) GetProject(uuid string) (*Project, error) {
	res, err := c.restClient.R().
		SetPathParams(map[string]string{
			"uuid": uuid,
		}).
		SetResult(&Project{}).
		Get("/api/v1/project/{uuid}")
	if err != nil {
		return nil, err
	}

	if err = c.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	project, ok := res.Result().(*Project)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return project, nil
}

// LookupProject gets a Project by its name and version
func (c Client) LookupProject(name, version string) (*Project, error) {
	res, err := c.restClient.R().
		SetQueryParams(map[string]string{
			"name":    name,
			"version": version,
		}).
		SetResult(&Project{}).
		Get("/api/v1/project/lookup")
	if err != nil {
		return nil, err
	}

	if err = c.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	project, ok := res.Result().(*Project)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return project, nil
}

// ResolveProject is a convenience method that either calls GetProject or LookupProject,
// depending on whether a UUID or name and version are provided
func (c Client) ResolveProject(uuid, name, version string) (*Project, error) {
	if uuid == "" && (name == "" || version == "") {
		return nil, fmt.Errorf("either project uuid or name AND version must be provided")
	}

	if uuid != "" {
		return c.GetProject(uuid)
	} else {
		return c.LookupProject(name, version)
	}
}

// GetProjects retrieves all projects
func (c Client) GetProjects() ([]Project, error) {
	projects := make([]Project, 0)

	req := c.restClient.R().
		SetResult([]Project{})

	err := c.getPaginatedResponse(req, "/api/v1/project", func(result interface{}) (int, error) {
		projectsOnPage, ok := result.(*[]Project)
		if !ok {
			return -1, ErrInvalidResponseType
		}
		projects = append(projects, *projectsOnPage...)
		return len(projects), nil
	})
	if err != nil {
		return nil, err
	}

	return projects, nil
}
