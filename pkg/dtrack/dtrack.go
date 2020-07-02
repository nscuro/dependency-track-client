package dtrack

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

var (
	ErrForbidden     = errors.New("forbidden")
	ErrInvalidStatus = errors.New("invalid status")
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
)

type Client struct {
	restClient *resty.Client
}

func NewClient(baseURL string, apiKey string) *Client {
	restClient := resty.New()
	restClient.SetHostURL(baseURL)
	restClient.SetHeader("X-Api-Key", apiKey)
	restClient.SetHeader("Accept", "application/json")

	return &Client{restClient: restClient}
}

func (c Client) checkResponse(response *resty.Response, expectedStati ...int) error {
	switch response.StatusCode() {
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	}

	if len(expectedStati) > 0 {
		for _, status := range expectedStati {
			if response.StatusCode() == status {
				return nil
			}
		}
		return ErrInvalidStatus
	}

	return nil
}
