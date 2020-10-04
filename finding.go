package dtrack

import "strconv"

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
func (c Client) GetFindings(uuid string) ([]Finding, error) {
	page := 1
	hasMorePages := true
	findings := make([]Finding, 0)

	for hasMorePages {
		res, err := c.restClient.R().
			SetPathParams(map[string]string{
				"uuid": uuid,
			}).
			SetQueryParams(map[string]string{
				"pageSize":   "100",
				"pageNumber": strconv.Itoa(page),
			}).
			SetResult(make([]Finding, 0)).
			Get("/api/v1/finding/project/{uuid}")
		if err != nil {
			return nil, err
		}

		if err = c.checkResponseStatus(res, 200); err != nil {
			return nil, err
		}

		findingsOnPage, ok := res.Result().(*[]Finding)
		if !ok {
			return nil, ErrInvalidResponseType
		}

		findings = append(findings, *findingsOnPage...)

		totalCount, _ := strconv.Atoi(res.Header().Get(totalCountHeader))
		hasMorePages = len(findings) < totalCount
		page++
	}

	return findings, nil
}
