package dtrack

type Finding struct {
	Component     Component     `json:"component"`
	Vulnerability Vulnerability `json:"vulnerability"`
	Analysis      Analysis      `json:"analysis"`
	Matrix        string        `json:"matrix"`
}

func (c Client) GetFindings(uuid string) ([]Finding, error) {
	res, err := c.restClient.R().
		SetPathParams(map[string]string{
			"uuid": uuid,
		}).
		SetResult(make([]Finding, 0)).
		Get("/api/v1/finding/project/{uuid}")
	if err != nil {
		return nil, err
	}

	if err = c.checkResponse(res, 200); err != nil {
		return nil, err
	}

	findings, ok := res.Result().(*[]Finding)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return *findings, nil
}
