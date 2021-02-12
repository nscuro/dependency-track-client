package dtrack

import "context"

const (
	PolicyViolationLicense     = "LICENSE"
	PolicyViolationOperational = "OPERATIONAL"
	PolicyViolationSecurity    = "SECURITY"
)

type PolicyViolation struct {
	Component *Component `json:"component"`
	Project   *Project   `json:"project"`
	Text      string     `json:"text"`
	Timestamp int64      `json:"timestamp"`
	Type      string     `json:"type"`
	UUID      string     `json:"uuid"`
}

type PolicyViolationService interface {
	GetForComponent(ctx context.Context, componentUUID string) ([]PolicyViolation, error)
	GetForProject(ctx context.Context, projectUUID string) ([]PolicyViolation, error)
}

type policyViolationSericeImpl struct {
	client *Client
}

func (p policyViolationSericeImpl) GetForComponent(ctx context.Context, componentUUID string) ([]PolicyViolation, error) {
	violations := make([]PolicyViolation, 0)

	req := p.client.restClient.R().
		SetContext(ctx).
		SetPathParam("uuid", componentUUID).
		SetResult([]PolicyViolation{})

	err := p.client.getPaginatedResponse(req, "/api/v1/violation/component/{uuid}", func(res interface{}) (int, error) {
		violationsOnPage, ok := res.(*[]PolicyViolation)
		if !ok {
			return -1, ErrInvalidResponseType
		}
		violations = append(violations, *violationsOnPage...)
		return len(violations), nil
	})
	if err != nil {
		return nil, err
	}

	return violations, nil
}

func (p policyViolationSericeImpl) GetForProject(ctx context.Context, projectUUID string) ([]PolicyViolation, error) {
	violations := make([]PolicyViolation, 0)

	req := p.client.restClient.R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"uuid": projectUUID,
		}).
		SetResult([]PolicyViolation{})

	err := p.client.getPaginatedResponse(req, "/api/v1/violation/project/{uuid}", func(res interface{}) (int, error) {
		violationsOnPage, ok := res.(*[]PolicyViolation)
		if !ok {
			return -1, ErrInvalidResponseType
		}
		violations = append(violations, *violationsOnPage...)
		return len(violations), nil
	})
	if err != nil {
		return nil, err
	}

	return violations, nil
}
