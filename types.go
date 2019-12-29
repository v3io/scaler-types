package scaler_types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type AutoScalerOptions struct {
	Namespace     string
	ScaleInterval Duration
	GroupKind     string
}

type ResourceScalerConfig struct {
	KubeconfigPath    string
	AutoScalerOptions AutoScalerOptions
	DLXOptions        DLXOptions
}

type DLXOptions struct {
	Namespace                string
	TargetNameHeader         string
	TargetPathHeader         string
	TargetPort               int
	ListenAddress            string
	ResourceReadinessTimeout Duration
}

type ResourceScaler interface {
	SetScale([]Resource, int) error
	GetResources() ([]Resource, error)
	GetConfig() (*ResourceScalerConfig, error)
}

type Resource struct {
	Name               string          `json:"name,omitempty"`
	ScaleResources     []ScaleResource `json:"scale_resources,omitempty"`
	LastScaleEvent     *ScaleEvent     `json:"last_scale_event,omitempty"`
	LastScaleEventTime *time.Time      `json:"last_scale_event_time,omitempty"`
}

func (r Resource) String() string {
	out, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(out)
}

type ScaleResource struct {
	MetricName string   `json:"metric_name,omitempty"`
	WindowSize Duration `json:"windows_size,omitempty"`
	Threshold  int      `json:"threshold,omitempty"`
}

func (sr ScaleResource) GetKubernetesMetricName() string {
	return fmt.Sprintf("%s_per_%s", sr.MetricName, shortDurationString(sr.WindowSize))
}

func shortDurationString(d Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

func (sr ScaleResource) String() string {
	out, err := json.Marshal(sr)
	if err != nil {
		panic(err)
	}
	return string(out)
}

type ScaleEvent string

const (
	ResourceUpdatedScaleEvent        ScaleEvent = "resourceUpdated"
	ScaleFromZeroStartedScaleEvent   ScaleEvent = "scaleFromZeroStarted"
	ScaleFromZeroCompletedScaleEvent ScaleEvent = "scaleFromZeroCompleted"
	ScaleToZeroStartedScaleEvent     ScaleEvent = "scaleToZeroStarted"
	ScaleToZeroCompletedScaleEvent   ScaleEvent = "scaleToZeroCompleted"
)

func ParseScaleEvent(scaleEventStr string) (ScaleEvent, error) {
	switch scaleEventStr {
	case string(ResourceUpdatedScaleEvent):
		return ResourceUpdatedScaleEvent, nil
	case string(ScaleFromZeroStartedScaleEvent):
		return ScaleFromZeroStartedScaleEvent, nil
	case string(ScaleFromZeroCompletedScaleEvent):
		return ScaleFromZeroCompletedScaleEvent, nil
	case string(ScaleToZeroStartedScaleEvent):
		return ScaleToZeroStartedScaleEvent, nil
	case string(ScaleToZeroCompletedScaleEvent):
		return ScaleToZeroCompletedScaleEvent, nil
	default:
		return "", errors.New(fmt.Sprintf("Unknown scale event: %s", scaleEventStr))
	}
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
