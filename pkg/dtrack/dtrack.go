package dtrack

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	defaultPageSize  = 50
	totalCountHeader = "X-Total-Count"
)

var (
	ErrConflict            = errors.New("conflict")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidResponseType = errors.New("invalid response type")

	ErrMissingTotalCountHeader = errors.New("response does not contain " + totalCountHeader + " header")
)

type Client struct {
	restClient *resty.Client
}

type ClientOptions struct {
	HttpClient *http.Client
}

func NewClient(baseURL, apiKey string, options ...ClientOptions) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("no base url provided")
	} else if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	if apiKey == "" {
		return nil, fmt.Errorf("no api key provided")
	}

	var restClient *resty.Client
	if len(options) > 0 && options[0].HttpClient != nil {
		restClient = resty.NewWithClient(options[0].HttpClient)
	} else {
		restClient = resty.New()
	}
	restClient.SetHostURL(baseURL)
	restClient.SetHeader("X-Api-Key", apiKey)
	restClient.SetHeader("Accept", "application/json")

	return &Client{restClient: restClient}, nil
}

func (c Client) checkResponseStatus(response *resty.Response, expectedStati ...int) error {
	if len(expectedStati) == 0 {
		return fmt.Errorf("no expected status code provided")
	}

	// Handle common API errors
	switch response.StatusCode() {
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 409:
		return ErrConflict
	case 500:
		return ErrInternalServerError
	}

	if len(expectedStati) > 0 {
		for _, status := range expectedStati {
			if response.StatusCode() == status {
				return nil
			}
		}
		return fmt.Errorf("expected response status to be any of %v, but was %d", expectedStati, response.StatusCode())
	}

	return nil
}

type paginatedResultHandler func(interface{}) (int, error)

func (c Client) getPaginatedResponse(req *resty.Request, url string, handler paginatedResultHandler) error {
	page := 1
	hasMorePages := true

	for hasMorePages {
		req.SetQueryParams(map[string]string{
			"pageSize":   strconv.Itoa(defaultPageSize),
			"pageNumber": strconv.Itoa(page),
		})
		res, err := req.Execute(resty.MethodGet, url)
		if err != nil {
			return err
		}

		if err = c.checkResponseStatus(res, 200); err != nil {
			return err
		}

		totalCountStr := res.Header().Get(totalCountHeader)
		if totalCountStr == "" {
			return ErrMissingTotalCountHeader
		}

		totalCount, err := strconv.Atoi(totalCountStr)
		if err != nil {
			return err
		}

		itemCount, err := handler(res.Result())
		if err != nil {
			return err
		}

		hasMorePages = itemCount < totalCount
		page++
	}

	return nil
}