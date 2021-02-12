package dtrack

import "context"

type Component struct {
	Author          string  `json:"author"`
	Blake2b_256     string  `json:"blake2b_256"`
	Blake2b_384     string  `json:"blake2b_384"`
	Blake2b_512     string  `json:"blake2b_512"`
	Blake3          string  `json:"blake3"`
	Classifier      string  `json:"classifier"`
	Copyright       string  `json:"copyright"`
	CPE             string  `json:"cpe"`
	Extension       string  `json:"extension"`
	Filename        string  `json:"filename"`
	Group           string  `json:"group"`
	Internal        bool    `json:"isInternal"`
	License         string  `json:"license"`
	MD5             string  `json:"md5"`
	Name            string  `json:"name"`
	PackageURL      string  `json:"purl"`
	Publisher       string  `json:"publisher"`
	ResolvedLicense License `json:"resolvedLicense"`
	SHA1            string  `json:"sha1"`
	SHA256          string  `json:"sha256"`
	SHA384          string  `json:"sha384"`
	SHA512          string  `json:"sha512"`
	SHA3_256        string  `json:"sha3_256"`
	SHA3_384        string  `json:"sha3_384"`
	SHA3_512        string  `json:"sha3_512"`
	SWIDTagID       string  `json:"swidTagId"`
	UUID            string  `json:"uuid"`
	Version         string  `json:"version"`
}

type ComponentService interface {
	GetAllForProject(ctx context.Context, projectUUID string) ([]Component, error)
	GetByUUID(ctx context.Context, uuid string) (*Component, error)
	GetByHash(ctx context.Context, hash string) (*Component, error)
}

type componentServiceImpl struct {
	client *Client
}

func (c componentServiceImpl) GetAllForProject(ctx context.Context, projectUUID string) ([]Component, error) {
	components := make([]Component, 0)

	req := c.client.restClient.R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"uuid": projectUUID,
		}).
		SetResult([]Component{})

	err := c.client.getPaginatedResponse(req, "/api/v1/component/project/{uuid}", func(res interface{}) (int, error) {
		componentsOnPage, ok := res.(*[]Component)
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

func (c componentServiceImpl) GetByUUID(ctx context.Context, uuid string) (*Component, error) {
	res, err := c.client.restClient.R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"uuid": uuid,
		}).
		SetResult(&Component{}).
		Get("/api/v1/component/{uuid}")
	if err != nil {
		return nil, err
	}

	if err = c.client.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	component, ok := res.Result().(*Component)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return component, nil
}

func (c componentServiceImpl) GetByHash(ctx context.Context, hash string) (*Component, error) {
	res, err := c.client.restClient.R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"hash": hash,
		}).
		SetResult(&Component{}).
		Get("/api/v1/component/hash/{hash}")
	if err != nil {
		return nil, err
	}

	if err = c.client.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	component, ok := res.Result().(*Component)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return component, nil
}
