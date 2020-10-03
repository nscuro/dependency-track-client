package dtrack

import (
	"fmt"
	"strconv"
)

type Project struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
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

	if err = c.checkResponse(res, 200); err != nil {
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

	if err = c.checkResponse(res, 200); err != nil {
		return nil, err
	}

	project, ok := res.Result().(*Project)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return project, nil
}

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

		if err = c.checkResponse(res, 200); err != nil {
			return nil, err
		}

		projectsOnPage, ok := res.Result().(*[]Project)
		if !ok {
			return nil, ErrInvalidResponseType
		}

		projects = append(projects, *projectsOnPage...)

		totalCount, _ := strconv.Atoi(res.Header().Get("X-Total-Count"))

		hasMorePages = len(projects) < totalCount
		page++
	}

	return projects, nil
}
