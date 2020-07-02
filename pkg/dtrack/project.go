package dtrack

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

	return res.Result().(*Project), nil
}
