package prometheus

import (
	"context"
	"testing"

	"github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func Test_NewExporter_should_create_empty_exporter(t *testing.T) {
	// When
	exporter := NewExporter()

	// Then
	assert.NotNil(t, exporter)
	assert.NotNil(t, exporter.registry)
	assert.Empty(t, exporter.metrics)
}

func Test_Exporter_should_export_metrics(t *testing.T) {
	// Given
	exporter := NewExporter()

	testMetrics := []domain.Metric{
		domain.NewBaseMetric(
			"test_counter",
			domain.TypeCounter,
			domain.NewValue(42.0, map[string]string{"label": "value"}),
			"Test counter metric",
		),
		domain.NewBaseMetric(
			"test_gauge",
			domain.TypeGauge,
			domain.NewValue(123.45, map[string]string{"label": "value"}),
			"Test gauge metric",
		),
	}

	// When
	err := exporter.Export(context.Background(), testMetrics)

	// Then
	assert.NoError(t, err)
	assert.Len(t, exporter.metrics, 2)

	metricFamilies, err := exporter.registry.Gather()
	assert.NoError(t, err)
	assert.Len(t, metricFamilies, 2)

	for _, family := range metricFamilies {
		assert.Contains(t, []string{"test_counter", "test_gauge"}, family.GetName())
		assert.Len(t, family.GetMetric(), 1)
	}
}

func Test_Exporter_should_handle_different_metric_types(t *testing.T) {
	// Given
	exporter := NewExporter()

	testMetrics := []domain.Metric{
		domain.NewBaseMetric(
			"test_counter",
			domain.TypeCounter,
			domain.NewValue(1.0, nil),
			"Test counter",
		),
		domain.NewBaseMetric(
			"test_gauge",
			domain.TypeGauge,
			domain.NewValue(2.0, nil),
			"Test gauge",
		),
		domain.NewBaseMetric(
			"test_histogram",
			domain.TypeHistogram,
			domain.NewValue(3.0, nil),
			"Test histogram",
		),
		domain.NewBaseMetric(
			"test_summary",
			domain.TypeSummary,
			domain.NewValue(4.0, nil),
			"Test summary",
		),
	}

	// When
	err := exporter.Export(context.Background(), testMetrics)

	// Then
	assert.NoError(t, err)
	assert.Len(t, exporter.metrics, 4)

	metricFamilies, err := exporter.registry.Gather()
	assert.NoError(t, err)
	assert.Len(t, metricFamilies, 4)
}

func Test_Exporter_should_update_existing_metrics(t *testing.T) {
	// Given
	exporter := NewExporter()

	metric := domain.NewBaseMetric(
		"test_counter",
		domain.TypeCounter,
		domain.NewValue(1.0, nil),
		"Test counter",
	)

	// When
	err := exporter.Export(context.Background(), []domain.Metric{metric})
	assert.NoError(t, err)

	metric = domain.NewBaseMetric(
		"test_counter",
		domain.TypeCounter,
		domain.NewValue(2.0, nil),
		"Test counter",
	)
	err = exporter.Export(context.Background(), []domain.Metric{metric})

	// Then
	assert.NoError(t, err)
	assert.Len(t, exporter.metrics, 1)

	metricFamilies, err := exporter.registry.Gather()
	assert.NoError(t, err)
	assert.Len(t, metricFamilies, 1)

	family := metricFamilies[0]
	assert.Equal(t, "test_counter", family.GetName())
	assert.Equal(t, 3.0, family.GetMetric()[0].GetCounter().GetValue())
}

type invalidMetric struct {
	name        string
	metricType  domain.Type
	value       domain.Value
	description string
}

func (m *invalidMetric) Name() string        { return m.name }
func (m *invalidMetric) Type() domain.Type   { return "invalid" }
func (m *invalidMetric) Value() domain.Value { return m.value }
func (m *invalidMetric) Description() string { return m.description }

func Test_Exporter_should_handle_invalid_metric_type(t *testing.T) {
	// Given
	exporter := NewExporter()

	metric := &invalidMetric{
		name:        "test_invalid",
		metricType:  "invalid",
		value:       domain.NewValue(1.0, nil),
		description: "Test invalid metric",
	}

	// When
	err := exporter.Export(context.Background(), []domain.Metric{metric})

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "지원하지 않는 메트릭 타입")
}

func Test_Exporter_should_handle_duplicate_registration(t *testing.T) {
	// Given
	exporter := NewExporter()

	metric := domain.NewBaseMetric(
		"test_counter",
		domain.TypeCounter,
		domain.NewValue(1.0, nil),
		"Test counter",
	)

	// When
	err := exporter.Export(context.Background(), []domain.Metric{metric})
	assert.NoError(t, err)

	// Try to register a different collector with the same name
	err = exporter.registry.Register(prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "Duplicate counter",
	}))

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "previously registered descriptor")
}

func TestExporter_Export(t *testing.T) {
	// 테스트 메트릭 생성
	metric := domain.NewBaseMetric(
		"test_metric",
		domain.TypeGauge,
		domain.NewValue(42.0, map[string]string{
			"label1": "value1",
			"label2": "value2",
		}),
		"Test metric description",
	)

	// 익스포터 생성
	exporter := NewExporter()

	// 메트릭 익스포트
	err := exporter.Export(context.Background(), []domain.Metric{metric})
	assert.NoError(t, err)
}
