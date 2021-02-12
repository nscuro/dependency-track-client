package dtrack

import "context"

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

type AboutService interface {
	Get(ctx context.Context) (*About, error)
}

type aboutServiceImpl struct {
	client *Client
}

func (a aboutServiceImpl) Get(ctx context.Context) (*About, error) {
	res, err := a.client.restClient.R().
		SetContext(ctx).
		SetResult(&About{}).
		Get("/api/version")
	if err != nil {
		return nil, err
	}

	if err = a.client.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	about, ok := res.Result().(*About)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return about, nil
}
