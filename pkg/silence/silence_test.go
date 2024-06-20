package silence

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/stretchr/testify/assert"
)

// Mock server for Alertmanager API
func mockAlertmanagerServer(existingSilences []*models.GettableSilence) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/v2/silences", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingSilences)
		} else if r.Method == "POST" {
			w.WriteHeader(http.StatusOK)
		}
	})
	return httptest.NewServer(handler)
}

func TestGenerateMatchers(t *testing.T) {
	validJSON := `[{"name": "instance", "value": "{{.NodeName}}", "isRegex": false}, {"name": "alertname", "value": "node_reboot", "isRegex": false}]`
	invalidJSON := `[{name: "instance", "value": "{{.NodeName}}", "isRegex": false}]`
	missingFieldsJSON := `[{"name": "instance", "value": "{{.NodeName}}"}`

	tests := []struct {
		name      string
		jsonInput string
		nodeName  string
		expectErr bool
	}{
		{"Valid JSON", validJSON, "node1", false},
		{"Invalid JSON", invalidJSON, "node1", true},
		{"Missing Fields JSON", missingFieldsJSON, "node1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchers, err := generateMatchers(tt.jsonInput, tt.nodeName)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "instance", *matchers[0].Name)
				assert.Equal(t, "node1", *matchers[0].Value)
				assert.Equal(t, "alertname", *matchers[1].Name)
				assert.Equal(t, "node_reboot", *matchers[1].Value)
			}
		})
	}
}

func TestNewAlertmanagerClient(t *testing.T) {
	tests := []struct {
		name            string
		alertmanagerURL string
		expectErr       bool
	}{
		{"Valid URL", "http://localhost:9093", false},
		{"Invalid URL", ":", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewAlertmanagerClient(tt.alertmanagerURL)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestSilenceExists(t *testing.T) {
	existingSilences := []*models.GettableSilence{
		{
			Silence: models.Silence{
				Matchers: []*models.Matcher{
					{Name: ptr.String("instance"), Value: ptr.String("node1")},
				},
				StartsAt:  (*strfmt.DateTime)(ptr.Time(time.Now().Add(-1 * time.Hour))),
				EndsAt:    (*strfmt.DateTime)(ptr.Time(time.Now().Add(1 * time.Hour))),
				CreatedBy: ptr.String("kured-alert-silencer"),
				Comment:   ptr.String("Silencing during node reboot"),
			},
		},
	}

	tests := []struct {
		name             string
		existingSilences []*models.GettableSilence
		matcher          *models.Matcher
		expectedExists   bool
	}{
		{
			name:             "No Existing Silences",
			existingSilences: []*models.GettableSilence{},
			matcher:          &models.Matcher{Name: ptr.String("instance"), Value: ptr.String("node1")},
			expectedExists:   false,
		},
		{
			name:             "Existing Silences",
			existingSilences: existingSilences,
			matcher:          &models.Matcher{Name: ptr.String("instance"), Value: ptr.String("node1")},
			expectedExists:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mockAlertmanagerServer(tt.existingSilences)
			defer server.Close()

			alertmanager, err := NewAlertmanagerClient(server.URL)
			assert.NoError(t, err)

			exists, err := silenceExists(alertmanager, tt.matcher)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedExists, exists)
		})
	}
}

func TestSilenceAlerts(t *testing.T) {
	existingSilences := []*models.GettableSilence{
		{
			Silence: models.Silence{
				Matchers: []*models.Matcher{
					{Name: ptr.String("instance"), Value: ptr.String("node1")},
				},
				StartsAt:  (*strfmt.DateTime)(ptr.Time(time.Now().Add(-1 * time.Hour))),
				EndsAt:    (*strfmt.DateTime)(ptr.Time(time.Now().Add(1 * time.Hour))),
				CreatedBy: ptr.String("kured-alert-silencer"),
				Comment:   ptr.String("Silencing during node reboot"),
			},
		},
	}

	server := mockAlertmanagerServer(existingSilences)
	defer server.Close()

	alertmanager, err := NewAlertmanagerClient(server.URL)
	assert.NoError(t, err)

	validMatchersJSON := `[{"name": "instance", "value": "{{.NodeName}}", "isRegex": false}, {"name": "alertname", "value": "node_reboot", "isRegex": false}]`
	invalidMatchersJSON := `[{name: "instance", "value": "{{.NodeName}}", "isRegex": false}]`

	tests := []struct {
		name         string
		matchersJSON string
		nodeName     string
		alertTTL     string
		expectErr    bool
	}{
		{"Valid Silence", validMatchersJSON, "node1", "1h", false},
		{"Invalid Matchers JSON", invalidMatchersJSON, "node1", "1h", true},
		{"Invalid Alert TTL", validMatchersJSON, "node1", "invalid_duration", true},
		{"Existing Silence", validMatchersJSON, "node1", "1h", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SilenceAlerts(alertmanager, tt.matchersJSON, tt.nodeName, tt.alertTTL)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
