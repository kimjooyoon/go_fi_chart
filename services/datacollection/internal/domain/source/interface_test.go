package source

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDataSource는 DataSource 인터페이스의 모의 구현입니다.
type MockDataSource struct {
	mock.Mock
}

// FetchHistoricalData는 과거 데이터를 가져오는 모의 메서드입니다.
func (m *MockDataSource) FetchHistoricalData(ctx context.Context, request HistoricalDataRequest) (*HistoricalDataResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*HistoricalDataResponse), args.Error(1)
}

// FetchRealTimeData는 실시간 데이터를 가져오는 모의 메서드입니다.
func (m *MockDataSource) FetchRealTimeData(ctx context.Context, request RealTimeDataRequest) (*RealTimeDataResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RealTimeDataResponse), args.Error(1)
}

// GetMetadata는 메타데이터를 가져오는 모의 메서드입니다.
func (m *MockDataSource) GetMetadata(ctx context.Context, request MetadataRequest) (*MetadataResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MetadataResponse), args.Error(1)
}

// SourceName은 데이터 소스 이름을 반환하는 모의 메서드입니다.
func (m *MockDataSource) SourceName() string {
	args := m.Called()
	return args.String(0)
}

// MockSourceConfig는 SourceConfig 인터페이스의 모의 구현입니다.
type MockSourceConfig struct {
	mock.Mock
}

// GetAPIKey는 API 키를 반환하는 모의 메서드입니다.
func (m *MockSourceConfig) GetAPIKey() string {
	args := m.Called()
	return args.String(0)
}

// GetBaseURL은 기본 URL을 반환하는 모의 메서드입니다.
func (m *MockSourceConfig) GetBaseURL() string {
	args := m.Called()
	return args.String(0)
}

// GetTimeout은 타임아웃을 반환하는 모의 메서드입니다.
func (m *MockSourceConfig) GetTimeout() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

// GetRateLimitPerMinute는 분당 요청 제한을 반환하는 모의 메서드입니다.
func (m *MockSourceConfig) GetRateLimitPerMinute() int {
	args := m.Called()
	return args.Int(0)
}

// GetRateLimitPerDay는 일일 요청 제한을 반환하는 모의 메서드입니다.
func (m *MockSourceConfig) GetRateLimitPerDay() int {
	args := m.Called()
	return args.Int(0)
}

// GetRetryCount는 재시도 횟수를 반환하는 모의 메서드입니다.
func (m *MockSourceConfig) GetRetryCount() int {
	args := m.Called()
	return args.Int(0)
}

// GetRetryDelay는 재시도 지연 시간을 반환하는 모의 메서드입니다.
func (m *MockSourceConfig) GetRetryDelay() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

// TestMockDataSource는 MockDataSource가 DataSource 인터페이스를 구현하는지 테스트합니다.
func TestMockDataSource(t *testing.T) {
	var _ DataSource = (*MockDataSource)(nil)
}

// TestMockSourceConfig는 MockSourceConfig가 SourceConfig 인터페이스를 구현하는지 테스트합니다.
func TestMockSourceConfig(t *testing.T) {
	var _ SourceConfig = (*MockSourceConfig)(nil)
}

// TestSourceErrorInterface는 다양한 에러 타입이 SourceError 인터페이스를 구현하는지 테스트합니다.
func TestSourceErrorInterface(t *testing.T) {
	var err SourceError

	// BaseSourceError 테스트
	baseErr := NewSourceError("test_source", "TEST_ERROR", "Test error message")
	err = baseErr
	assert.Equal(t, "test_source", err.SourceName())
	assert.Equal(t, "TEST_ERROR", err.ErrorCode())
	assert.Contains(t, err.Error(), "Test error message")

	// RateLimitError 테스트
	rateLimitErr := NewRateLimitError("test_source", 5*time.Minute)
	err = rateLimitErr
	assert.Equal(t, "test_source", err.SourceName())
	assert.Equal(t, "RATE_LIMIT_EXCEEDED", err.ErrorCode())
	assert.Contains(t, err.Error(), "Rate limit exceeded")
	assert.Equal(t, 5*time.Minute, rateLimitErr.GetRetryAfter())

	// NetworkError 테스트
	networkErr := NewNetworkError("test_source", "Connection timeout", true)
	err = networkErr
	assert.Equal(t, "test_source", err.SourceName())
	assert.Equal(t, "NETWORK_ERROR", err.ErrorCode())
	assert.Contains(t, err.Error(), "Connection timeout")
	assert.True(t, networkErr.IsRetryable())

	// ParseError 테스트
	parseErr := NewParseError("test_source", "Invalid JSON format", "{invalid:json}")
	err = parseErr
	assert.Equal(t, "test_source", err.SourceName())
	assert.Equal(t, "PARSE_ERROR", err.ErrorCode())
	assert.Contains(t, err.Error(), "Invalid JSON format")
	assert.Equal(t, "{invalid:json}", parseErr.GetRawData())
}
