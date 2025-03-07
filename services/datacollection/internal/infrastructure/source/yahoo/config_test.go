package yahoo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	// 기본 설정 값 검증
	assert.Equal(t, "", config.APIKey)
	assert.Equal(t, "https://query1.finance.yahoo.com/v8/finance", config.BaseURL)
	assert.Equal(t, 10*time.Second, config.Timeout)
	assert.Equal(t, 100, config.RateLimitPerMin)
	assert.Equal(t, 2000, config.RateLimitPerDay)
	assert.Equal(t, 3, config.RetryCount)
	assert.Equal(t, 2*time.Second, config.RetryDelay)
	assert.Equal(t, 7300, config.MaxHistoryDays)
	assert.Equal(t, "GoFiChart/1.0", config.UserAgent)
	assert.Equal(t, 30*time.Minute, config.CacheDuration)
	assert.Equal(t, "", config.ProxyURL)
}

func TestConfigGetters(t *testing.T) {
	config := &Config{
		APIKey:          "test_key",
		BaseURL:         "https://test.example.com",
		Timeout:         5 * time.Second,
		RateLimitPerMin: 50,
		RateLimitPerDay: 1000,
		RetryCount:      5,
		RetryDelay:      3 * time.Second,
	}

	// 게터 메서드 검증
	assert.Equal(t, "test_key", config.GetAPIKey())
	assert.Equal(t, "https://test.example.com", config.GetBaseURL())
	assert.Equal(t, 5*time.Second, config.GetTimeout())
	assert.Equal(t, 50, config.GetRateLimitPerMinute())
	assert.Equal(t, 1000, config.GetRateLimitPerDay())
	assert.Equal(t, 5, config.GetRetryCount())
	assert.Equal(t, 3*time.Second, config.GetRetryDelay())
}
