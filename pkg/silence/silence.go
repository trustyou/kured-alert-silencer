package silence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/smithy-go/ptr"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/client"
	"github.com/prometheus/alertmanager/api/v2/client/silence"
	"github.com/prometheus/alertmanager/api/v2/models"
)

// generate models.Matcher form JSON string with format `[{"name": "instance", "value": "{{.NodeName}}"}, {"name": "alertname", "value": "node_reboot"}]`
func generateMatchers(matchersJSON string, nodeName string) ([]*models.Matcher, error) {
	tmpl, err := template.New("matchers").Parse(matchersJSON)
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"NodeName": nodeName,
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return nil, err
	}

	var matchers []*models.Matcher
	log.Debugf("rendered matchers: %s", tpl.String())
	err = json.Unmarshal(tpl.Bytes(), &matchers)
	if err != nil {
		return nil, err
	}

	// check that matchers contain required fields
	for _, matcher := range matchers {
		if matcher.Name == nil || matcher.Value == nil || matcher.IsRegex == nil {
			return nil, fmt.Errorf("matcher is missing required fields")
		}
	}

	return matchers, nil
}

// create client.AlertmanagerAPI from alertmanagerURL
func NewAlertmanagerClient(alertmanagerURL string) (*client.AlertmanagerAPI, error) {
	u, err := url.Parse(alertmanagerURL)
	if err != nil {
		return nil, err
	}

	scheme := u.Scheme
	host := u.Host

	log.Infof("Alertmanager scheme: %s", scheme)
	log.Infof("Alertmanager host: %s", host)
	config := client.DefaultTransportConfig().WithSchemes([]string{scheme}).WithHost(host)

	alertmanager := client.NewHTTPClientWithConfig(nil, config)
	return alertmanager, nil
}

// Get silences from Alertmanager that match the given matcher
func silenceExists(alertmanager *client.AlertmanagerAPI, matcher *models.Matcher) (bool, error) {
	getSilencesParams := silence.NewGetSilencesParams()
	matchersStr := []string{fmt.Sprintf("%s=%s", *matcher.Name, *matcher.Value)}

	getSilencesResp, err := alertmanager.Silence.GetSilences(getSilencesParams.WithFilter(matchersStr))
	if err != nil {
		return true, err
	}

	return len(getSilencesResp.Payload) > 0, nil
}

// SilenceAlerts silences alerts in Alertmanager
func SilenceAlerts(alertmanager *client.AlertmanagerAPI, matchersJSON string, nodeName string, alertTTL string) error {
	startsAt := time.Now()
	alertTTLtime, err := time.ParseDuration(alertTTL)
	if err != nil {
		return err
	}
	endsAt := startsAt.Add(alertTTLtime)

	matchers, err := generateMatchers(matchersJSON, nodeName)
	if err != nil {
		return err
	}
	log.Infof("Silencing alerts with matchers: %v", len(matchers))

	for _, matcher := range matchers {
		log.Debugf(
			"matcher: %sIsRegex: %t, Name: %s, Value: %s",
			func(b *bool) interface{} {
				if b == nil {
					return ""
				} else {
					return fmt.Sprintf("IsEqual: %t, ", *b)
				}
			}(matcher.IsEqual),
			*matcher.IsRegex,
			*matcher.Name,
			*matcher.Value,
		)

		exists, err := silenceExists(alertmanager, matcher)
		if err != nil {
			return err
		}

		if exists {
			log.Debugf("Silence already exists for matcher: %s=%s\n", *matcher.Name, *matcher.Value)
			continue
		}

		postSilenceParams := silence.NewPostSilencesParams().WithSilence(
			&models.PostableSilence{
				Silence: models.Silence{
					Matchers:  []*models.Matcher{matcher},
					StartsAt:  (*strfmt.DateTime)(&startsAt),
					EndsAt:    (*strfmt.DateTime)(&endsAt),
					CreatedBy: ptr.String("kured-alert-silencer"),
					Comment:   ptr.String("Silencing during node reboot"),
				},
			},
		)

		_, err = alertmanager.Silence.PostSilences(postSilenceParams)
		if err != nil {
			return err
		}

		log.Debugf("Silence created for matcher: %s=%s\n", *matcher.Name, *matcher.Value)
		log.Info("Silence created successfully")
	}
	return nil
}
