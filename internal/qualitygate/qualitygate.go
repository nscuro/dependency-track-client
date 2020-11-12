package qualitygate

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/nscuro/dependency-track-client"
)

var (
	severities = []string{
		dtrack.UnassignedSeverity,
		dtrack.InfoSeverity,
		dtrack.LowSeverity,
		dtrack.MediumSeverity,
		dtrack.HighSeverity,
		dtrack.CriticalSeverity,
	}
	violationTypes = []string{
		dtrack.LicensePolicyViolation,
		dtrack.OperationalPolicyViolation,
		dtrack.SecurityPolicyViolation,
	}
)

type Gate struct {
	MaxRiskScore        float32        `json:"max_risk_score" yaml:"max-risk-score"`
	MaxSeverity         string         `json:"max_severity" yaml:"max-severity"`
	SeverityThresholds  map[string]int `json:"severity_thresholds" yaml:"severity-thresholds"`
	ViolationThresholds map[string]int `json:"violation_thresholds" yaml:"violation-thresholds"`
}

func LoadGateFromFile(filePath string) (*Gate, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if strings.HasSuffix(filePath, ".json") {
		return loadGateFromJSONFile(file)
	} else {
		return loadGateFromYAMLFile(file)
	}
}

func loadGateFromJSONFile(file *os.File) (*Gate, error) {
	var gate Gate
	if err := json.NewDecoder(file).Decode(&gate); err != nil {
		return nil, err
	}
	return &gate, nil
}

func loadGateFromYAMLFile(file *os.File) (*Gate, error) {
	var gate Gate
	if err := yaml.NewDecoder(file).Decode(&gate); err != nil {
		return nil, err
	}
	return &gate, nil
}

type Evaluator struct {
	dtrackClient *dtrack.Client
}

func NewEvaluator(dtrackClient *dtrack.Client) *Evaluator {
	return &Evaluator{
		dtrackClient: dtrackClient,
	}
}

// Evaluate evaluates a given Gate for a given project
func (e Evaluator) Evaluate(projectUUID string, gate *Gate) error {
	projectMetrics, err := e.dtrackClient.GetCurrentProjectMetrics(projectUUID)
	if err != nil {
		return fmt.Errorf("failed to retrieve project metrics: %w", err)
	}

	e.evaluateMaxRiskScore(projectMetrics, gate)
	e.evaluateMaxSeverity(projectMetrics, gate)
	e.evaluateSeverityThresholds(projectMetrics, gate)
	e.evaluateViolationThresholds(projectMetrics, gate)

	log.Println("quality gate passed")
	return nil
}

func (e Evaluator) evaluateMaxRiskScore(metrics *dtrack.ProjectMetrics, gate *Gate) error {
	if gate.MaxRiskScore <= -1 {
		return nil
	}

	log.Println("evaluating max risk score")
	if metrics.InheritedRiskScore > gate.MaxRiskScore {
		return fmt.Errorf("expected risk score to be <= %.2f, but was %.2f", gate.MaxRiskScore, metrics.InheritedRiskScore)
	}
	return nil
}

func (e Evaluator) evaluateMaxSeverity(metrics *dtrack.ProjectMetrics, gate *Gate) error {
	if gate.MaxSeverity == "" {
		return nil
	}

	log.Println("evaluating max severity")
	switch gate.MaxSeverity {
	case dtrack.LowSeverity:
		if metrics.Low > 0 {
			return maxSeverityExceededError(gate.MaxSeverity, dtrack.LowSeverity)
		}
		fallthrough
	case dtrack.MediumSeverity:
		if metrics.Medium > 0 {
			return maxSeverityExceededError(gate.MaxSeverity, dtrack.MediumSeverity)
		}
		fallthrough
	case dtrack.HighSeverity:
		if metrics.High > 0 {
			return maxSeverityExceededError(gate.MaxSeverity, dtrack.HighSeverity)
		}
		fallthrough
	case dtrack.CriticalSeverity:
		if metrics.Critical > 0 {
			return maxSeverityExceededError(gate.MaxSeverity, dtrack.CriticalSeverity)
		}
	default:
		return fmt.Errorf("invalid severity \"%s\"", gate.MaxSeverity)
	}
	return nil
}

func (e Evaluator) evaluateSeverityThresholds(metrics *dtrack.ProjectMetrics, gate *Gate) error {
	if len(gate.SeverityThresholds) == 0 {
		return nil
	}

	log.Println("evaluating severity thresholds")
	for _, severity := range severities {
		threshold, ok := gate.SeverityThresholds[severity]
		if !ok {
			continue
		}

		count, err := metrics.GetSeverityCount(severity)
		if err != nil {
			continue
		}

		if count > threshold {
			return fmt.Errorf("threshold for severity %s exceeded: allowed=%d actual=%d", severity, threshold, count)
		}
	}
	return nil
}

func (e Evaluator) evaluateViolationThresholds(metrics *dtrack.ProjectMetrics, gate *Gate) error {
	if len(gate.ViolationThresholds) == 0 {
		return nil
	}

	log.Println("evaluating violation thresholds")
	for _, violationType := range violationTypes {
		threshold, ok := gate.ViolationThresholds[violationType]
		if !ok {
			continue
		}

		count, err := metrics.GetViolationCount(violationType)
		if err != nil {
			continue
		}

		if count > threshold {
			return fmt.Errorf("threshold for violation type %s exceeded: allowed=%d actual=%d", violationType, threshold, count)
		}
	}
	return nil
}

func maxSeverityExceededError(max, actual string) error {
	return fmt.Errorf("maximum severity exceeded: allowed=%s actual=%s", max, actual)
}
