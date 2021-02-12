package dtrack

import "context"

type PolicyViolationAnalysis struct {
	Comments   []PolicyViolationAnalysisComment `json:"analysisComments"`
	State      string                           `json:"analysisState"`
	Suppressed bool                             `json:"isSuppressed"`
}

type PolicyViolationAnalysisComment struct {
	Comment   string `json:"comment"`
	Commenter string `json:"commenter"`
	Timestamp string `json:"timestamp"`
}

type PolicyViolationAnalysisRequest struct {
	Comment       string `json:"comment,omitempty"`
	ComponentUUID string `json:"component"`
	ViolationUUID string `json:"policyViolation"`
	State         string `json:"analysisState,omitempty"`
	Suppressed    bool   `json:"isSuppressed"`
}

type PolicyViolationAnalysisService interface {
	Create(ctx context.Context, req PolicyViolationAnalysisRequest) (*PolicyViolationAnalysis, error)
	Get(ctx context.Context, componentUUID, violationUUID string) (*PolicyViolationAnalysis, error)
}

type policyViolationAnalysisServiceImpl struct {
	client *Client
}

func (p policyViolationAnalysisServiceImpl) Create(ctx context.Context, req PolicyViolationAnalysisRequest) (*PolicyViolationAnalysis, error) {
	panic("implement me")
}

func (p policyViolationAnalysisServiceImpl) Get(ctx context.Context, componentUUID, violationUUID string) (*PolicyViolationAnalysis, error) {
	panic("implement me")
}
