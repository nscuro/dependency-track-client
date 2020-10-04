package dtrack

import (
	"fmt"
	"strconv"
)

type Project struct {
	Description            string          `json:"description"`
	LastBOMImport          int64           `json:"lastBomImport"`
	LastBOMImportFormat    string          `json:"lastBomImportFormat"`
	LastInheritedRiskScore int             `json:"lastInheritedRiskScore"`
	Metrics                *ProjectMetrics `json:"metrics"`
	Name                   string          `json:"name"`
	UUID                   string          `json:"uuid"`
	Version                string          `json:"version"`
}

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

// GetProject retrieves all projects
func (c Client) GetProjects() ([]Project, error) {
	page := 1
	hasMorePages := true
	projects := make([]Project, 0)
	for hasMorePages {
		res, err := c.restClient.R().
			SetResult(make([]Project, 0)).
			SetQueryParams(map[string]string{
				"pageSize":   "100",
				"pageNumber": strconv.Itoa(page),
			}).
			Get("/api/v1/project")
		if err != nil {
			return nil, err
		}

		if err = c.checkResponseStatus(res, 200); err != nil {
			return nil, err
		}

		projectsOnPage, ok := res.Result().(*[]Project)
		if !ok {
			return nil, ErrInvalidResponseType
		}

		projects = append(projects, *projectsOnPage...)

		totalCount, _ := strconv.Atoi(res.Header().Get(totalCountHeader))
		hasMorePages = len(projects) < totalCount
		page++
	}

	return projects, nil
}
