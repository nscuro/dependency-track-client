package dtrack

type ProjectMetrics struct {
	Components           int   `json:"components"`
	Critical             int   `json:"critical"`
	FindingsAudited      int   `json:"findingsAudited"`
	FindingsTotal        int   `json:"findingsTotal"`
	FindingsUnaudited    int   `json:"findingsUnaudited"`
	High                 int   `json:"high"`
	InheritedRiskScore   int64 `json:"inheritedRiskScore"`
	Low                  int   `json:"low"`
	Medium               int   `json:"medium"`
	Suppressed           int   `json:"suppressed"`
	Unassigned           int   `json:"unassigned"`
	VulnerableComponents int   `json:"vulnerableComponents"`
}

func (c Client) GetCurrentProjectMetrics(uuid string) (*ProjectMetrics, error) {
	res, err := c.restClient.R().
		SetPathParams(map[string]string{
			"uuid": uuid,
		}).
		SetResult(&ProjectMetrics{}).
		Get("/api/v1/metrics/project/{uuid}/current")
	if err != nil {
		return nil, err
	}

	if err = c.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	metrics, ok := res.Result().(*ProjectMetrics)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return metrics, nil
}
