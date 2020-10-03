package dtrack

type About struct {
	UUID        string         `json:"uuid"`
	SystemUUID  string         `json:"systemUuid"`
	Application string         `json:"application"`
	Version     string         `json:"version"`
	Timestamp   string         `json:"timestamp"`
	Framework   AboutFramework `json:"framework"`
}

type AboutFramework struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

func (c Client) GetAbout() (*About, error) {
	res, err := c.restClient.R().
		SetResult(&About{}).
		Get("/api/version")
	if err != nil {
		return nil, err
	}

	if err = c.checkResponse(res, 200); err != nil {
		return nil, err
	}

	about, ok := res.Result().(*About)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return about, nil
}
