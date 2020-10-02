package audit

import (
	"errors"
	"fmt"

	"github.com/nscuro/dependency-track-client/pkg/dtrack"
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

type QualityGate struct {
	MaxRiskScore       int64          `yaml:"max-risk-score"`
	MaxSeverity        string         `yaml:"max-severity"`
	SeverityThresholds map[string]int `yaml:"severity-thresholds"`
}

func (q QualityGate) Evaluate(findings []dtrack.Finding) error {
	if q.MaxRiskScore > 0 && calculateInheritedRiskScore(findings) > q.MaxRiskScore {
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

func calculateInheritedRiskScore(findings []dtrack.Finding) int64 {
	critical, high, medium, low, info, unassigned := 0, 0, 0, 0, 0, 0

	for _, finding := range findings {
		switch finding.Vulnerability.Severity {
		case dtrack.CriticalSeverity:
			critical++
		case dtrack.HighSeverity:
			high++
		case dtrack.MediumSeverity:
			medium++
		case dtrack.LowSeverity:
			low++
		case dtrack.InfoSeverity:
			info++
		case dtrack.UnassignedSeverity:
			unassigned++
		}
	}

	// https://github.com/DependencyTrack/dependency-track/blob/master/src/main/java/org/dependencytrack/metrics/Metrics.java#L32
	return int64((critical * 10) + (high * 5) + (medium * 3) + (low * 1) + (unassigned * 5))
}
