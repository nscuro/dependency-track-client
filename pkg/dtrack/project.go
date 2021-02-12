package dtrack

import (
	"context"
	"fmt"
)

type Project struct {
	Author                 string            `json:"author"`
	Classifier             string            `json:"classifier"`
	CPE                    string            `json:"cpe"`
	Description            string            `json:"description"`
	Group                  string            `json:"group"`
	LastBOMImport          int64             `json:"lastBomImport"`
	LastBOMImportFormat    string            `json:"lastBomImportFormat"`
	LastInheritedRiskScore float32           `json:"lastInheritedRiskScore"`
	Metrics                *ProjectMetrics   `json:"metrics"`
	Name                   string            `json:"name"`
	PackageURL             string            `json:"purl"`
	Properties             []ProjectProperty `json:"properties"`
	Publisher              string            `json:"publisher"`
	SWIDTagID              string            `json:"swidTagId"`
	Tags                   []ProjectTag      `json:"tags"`
	UUID                   string            `json:"uuid"`
	Version                string            `json:"version"`
}

type ProjectProperty struct {
	Group string `json:"groupName"`
	Name  string `json:"propertyName"`
	Type  string `json:"propertyType"`
	Value string `json:"propertyValue"`
}

type ProjectTag struct {
	Name string `json:"name"`
}

type ProjectService interface {
	GetAll(ctx context.Context) ([]Project, error)
	GetByUUID(ctx context.Context, uuid string) (*Project, error)
	Lookup(ctx context.Context, name, version string) (*Project, error)
	Resolve(ctx context.Context, uuid, name, version string) (*Project, error)
}

type projectServiceImpl struct {
	client *Client
}

func (p projectServiceImpl) GetAll(ctx context.Context) ([]Project, error) {
	projects := make([]Project, 0)

	req := p.client.restClient.R().
		SetContext(ctx).
		SetResult([]Project{})

	err := p.client.getPaginatedResponse(req, "/api/v1/project", func(result interface{}) (int, error) {
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

func (p projectServiceImpl) GetByUUID(ctx context.Context, uuid string) (*Project, error) {
	res, err := p.client.restClient.R().
		SetContext(ctx).
		SetPathParam("uuid", uuid).
		SetResult(&Project{}).
		Get("/api/v1/project/{uuid}")
	if err != nil {
		return nil, err
	}

	if err = p.client.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	project, ok := res.Result().(*Project)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return project, nil
}

func (p projectServiceImpl) Lookup(ctx context.Context, name, version string) (*Project, error) {
	res, err := p.client.restClient.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"name":    name,
			"version": version,
		}).
		SetResult(&Project{}).
		Get("/api/v1/project/lookup")
	if err != nil {
		return nil, err
	}

	if err = p.client.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	project, ok := res.Result().(*Project)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return project, nil
}

func (p projectServiceImpl) Resolve(ctx context.Context, uuid, name, version string) (*Project, error) {
	if uuid == "" && (name == "" || version == "") {
		return nil, fmt.Errorf("either project uuid or name AND version must be provided")
	}

	if uuid != "" {
		return p.GetByUUID(ctx, uuid)
	} else {
		return p.Lookup(ctx, name, version)
	}
}
