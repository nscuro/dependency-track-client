package dtrack

type License struct {
	Groups      []LicenseGroup `json:"licenseGroups"`
	LicenseID   string         `json:"licenseId"`
	Name        string         `json:"name"`
	OSIApproved bool           `json:"isOsiApproved"`
	Text        string         `json:"licenseText"`
	UUID        string         `json:"uuid"`
}

func (c Client) GetLicenses() ([]License, error) {
	licenses := make([]License, 0)

	request := c.restClient.R().
		SetResult(make([]License, 0))

	err := c.getPaginatedResponse(request, "/api/v1/license", func(result interface{}) (int, error) {
		l, ok := result.(*[]License)
		if !ok {
			return -1, ErrInvalidResponseType
		}
		licenses = append(licenses, *l...)
		return len(licenses), nil
	})
	if err != nil {
		return nil, err
	}

	return licenses, nil
}

func (c Client) GetLicense(id string) (*License, error) {
	return nil, nil
}
