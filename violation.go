package dtrack

type PolicyViolation struct {
	Component *Component `json:"component"`
	Project   *Project   `json:"project"`
	Text      string     `json:"text"`
	Timestamp int64      `json:"timestamp"`
	Type      string     `json:"type"`
	UUID      string     `json:"uuid"`
}

func (c Client) GetPolicyViolationsForComponent(componentUUID string) ([]PolicyViolation, error) {
	return nil, nil
}

func (c Client) GetPolicyViolationsForProject(projectUUID string) ([]PolicyViolation, error) {
	violations := make([]PolicyViolation, 0)

	req := c.restClient.R().
		SetPathParams(map[string]string{
			"uuid": projectUUID,
		}).
		SetResult([]PolicyViolation{})

	err := c.getPaginatedResponse(req, "/api/v1/violation/project/{uuid}", func(result interface{}) (int, error) {
		violationsOnPage, ok := result.(*[]PolicyViolation)
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

type ViolationAnalysis struct {
	Comments   []ViolationAnalysisComment `json:"analysisComments"`
	State      string                     `json:"analysisState"`
	Suppressed bool                       `json:"isSuppressed"`
}

type ViolationAnalysisComment struct {
	Comment   string `json:"comment"`
	Commenter string `json:"commenter"`
	Timestamp string `json:"timestamp"`
}

type ViolationAnalysisRequest struct {
	Comment       string `json:"comment,omitempty"`
	ComponentUUID string `json:"component"`
	ViolationUUID string `json:"policyViolation"`
	State         string `json:"analysisState,omitempty"`
	Suppressed    bool   `json:"isSuppressed"`
}

func (c Client) GetViolationAnalysis(componentUUID, violationUUID string) (*ViolationAnalysis, error) {
	return nil, nil
}

func (c Client) RecordViolationAnalysis(req ViolationAnalysisRequest) (*ViolationAnalysis, error) {
	return nil, nil
}
