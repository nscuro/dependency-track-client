package dtrack

type Component struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Group   string `json:"group"`
	Version string `json:"version"`
}

func (c Client) GetComponent(uuid string) (*Component, error) {
	res, err := c.restClient.R().
		SetPathParams(map[string]string{
			"uuid": uuid,
		}).
		Get("/api/v1/component/{uuid}")
	if err != nil {
		return nil, err
	}

	if err = c.checkResponse(res, 200); err != nil {
		return nil, err
	}

	return res.Result().(*Component), nil
}
