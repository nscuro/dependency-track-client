package dtrack

type Analysis struct {
	State      string `json:"state"`
	Suppressed bool   `json:"isSuppressed"`
}
