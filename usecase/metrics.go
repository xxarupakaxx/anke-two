package usecase

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	gormPrometheus "gorm.io/plugin/prometheus"
)


type MetricsCollector struct {
	Prefix             string
	Interval           uint32
	questionnaireGauge *prometheus.GaugeVec
	questionGauge      *prometheus.GaugeVec
	respondentGauge    *prometheus.GaugeVec
	responseGauge      *prometheus.GaugeVec
	administratorGauge *prometheus.GaugeVec
}

func (mc *MetricsCollector) Metrics(p *gormPrometheus.Prometheus) []prometheus.Collector {
	return []prometheus.Collector{
		mc.questionnaireGauge,
		mc.questionGauge,
		mc.respondentGauge,
		mc.responseGauge,
		mc.administratorGauge,
	}
}

func (mc *MetricsCollector) Collect(p *gormPrometheus.Prometheus) {

}

func (mc *MetricsCollector) collectQuestionnaireMetrics(ctx context.Context, p *gormPrometheus.Prometheus) error {
	return nil
}

func (mc *MetricsCollector) collectQuestionMetrics(ctx context.Context, p *gormPrometheus.Prometheus) error {
	return nil
}

func (mc *MetricsCollector) collectRespondentMetrics(ctx context.Context, p *gormPrometheus.Prometheus) error {
	return nil
}

func (mc *MetricsCollector) collectResponseMetrics(ctx context.Context, p *gormPrometheus.Prometheus) error {
	return nil
}

func (mc *MetricsCollector) collectAdministratorMetrics(ctx context.Context, p *gormPrometheus.Prometheus) error {
	return nil
}

