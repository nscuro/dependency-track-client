package dtrack

import (
	"context"
	"fmt"
	"time"
)

type ProjectMetrics struct {
	Components                       int     `json:"components"`
	Critical                         int     `json:"critical"`
	FindingsAudited                  int     `json:"findingsAudited"`
	FindingsTotal                    int     `json:"findingsTotal"`
	FindingsUnaudited                int     `json:"findingsUnaudited"`
	High                             int     `json:"high"`
	InheritedRiskScore               float32 `json:"inheritedRiskScore"`
	Low                              int     `json:"low"`
	Medium                           int     `json:"medium"`
	PolicyViolationsFail             int     `json:"policyViolationsFail"`
	PolicyViolationsInfo             int     `json:"policyViolationsInfo"`
	PolicyViolationsLicenseTotal     int     `json:"policyViolationsLicenseTotal"`
	PolicyViolationsOperationalTotal int     `json:"policyViolationsOperationalTotal"`
	PolicyViolationsSecurityTotal    int     `json:"policyViolationsSecurityTotal"`
	PolicyViolationsTotal            int     `json:"policyViolationsTotal"`
	PolicyViolationsWarn             int     `json:"policyViolationsWarn"`
	Suppressed                       int     `json:"suppressed"`
	Unassigned                       int     `json:"unassigned"`
	VulnerableComponents             int     `json:"vulnerableComponents"`
}

func (pm ProjectMetrics) GetSeverityCount(severity string) (count int, err error) {
	switch severity {
	case SeverityCritical:
		count = pm.Critical
	case SeverityHigh:
		count = pm.High
	case SeverityMedium:
		count = pm.Medium
	case SeverityLow:
		count = pm.Low
	case SeverityUnassigned:
		count = pm.Unassigned
	default:
		err = fmt.Errorf("cannot determine count for severity %s", severity)
	}
	return
}

func (pm ProjectMetrics) GetViolationCount(violationType string) (count int, err error) {
	switch violationType {
	case PolicyViolationLicense:
		count = pm.PolicyViolationsLicenseTotal
	case PolicyViolationOperational:
		count = pm.PolicyViolationsOperationalTotal
	case PolicyViolationSecurity:
		count = pm.PolicyViolationsSecurityTotal
	default:
		err = fmt.Errorf("cannot determine count for violation type %s", violationType)
	}
	return
}

type ProjectMetricsService interface {
	GetCurrent(ctx context.Context, projectUUID string) (*ProjectMetrics, error)
	GetForDays(ctx context.Context, projectUUID string, days int) (*ProjectMetrics, error)
	GetSince(ctx context.Context, projectUUID string, date time.Time) (*ProjectMetrics, error)
}

type projectMetricsServiceImpl struct {
	client *Client
}

func (p projectMetricsServiceImpl) GetCurrent(ctx context.Context, projectUUID string) (*ProjectMetrics, error) {
	res, err := p.client.restClient.R().
		SetContext(ctx).
		SetPathParam("uuid", projectUUID).
		SetResult(&ProjectMetrics{}).
		Get("/api/v1/metrics/project/{uuid}/current")
	if err != nil {
		return nil, err
	}

	if err = p.client.checkResponseStatus(res, 200); err != nil {
		return nil, err
	}

	metrics, ok := res.Result().(*ProjectMetrics)
	if !ok {
		return nil, ErrInvalidResponseType
	}

	return metrics, nil
}

func (p projectMetricsServiceImpl) GetForDays(ctx context.Context, projectUUID string, days int) (*ProjectMetrics, error) {
	panic("implement me")
}

func (p projectMetricsServiceImpl) GetSince(ctx context.Context, projectUUID string, date time.Time) (*ProjectMetrics, error) {
	panic("implement me")
}
