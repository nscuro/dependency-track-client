package dtrack

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAbout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/version", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `
		{
  			"timestamp": "2020-09-29T22:02:23Z",
  			"version": "4.0.0-SNAPSHOT",
  			"uuid": "c35ce882-398d-46ed-8c36-6148bf73f941",
  			"systemUuid": "f2eee6f6-a161-418a-baf5-b57a2e30de82",
			"application": "Dependency-Track",
  			"framework": {
    			"timestamp": "2020-07-20T15:56:44Z",
    			"version": "1.8.0-SNAPSHOT",
    			"uuid": "beee8786-ca9c-473a-b7a5-efcc95e8c469",
    			"name": "Alpine"
  			}
		}`)
	}))
	defer mockServer.Close()

	client, _ := NewClient(mockServer.URL, "apiKey")

	about, err := client.About.Get(context.Background())
	assert.NoError(t, err)

	assert.NotNil(t, about)
	assert.Equal(t, "2020-09-29T22:02:23Z", about.Timestamp)
	assert.Equal(t, "4.0.0-SNAPSHOT", about.Version)
	assert.Equal(t, "c35ce882-398d-46ed-8c36-6148bf73f941", about.UUID)
	assert.Equal(t, "f2eee6f6-a161-418a-baf5-b57a2e30de82", about.SystemUUID)
	assert.Equal(t, "Dependency-Track", about.Application)

	assert.Equal(t, "2020-07-20T15:56:44Z", about.Framework.Timestamp)
	assert.Equal(t, "1.8.0-SNAPSHOT", about.Framework.Version)
	assert.Equal(t, "beee8786-ca9c-473a-b7a5-efcc95e8c469", about.Framework.UUID)
	assert.Equal(t, "Alpine", about.Framework.Name)
}
