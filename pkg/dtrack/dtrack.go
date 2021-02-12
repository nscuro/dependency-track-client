package dtrack

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	defaultPageSize = 50

	headerAPIKey     = "X-Api-Key"
	headerTotalCount = "X-Total-Count"
)

var (
	ErrConflict            = errors.New("conflict")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidResponseType = errors.New("invalid response type")

	ErrMissingTotalCountHeader = errors.New("response does not contain " + headerTotalCount + " header")
)

type Client struct {
	restClient *resty.Client

	About                   AboutService
	Analysis                AnalysisService
	BOM                     BOMService
	Component               ComponentService
	Finding                 FindingService
	License                 LicenseService
	LicenseGroup            LicenseGroupService
	PolicyViolation         PolicyViolationService
	PolicyViolationAnalysis PolicyViolationAnalysisService
	Project                 ProjectService
	ProjectMetrics          ProjectMetricsService
	Vulnerability           VulnerabilityService
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
	restClient.SetHeader(headerAPIKey, apiKey)
	restClient.SetHeader("Accept", "application/json")

	client := &Client{restClient: restClient}
	client.About = &aboutServiceImpl{client: client}
	client.Analysis = &analysisServiceImpl{client: client}
	client.BOM = &bomServiceImpl{client: client}
	client.Component = &componentServiceImpl{client: client}
	client.Finding = &findingServiceImpl{client: client}
	client.License = &licenseServiceImpl{client: client}
	client.LicenseGroup = &licenseGroupServiceImpl{client: client}
	client.PolicyViolation = &policyViolationSericeImpl{client: client}
	client.PolicyViolationAnalysis = &policyViolationAnalysisServiceImpl{client: client}
	client.Project = &projectServiceImpl{client: client}
	client.ProjectMetrics = &projectMetricsServiceImpl{client: client}
	client.Vulnerability = &vulnerabilityServiceImpl{client: client}
	return client, nil
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

func (c Client) getPaginatedResponse(req *resty.Request, url string, handler func(interface{}) (int, error)) error {
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

		totalCountStr := res.Header().Get(headerTotalCount)
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
