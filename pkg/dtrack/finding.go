package dtrack

import (
	"context"
	"strconv"
)

type Finding struct {
	Attribution   *FindingAttribution `json:"attribution"`
	Analysis      *Analysis           `json:"analysis"`
	Component     Component           `json:"component"`
	Matrix        string              `json:"matrix"`
	Vulnerability Vulnerability       `json:"vulnerability"`
}

type FindingAttribution struct {
	AlternateIdentifier string `json:"alternateIdentifier"`
	AnalyzerIdentity    string `json:"analyzerIdentity"`
	AttributedOn        int    `json:"attributedOn"`
	ReferenceURL        string `json:"referenceUrl"`
	UUID                string `json:"uuid"`
}

type FindingService interface {
	GetForProject(ctx context.Context, projectUUID string, includeSuppressed bool) ([]Finding, error)
	ExportForProject(ctx context.Context, projectUUID string) (string, error)
}

type findingServiceImpl struct {
	client *Client
}

func (f findingServiceImpl) GetForProject(ctx context.Context, projectUUID string, includeSuppressed bool) ([]Finding, error) {
	findings := make([]Finding, 0)

	req := f.client.restClient.R().
		SetContext(ctx).
		SetPathParam("uuid", projectUUID).
		SetQueryParam("suppressed", strconv.FormatBool(includeSuppressed)).
		SetResult([]Finding{})

	err := f.client.getPaginatedResponse(req, "/api/v1/finding/project/{uuid}", func(res interface{}) (int, error) {
		findingsOnPage, ok := res.(*[]Finding)
		if !ok {
			return -1, ErrInvalidResponseType
		}
		findings = append(findings, *findingsOnPage...)
		return len(findings), nil
	})
	if err != nil {
		return nil, err
	}

	return findings, nil
}

func (f findingServiceImpl) ExportForProject(ctx context.Context, projectUUID string) (string, error) {
	panic("implement me")
}
