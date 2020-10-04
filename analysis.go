package dtrack

type Analysis struct {
	Comments   []AnalysisComment
	State      string `json:"state"`
	Suppressed bool   `json:"isSuppressed"`
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
	Comment           string `json:"comment"`
	State             string `json:"analysisState"`
	Suppressed        bool   `json:"isSuppressed"`
}

func (c Client) GetAnalysis(componentUUID, projectUUID, vulnerabilityUUID string) (*Analysis, error) {
	res, err := c.restClient.R().
		SetQueryParams(map[string]string{
			"component":     componentUUID,
			"project":       projectUUID,
			"vulnerability": vulnerabilityUUID,
		}).
		SetResult(&Analysis{}).
		Get("/api/v1/analysis")
	if err != nil {
		return nil, err
	}

	if err = c.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	analysis, ok := res.Result().(*Analysis)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return analysis, nil
}

func (c Client) RecordAnalysis(req AnalysisRequest) (*Analysis, error) {
	res, err := c.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&Analysis{}).
		Put("/api/v1/analysis")
	if err != nil {
		return nil, err
	}

	if err = c.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	analysis, ok := res.Result().(*Analysis)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return analysis, nil
}
