package dtrack

import "context"

type RepositoryMetaComponent struct {
	LatestVersion string `json:"latestVersion"`
}

type RepositoryService interface {
	GetMetaComponent(ctx context.Context, purl string) (*RepositoryMetaComponent, error)
}

type repositoryServiceImpl struct {
	client *Client
}

func (r repositoryServiceImpl) GetMetaComponent(ctx context.Context, purl string) (*RepositoryMetaComponent, error) {
	res, err := r.client.restClient.R().
		SetContext(ctx).
		SetQueryParam("purl", purl).
		SetResult(&RepositoryMetaComponent{}).
		Get("/api/v1/repository/latest")
	if err != nil {
		return nil, err
	}

	if err = r.client.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	project, ok := res.Result().(*RepositoryMetaComponent)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return project, nil
}
