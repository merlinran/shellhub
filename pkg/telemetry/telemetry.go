package telemetry

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	UsersCount      = "users_count"
	NamespacesCount = "namespaces_count"
	DevicesCount    = "devices_count"
)

type Track struct {
	InstanceID string                 `json:"instance_id"`
	Event      string                 `json:"event"`
	Properties map[string]interface{} `json:"properties"`
}

type Telemetry interface {
	SendTelemetry(name string, properties map[string]interface{})
}

type telemetryService struct {
	instanceID string
	hostURL    string
	cli        *resty.Client
}

func New(instanceID string, hostURL string) Telemetry {
	return &telemetryService{
		instanceID: instanceID,
		hostURL:    hostURL,
		cli:        resty.New(),
	}
}

func (t *telemetryService) SendTelemetry(event string, properties map[string]interface{}) {
	track := &Track{
		InstanceID: t.instanceID,
		Event:      event,
		Properties: properties,
	}

	_, err := t.cli.
		SetHostURL(t.hostURL).
		R().
		SetBody(track).Post("/track")
	if err != nil {
		logrus.Error(err)
	}
}
