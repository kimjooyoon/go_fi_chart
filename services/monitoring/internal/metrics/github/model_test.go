package github

import (
	"testing"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGitHubMetric(t *testing.T) {
	// Given
	repository := "test-repo"
	metricType := MetricTypeStars
	value := 42.0
	timestamp := time.Now()

	// When
	metric := NewMetric(repository, metricType, value, timestamp)

	// Then
	assert.NotNil(t, metric)
	assert.Equal(t, repository, metric.Repository())
	assert.Equal(t, metricType, metric.MetricType())
	assert.Equal(t, value, metric.Value())
	assert.Equal(t, timestamp, metric.Timestamp())
}

func TestGitHubMetricType(t *testing.T) {
	t.Run("유효한 메트릭 타입", func(t *testing.T) {
		validTypes := []MetricType{
			MetricTypeStars,
			MetricTypeForks,
			MetricTypeIssues,
			MetricTypePullRequests,
			MetricTypeContributors,
		}

		for _, metricType := range validTypes {
			metric := NewMetric("test-repo", metricType, 1.0, time.Now())
			assert.NoError(t, metric.Validate())
			assert.True(t, metricType.IsValid())
		}
	})

	t.Run("잘못된 메트릭 타입", func(t *testing.T) {
		metric := NewMetric("test-repo", "INVALID_TYPE", 1.0, time.Now())
		assert.Error(t, metric.Validate())
		assert.False(t, MetricType("INVALID_TYPE").IsValid())
	})
}

func TestGitHubMetricToDomain(t *testing.T) {
	// Given
	repository := "test-repo"
	metricType := MetricTypeStars
	value := 42.0
	timestamp := time.Now()
	metric := NewMetric(repository, metricType, value, timestamp)

	// When
	domainMetric := metric.ToDomain()

	// Then
	assert.NotNil(t, domainMetric)
	assert.Equal(t, repository+"."+string(metricType), domainMetric.ID())
	assert.Equal(t, domain.MetricTypeGitHub, domainMetric.Type())
	assert.Equal(t, value, domainMetric.Value().Value())
	assert.Equal(t, timestamp, domainMetric.Timestamp())

	labels := domainMetric.Value().Labels()
	assert.Equal(t, repository, labels["repository"])
	assert.Equal(t, string(metricType), labels["type"])
}
