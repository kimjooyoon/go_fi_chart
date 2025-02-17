package metrics

import (
	"testing"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMetricType(t *testing.T) {
	t.Run("메트릭 타입 문자열 변환", func(t *testing.T) {
		tests := []struct {
			metricType domain.MetricType
			expected   string
		}{
			{domain.MetricTypeAssetValue, string(domain.MetricTypeAssetValue)},
			{domain.MetricTypeTransactionCount, string(domain.MetricTypeTransactionCount)},
			{domain.MetricTypePortfolioValue, string(domain.MetricTypePortfolioValue)},
			{domain.MetricTypeUserCount, string(domain.MetricTypeUserCount)},
		}

		for _, test := range tests {
			assert.Equal(t, test.expected, string(test.metricType))
		}
	})

	t.Run("메트릭 타입 유효성 검사", func(t *testing.T) {
		tests := []struct {
			metricType domain.MetricType
			isValid    bool
		}{
			{domain.MetricTypeAssetValue, true},
			{domain.MetricTypeTransactionCount, true},
			{domain.MetricTypePortfolioValue, true},
			{domain.MetricTypeUserCount, true},
			{"INVALID_TYPE", false},
		}

		for _, test := range tests {
			err := domain.NewMetric(
				"test-metric",
				test.metricType,
				domain.NewMetricValue(100.0, nil),
				time.Now(),
			).Validate()
			assert.Equal(t, test.isValid, err == nil)
		}
	})
}

func TestMetricValue(t *testing.T) {
	t.Run("메트릭 값 생성", func(t *testing.T) {
		value := domain.NewMetricValue(100.0, nil)
		assert.Equal(t, 100.0, value.Value())
	})

	t.Run("메트릭 값 비교", func(t *testing.T) {
		value1 := domain.NewMetricValue(100.0, nil)
		value2 := domain.NewMetricValue(100.0, nil)
		value3 := domain.NewMetricValue(200.0, nil)

		assert.True(t, value1.Equals(value2))
		assert.False(t, value1.Equals(value3))
	})

	t.Run("메트릭 값 연산", func(t *testing.T) {
		value1 := domain.NewMetricValue(100.0, nil)
		value2 := domain.NewMetricValue(50.0, nil)

		sum := value1.Add(value2)
		assert.Equal(t, 150.0, sum.Value())
	})
}

func TestMetric(t *testing.T) {
	t.Run("메트릭 생성", func(t *testing.T) {
		now := time.Now()
		metric := domain.NewMetric(
			"test-metric",
			domain.MetricTypeAssetValue,
			domain.NewMetricValue(100.0, nil),
			now,
		)

		assert.Equal(t, "test-metric", metric.ID())
		assert.Equal(t, domain.MetricTypeAssetValue, metric.Type())
		assert.Equal(t, 100.0, metric.Value().Value())
		assert.Equal(t, now, metric.Timestamp())
	})

	t.Run("메트릭 유효성 검사", func(t *testing.T) {
		tests := []struct {
			name      string
			metric    *domain.Metric
			isValid   bool
			errorType error
		}{
			{
				name: "유효한 메트릭",
				metric: domain.NewMetric(
					"test-metric",
					domain.MetricTypeAssetValue,
					domain.NewMetricValue(100.0, nil),
					time.Now(),
				),
				isValid: true,
			},
			{
				name: "잘못된 메트릭 타입",
				metric: domain.NewMetric(
					"test-metric",
					"INVALID_TYPE",
					domain.NewMetricValue(100.0, nil),
					time.Now(),
				),
				isValid:   false,
				errorType: domain.ErrInvalidMetricType,
			},
			{
				name: "빈 메트릭 ID",
				metric: domain.NewMetric(
					"",
					domain.MetricTypeAssetValue,
					domain.NewMetricValue(100.0, nil),
					time.Now(),
				),
				isValid:   false,
				errorType: domain.ErrInvalidMetricID,
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				err := test.metric.Validate()
				if test.isValid {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
					assert.Equal(t, test.errorType, err)
				}
			})
		}
	})
}

func TestMetricLabels(t *testing.T) {
	t.Run("메트릭 레이블 추가", func(t *testing.T) {
		labels := make(map[string]string)
		metric := domain.NewMetric(
			"test-metric",
			domain.MetricTypeAssetValue,
			domain.NewMetricValue(100.0, labels),
			time.Now(),
		)

		labels["env"] = "production"
		labels["service"] = "asset"

		assert.Equal(t, "production", metric.Value().Labels()["env"])
		assert.Equal(t, "asset", metric.Value().Labels()["service"])
	})

	t.Run("메트릭 레이블 제거", func(t *testing.T) {
		labels := map[string]string{"env": "production"}
		metric := domain.NewMetric(
			"test-metric",
			domain.MetricTypeAssetValue,
			domain.NewMetricValue(100.0, labels),
			time.Now(),
		)

		delete(labels, "env")

		_, exists := metric.Value().Labels()["env"]
		assert.False(t, exists)
	})
}
