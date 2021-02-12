package dtrack

import "context"

type Analysis struct {
	Comments   []AnalysisComment `json:"comments"`
	State      string            `json:"state"`
	Suppressed bool              `json:"isSuppressed"`
}

type AnalysisComment struct {
	Comment   string `json:"comment"`
	Commenter string `json:"commenter"`
	Timestamp string `json:"timestamp"`
}

type AnalysisRequest struct {
	ComponentUUID     string `json:"component"`
	ProjectUUID       string `json:"project"`
	VulnerabilityUUID string `json:"vulnerability"`
	Comment           string `json:"comment,omitempty"`
	State             string `json:"analysisState,omitempty"`
	Suppressed        bool   `json:"isSuppressed"`
}

type AnalysisService interface {
	Create(ctx context.Context, req AnalysisRequest) (*Analysis, error)
	Get(ctx context.Context, cUUID, pUUID, vUUID string) (*Analysis, error)
}

type analysisServiceImpl struct {
	client *Client
}

func (a analysisServiceImpl) Create(ctx context.Context, req AnalysisRequest) (*Analysis, error) {
	res, err := a.client.restClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&Analysis{}).
		Put("/api/v1/analysis")
	if err != nil {
		return nil, err
	}

	if err = a.client.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	analysis, ok := res.Result().(*Analysis)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return analysis, nil
}

func (a analysisServiceImpl) Get(ctx context.Context, cUUID, pUUID, vUUID string) (*Analysis, error) {
	res, err := a.client.restClient.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"component":     cUUID,
			"project":       pUUID,
			"vulnerability": vUUID,
		}).
		SetResult(&Analysis{}).
		Get("/api/v1/analysis")
	if err != nil {
		return nil, err
	}

	if err = a.client.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	analysis, ok := res.Result().(*Analysis)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return analysis, nil
}
