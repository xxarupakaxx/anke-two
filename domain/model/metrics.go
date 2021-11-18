package model

import "github.com/prometheus/client_golang/prometheus"

type MetricsCollector struct {
	Prefix             string
	Interval           uint32
	questionnaireGauge *prometheus.GaugeVec
	questionGauge      *prometheus.GaugeVec
	respondentGauge    *prometheus.GaugeVec
	responseGauge      *prometheus.GaugeVec
	administratorGauge *prometheus.GaugeVec
}