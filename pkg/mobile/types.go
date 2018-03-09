package mobile

import (
	"encoding/json"
	"fmt"
)

// swagger:ignore
type AppConfig struct {
	DBConnectionString string
}

// ClientMetric struct is what the client payload should be parsed into
// Need to figure out how to structure this
//
// swagger:model
type Metric struct {
	// required: true
	// example: 1520853523661
	ClientTimestamp json.Number `json:"timestamp,omitempty"`
	// required: true
	// example: some-unique-client-id
	ClientId string `json:"clientId"`
	// required: true
	Data *MetricData `json:"data,omitempty"`
	// required: true
	// example: init
	EventType string `json:"type"`
}

// swagger:model
type MetricData struct {
	// required: true
	App *AppMetric `json:"app,omitempty"`
	// required: true
	Device *DeviceMetric `json:"device,omitempty"`
	// required: true
	Security *SecurityMetrics `json:"security,omitempty"`
}

// swagger:model
type AppMetric struct {
	// required: true
	// example: com.example.myapp
	ID string `json:"appId"`
	// required: true
	// example: 1.0.0
	SDKVersion string `json:"sdkVersion"`
	// required: true
	// example: 2.1.0
	AppVersion string `json:"appVersion"`
}

// swagger:model
type DeviceMetric struct {
	// required: true
	// example: android
	Platform string `json:"platform"`
	// required: true
	// example: 27
	PlatformVersion string `json:"platformVersion"`
}

type SecurityMetrics []SecurityMetric

// swagger:model
type SecurityMetric struct {
	// required: true
	// example: com.example.DeveloperMode
	Id *string `json:"id,omitempty"`
	// required: true
	// example: Developer Mode
	Name *string `json:"name,omitempty"`
	// required: true
	Passed *bool `json:"passed,omitempty"`
}

const clientIdMaxLength = 128
const eventTypeMaxLength = 128
const securityMetricsMaxLength = 30

const missingClientIdError = "missing clientId in payload"
const missingEventTypeError = "missing type in payload"
const unknownTypeError = "payload type unknown"
const missingAppError = "missing data.app in init-type payload"
const missingDeviceError = "missing data.device in init-type payload"
const missingSecurityError = "missing data.security in security-type payload"
const missingDataError = "missing metrics data in payload"
const initMetricsIncompleteError = "data.app and data.security must both be present simultaneously"

const invalidTimestampError = "timestamp must be a valid number"
const securityMetricsEmptyError = "data.security cannot be empty"
const securityMetricMissingIdError = "invalid element in data.security at position %v, id must be included"
const securityMetricMissingNameError = "invalid element in data.security at position %v, name must be included"
const securityMetricMissingPassedError = "invalid element in data.security at position %v, passed must be included"

var clientIdLengthError = fmt.Sprintf("clientId exceeded maximum length of %v", clientIdMaxLength)

var eventTypeLengthError = fmt.Sprintf("type exceeded maximum length of %v", eventTypeMaxLength)

var securityMetricsLengthError = fmt.Sprintf("maximum length of data.security %v", securityMetricsMaxLength)

func (m *Metric) Validate() (valid bool, reason string) {
	if m.ClientId == "" {
		return false, missingClientIdError
	}

	if len(m.ClientId) > clientIdMaxLength {
		return false, clientIdLengthError
	}

	if m.EventType == "" {
		return false, missingEventTypeError
	}

	if len(m.EventType) > eventTypeMaxLength {
		return false, eventTypeLengthError
	}

	if m.ClientTimestamp != "" {
		if _, err := m.ClientTimestamp.Int64(); err != nil {
			return false, invalidTimestampError
		}
	}

	// check if data field was missing or empty object
	if m.Data == nil || (MetricData{}) == *m.Data {
		return false, missingDataError
	}

	switch m.EventType {
	case "init":
		return validateInitMetric(m.Data)
	case "security":
		return validateSecurityMetric(m.Data)
	default:
		return false, unknownTypeError
	}
}

func validateInitMetric(data *MetricData) (valid bool, reason string) {
	if data.App == nil {
		return false, missingAppError
	}
	if data.Device == nil {
		return false, missingDeviceError
	}
	return true, ""
}

func validateSecurityMetric(data *MetricData) (valid bool, reason string) {
	// security type includes data from 'init'
	if ok, reason := validateInitMetric(data); !ok {
		return false, reason
	}
	if data.Security == nil {
		return false, missingSecurityError
	}

	if len(*data.Security) == 0 {
		return false, securityMetricsEmptyError
	}
	if len(*data.Security) > securityMetricsMaxLength {
		return false, securityMetricsLengthError
	}
	for i, sm := range *data.Security {
		if sm.Id == nil {
			return false, fmt.Sprintf(securityMetricMissingIdError, i)
		}
		if sm.Name == nil {
			return false, fmt.Sprintf(securityMetricMissingNameError, i)
		}
		if sm.Passed == nil {
			return false, fmt.Sprintf(securityMetricMissingPassedError, i)
		}
	}
	return true, ""
}
