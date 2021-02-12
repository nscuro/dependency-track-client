package dtrack

import "context"

type LicenseGroup struct {
	Name       string `json:"name"`
	RiskWeight int    `json:"riskWeight"`
	UUID       string `json:"uuid"`
}

type LicenseGroupService interface {
	GetAll(ctx context.Context) ([]LicenseGroup, error)
}

type licenseGroupServiceImpl struct {
	client *Client
}

func (l licenseGroupServiceImpl) GetAll(ctx context.Context) ([]LicenseGroup, error) {
	groups := make([]LicenseGroup, 0)

	req := l.client.restClient.R().
		SetContext(ctx).
		SetResult([]LicenseGroup{})

	err := l.client.getPaginatedResponse(req, "/api/v1/licenseGroup", func(res interface{}) (int, error) {
		groupsOnPage, ok := res.(*[]LicenseGroup)
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
