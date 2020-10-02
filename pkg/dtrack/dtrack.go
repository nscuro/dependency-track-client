package dtrack

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

var (
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidResponseType = errors.New("invalid response type")
)

type Client struct {
	restClient *resty.Client
}

func NewClient(baseURL, apiKey string) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("no base url provided")
	} else if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	if apiKey == "" {
		return nil, fmt.Errorf("no api key provided")
	}

	restClient := resty.New()
	restClient.SetHostURL(baseURL)
	restClient.SetHeader("X-Api-Key", apiKey)
	restClient.SetHeader("Accept", "application/json")

	return &Client{restClient: restClient}, nil
}

func (c Client) checkResponse(response *resty.Response, expectedStati ...int) error {
	switch response.StatusCode() {
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 500:
		return ErrInternalServerError
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
