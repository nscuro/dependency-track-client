package dtrack

type Component struct {
	UUID            string  `json:"uuid"`
	Name            string  `json:"name"`
	Group           string  `json:"group"`
	Version         string  `json:"version"`
	PackageURL      string  `json:"purl"`
	Internal        bool    `json:"isInternal"`
	ResolvedLicense License `json:"resolvedLicense"`
}

func (c Client) GetComponent(uuid string) (*Component, error) {
	res, err := c.restClient.R().
		SetPathParams(map[string]string{
			"uuid": uuid,
		}).
		SetResult(&Component{}).
		Get("/api/v1/component/{uuid}")
	if err != nil {
		return nil, err
	}

	if err = c.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	component, ok := res.Result().(*Component)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return component, nil
}

func (c Client) GetComponentsForProject(projectUUID string) ([]Component, error) {
	components := make([]Component, 0)

	req := c.restClient.R().
		SetPathParams(map[string]string{
			"uuid": projectUUID,
		}).
		SetResult([]Component{})

	err := c.getPaginatedResponse(req, "/api/v1/component/project/{uuid}", func(result interface{}) (int, error) {
		componentsOnPage, ok := result.(*[]Component)
		if !ok {
			return -1, ErrInvalidResponseType
		}
		components = append(components, *componentsOnPage...)
		return len(components), nil
	})
	if err != nil {
		return nil, err
	}

	return components, nil
}
