package dtrack

type License struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Text      string `json:"licenseText"`
	LicenseID string `json:"licenseId"`
}
