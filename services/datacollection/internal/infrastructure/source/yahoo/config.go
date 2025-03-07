package yahoo

import (
	"time"

	"github.com/kimjooyoon/go_fi_chart/services/datacollection/internal/domain/source"
)

// Config는 Yahoo Finance API에 대한 설정을 나타냅니다.
type Config struct {
	APIKey          string        // API 키
	BaseURL         string        // 기본 URL
	Timeout         time.Duration // API 요청 타임아웃
	RateLimitPerMin int           // 분당 요청 제한
	RateLimitPerDay int           // 일일 요청 제한
	RetryCount      int           // 재시도 횟수
	RetryDelay      time.Duration // 재시도 지연
	MaxHistoryDays  int           // 최대 과거 데이터 요청 일수
	UserAgent       string        // User-Agent 헤더
	CacheDuration   time.Duration // 캐시 유지 시간
	ProxyURL        string        // 프록시 URL (선택사항)
}

// NewDefaultConfig는 기본 설정으로 Config를 생성합니다.
func NewDefaultConfig() *Config {
	return &Config{
		APIKey:          "",                                            // Yahoo Finance는 일부 API에 키가 필요 없음
		BaseURL:         "https://query1.finance.yahoo.com/v8/finance", // 기본 API URL
		Timeout:         10 * time.Second,                              // 10초 타임아웃
		RateLimitPerMin: 100,                                           // 분당 100 요청 (보수적 추정)
		RateLimitPerDay: 2000,                                          // 일일 2000 요청 (보수적 추정)
		RetryCount:      3,                                             // 3회 재시도
		RetryDelay:      2 * time.Second,                               // 2초 재시도 지연
		MaxHistoryDays:  7300,                                          // 약 20년 데이터
		UserAgent:       "GoFiChart/1.0",                               // User-Agent 헤더
		CacheDuration:   30 * time.Minute,                              // 30분 캐시
		ProxyURL:        "",                                            // 프록시 없음
	}
}

// GetAPIKey는 API 키를 반환합니다.
func (c *Config) GetAPIKey() string {
	return c.APIKey
}

// GetBaseURL은 기본 URL을 반환합니다.
func (c *Config) GetBaseURL() string {
	return c.BaseURL
}

// GetTimeout은 타임아웃 설정을 반환합니다.
func (c *Config) GetTimeout() time.Duration {
	return c.Timeout
}

// GetRateLimitPerMinute는 분당 요청 제한을 반환합니다.
func (c *Config) GetRateLimitPerMinute() int {
	return c.RateLimitPerMin
}

// GetRateLimitPerDay는 일일 요청 제한을 반환합니다.
func (c *Config) GetRateLimitPerDay() int {
	return c.RateLimitPerDay
}

// GetRetryCount는 재시도 횟수를 반환합니다.
func (c *Config) GetRetryCount() int {
	return c.RetryCount
}

// GetRetryDelay는 재시도 지연 시간을 반환합니다.
func (c *Config) GetRetryDelay() time.Duration {
	return c.RetryDelay
}

// Ensure Config implements SourceConfig interface
var _ source.SourceConfig = (*Config)(nil)
