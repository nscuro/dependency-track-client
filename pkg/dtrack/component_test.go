package dtrack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetComponent(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/component/4d5cd8df-cff7-4212-a038-91ae4ab79396", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `
		{
			"uuid": "4d5cd8df-cff7-4212-a038-91ae4ab79396",
			"group": "apache",
			"name": "axis",
			"version": "1.4",
			"md5": "03dcfdd88502505cc5a805a128bfdd8d",
			"sha1": "94a9ce681a42d0352b3ad22659f67835e560d107",
			"sha256": "05aebb421d0615875b4bf03497e041fe861bf0556c3045d8dda47e29241ffdd3",
			"purl": "pkg:maven/apache/axis@1.4",
			"isInternal": false
		}
		`)
	}))
	defer mockServer.Close()

	client, _ := NewClient(mockServer.URL, "apiKey")

	component, err := client.GetComponent("4d5cd8df-cff7-4212-a038-91ae4ab79396")
	assert.NoError(t, err)
	assert.NotNil(t, component)
	assert.Equal(t, "4d5cd8df-cff7-4212-a038-91ae4ab79396", component.UUID)
	assert.Equal(t, "axis", component.Name)
	assert.Equal(t, "apache", component.Group)
	assert.Equal(t, "1.4", component.Version)
	assert.Equal(t, "pkg:maven/apache/axis@1.4", component.PackageURL)
	assert.Equal(t, false, component.Internal)
}
