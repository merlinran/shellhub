package telemetry

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendTelemetry(t *testing.T) {
	var req *http.Request
	var body []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		req = r
		body, err = io.ReadAll(r.Body)
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}))

	defer srv.Close()

	properties := map[string]interface{}{
		"value": 1.0,
	}

	telemetry := New("id", srv.URL)
	telemetry.SendTelemetry("test", properties)

	assert.NotNil(t, req)

	track := Track{}
	err := json.Unmarshal(body, &track)
	assert.NoError(t, err)

	assert.Equal(t, "/track", req.URL.Path)
	assert.Equal(t, "id", track.InstanceID)
	assert.Equal(t, "test", track.Event)
	assert.Equal(t, properties, track.Properties)
}
