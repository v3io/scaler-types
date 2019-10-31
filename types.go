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
	ScaleInterval time.Duration
	GroupKind     string
}

type ResourceScalerConfig struct {
	KubeconfigPath    string
	AutoScalerOptions AutoScalerOptions
	DLXOptions        DLXOptions
}

type DLXOptions struct {
	Namespace        string
	TargetNameHeader string
	TargetPathHeader string
	TargetPort       int
	ListenAddress    string
}

type ResourceScaler interface {
	SetScale(Resource, int) error
	GetResources() ([]Resource, error)
	GetConfig() (*ResourceScalerConfig, error)
}

type Resource struct {
	Name               string          `json:"name,omitempty"`
	ScaleResources     []ScaleResource `json:"scale_resources,omitempty"`
	LastScaleState     ScaleState      `json:"last_scale_state,omitempty"`
	LastScaleStateTime time.Time       `json:"last_scale_state_time,omitempty"`
}

func (r Resource) String() string {
	out, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(out)
}

type ScaleResource struct {
	MetricName string        `json:"metric_name,omitempty"`
	WindowSize time.Duration `json:"windows_size,omitempty"`
	Threshold  int           `json:"threshold,omitempty"`
}

func (sr ScaleResource) GetKubernetesMetricName() string {
	return fmt.Sprintf("%s_per_%s", sr.MetricName, shortDurationString(sr.WindowSize))
}

func shortDurationString(d time.Duration) string {
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

type ScaleState string

const (
	NonScaleState             ScaleState = "non"
	ScalingFromZeroScaleState ScaleState = "scalingFromZero"
	ScaledFromZeroScaleState  ScaleState = "scaledFromZero"
	ScalingToZeroScaleState   ScaleState = "scalingToZero"
	ScaledToZeroScaleState    ScaleState = "scaledToZero"
)

func ParseScaleState(scaleStateStr string) (ScaleState, error) {
	switch scaleStateStr {
	case string(NonScaleState):
		return NonScaleState, nil
	case string(ScalingFromZeroScaleState):
		return ScalingFromZeroScaleState, nil
	case string(ScaledFromZeroScaleState):
		return ScaledFromZeroScaleState, nil
	case string(ScalingToZeroScaleState):
		return ScalingToZeroScaleState, nil
	case string(ScaledToZeroScaleState):
		return ScaledToZeroScaleState, nil
	default:
		return "", errors.New(fmt.Sprintf("Unknown scale state: %s", scaleStateStr))
	}
}
