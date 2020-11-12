package dtrack

type LicenseGroup struct {
	Name       string `json:"name"`
	RiskWeight int    `json:"riskWeight"`
	UUID       string `json:"uuid"`
}

func (c Client) GetLicenseGroups() ([]LicenseGroup, error) {
	groups := make([]LicenseGroup, 0)

	req := c.restClient.R().
		SetResult([]LicenseGroup{})

	err := c.getPaginatedResponse(req, "/api/v1/licenseGroup", func(result interface{}) (int, error) {
		groupsOnPage, ok := result.(*[]LicenseGroup)
		if !ok {
			return -1, ErrInvalidResponseType
		}
		groups = append(groups, *groupsOnPage...)
		return len(groups), nil
	})
	if err != nil {
		return nil, err
	}

	return groups, nil
}
