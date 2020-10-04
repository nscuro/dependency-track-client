package audit

import (
	"errors"
	"fmt"

	"github.com/nscuro/dependency-track-client"
)

var (
	ErrMaxRiskScoreExceeded      = errors.New("max risk score exceeded")
	ErrMaxSeverityExceeded       = errors.New("max severity exceeded")
	ErrSeverityThresholdExceeded = errors.New("severity threshold exceeded")

	severities = map[string]int{
		dtrack.UnassignedSeverity: 0,
		dtrack.InfoSeverity:       1,
		dtrack.LowSeverity:        2,
		dtrack.MediumSeverity:     3,
		dtrack.HighSeverity:       4,
		dtrack.CriticalSeverity:   5,
	}
)

// TODO: Add functionality to read from YAML file
type QualityGate struct {
	MaxRiskScore       int64          `yaml:"max-risk-score"`
	MaxSeverity        string         `yaml:"max-severity"`
	SeverityThresholds map[string]int `yaml:"severity-thresholds"`
}

// TODO: Run through all quality gates, don't abort after the first failure
// TODO: Return a structure with detailed info about the failures
func (q QualityGate) Evaluate(metrics dtrack.ProjectMetrics, findings []dtrack.Finding) error {
	if q.MaxRiskScore > 0 && metrics.InheritedRiskScore > q.MaxRiskScore {
		return ErrMaxRiskScoreExceeded
	}

	if q.MaxSeverity != "" {
		maxSeverityValue, ok := severities[q.MaxSeverity]
		if !ok {
			return fmt.Errorf("invalid severity \"%s\"", q.MaxSeverity)
		}

		for _, finding := range findings {
			severityValue, ok := severities[finding.Vulnerability.Severity]
			if !ok {
				return fmt.Errorf("invalid severity \"%s\"", finding.Vulnerability.Severity)
			}

			if severityValue > maxSeverityValue {
				return ErrMaxSeverityExceeded
			}
		}
	}

	if len(q.SeverityThresholds) > 0 {
		// TODO: Verify
		return ErrSeverityThresholdExceeded
	}

	return nil
}
