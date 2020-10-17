package dtrack

type Finding struct {
	Attribution   *FindingAttribution `json:"attribution"`
	Analysis      *Analysis           `json:"analysis"`
	Component     Component           `json:"component"`
	Matrix        string              `json:"matrix"`
	Vulnerability Vulnerability       `json:"vulnerability"`
}

type FindingAttribution struct {
	UUID             string `json:"uuid"`
	AnalyzerIdentity string `json:"analyzerIdentity"`
}

// GetFindings retrieves all findings associated with a given project
func (c Client) GetFindings(projectUUID string) ([]Finding, error) {
	findings := make([]Finding, 0)

	req := c.restClient.R().
		SetPathParams(map[string]string{
			"uuid": projectUUID,
		}).
		SetResult([]Finding{})

	err := c.getPaginatedResponse(req, "/api/v1/finding/project/{uuid}", func(result interface{}) (int, error) {
		findingsOnPage, ok := result.(*[]Finding)
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
