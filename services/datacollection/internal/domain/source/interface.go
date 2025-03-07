package source

import (
	"context"
	"time"
)

// DataSource는 금융 데이터 소스(예: Yahoo Finance, Alpha Vantage 등)에 대한 인터페이스를 정의합니다.
type DataSource interface {
	// FetchHistoricalData는 주어진 자산과 시간 범위에 대한 과거 가격 데이터를 가져옵니다.
	FetchHistoricalData(ctx context.Context, request HistoricalDataRequest) (*HistoricalDataResponse, error)

	// FetchRealTimeData는 주어진 자산에 대한 실시간 가격 데이터를 가져옵니다.
	FetchRealTimeData(ctx context.Context, request RealTimeDataRequest) (*RealTimeDataResponse, error)

	// GetMetadata는 주어진 자산에 대한 메타데이터를 가져옵니다.
	GetMetadata(ctx context.Context, request MetadataRequest) (*MetadataResponse, error)

	// SourceName은 데이터 소스의 이름을 반환합니다.
	SourceName() string
}

// SourceConfig는 데이터 소스에 대한 기본 구성을 정의합니다.
type SourceConfig interface {
	// GetAPIKey는 데이터 소스에 접근하기 위한 API 키를 반환합니다.
	GetAPIKey() string

	// GetBaseURL은 데이터 소스의 기본 URL을 반환합니다.
	GetBaseURL() string

	// GetTimeout은 API 요청에 대한 타임아웃을 반환합니다.
	GetTimeout() time.Duration

	// GetRateLimitPerMinute는 분당 허용되는 API 요청 수를 반환합니다.
	GetRateLimitPerMinute() int

	// GetRateLimitPerDay는 일일 허용되는 API 요청 수를 반환합니다.
	GetRateLimitPerDay() int

	// GetRetryCount는 API 요청 실패 시 재시도 횟수를 반환합니다.
	GetRetryCount() int

	// GetRetryDelay는 API 요청 재시도 간의 지연 시간을 반환합니다.
	GetRetryDelay() time.Duration
}
