package dtrack

import "context"

type License struct {
	Comment             string         `json:"comment"`
	DeprecatedLicenseID bool           `json:"deprecatedLicenseId"`
	FSFLibre            bool           `json:"fsfLibre"`
	Groups              []LicenseGroup `json:"licenseGroups"`
	Header              string         `json:"header"`
	LicenseID           string         `json:"licenseId"`
	Name                string         `json:"name"`
	OSIApproved         bool           `json:"isOsiApproved"`
	SeeAlso             []string       `json:"seeAlso"`
	Text                string         `json:"licenseText"`
	UUID                string         `json:"uuid"`
}

type LicenseService interface {
	GetAll(ctx context.Context) ([]License, error)
	GetByID(ctx context.Context, id string) (*License, error)
}

type licenseServiceImpl struct {
	client *Client
}

func (l licenseServiceImpl) GetAll(ctx context.Context) ([]License, error) {
	licenses := make([]License, 0)

	request := l.client.restClient.R().
		SetContext(ctx).
		SetResult(make([]License, 0))

	err := l.client.getPaginatedResponse(request, "/api/v1/license", func(res interface{}) (int, error) {
		l, ok := res.(*[]License)
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

func (l licenseServiceImpl) GetByID(ctx context.Context, id string) (*License, error) {
	panic("implement me")
}
