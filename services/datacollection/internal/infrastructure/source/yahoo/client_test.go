package yahoo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kimjooyoon/go_fi_chart/services/datacollection/internal/domain/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	config := NewDefaultConfig()
	client := NewClient(config)

	assert.NotNil(t, client)
	assert.Equal(t, config, client.config)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.cache)
}

func TestYahooIntervalFromDomain(t *testing.T) {
	testCases := []struct {
		interval source.Interval
		expected string
	}{
		{source.Interval1Min, "1m"},
		{source.Interval5Min, "5m"},
		{source.Interval15Min, "15m"},
		{source.Interval30Min, "30m"},
		{source.Interval1Hour, "1h"},
		{source.Interval4Hour, "4h"},
		{source.IntervalDaily, "1d"},
		{source.IntervalWeekly, "1wk"},
		{source.IntervalMonthly, "1mo"},
		{source.IntervalQuarterly, "3mo"},
		{source.IntervalYearly, "1y"},
		{"unknown", "1d"}, // 기본값 테스트
	}

	for _, tc := range testCases {
		t.Run(string(tc.interval), func(t *testing.T) {
			result := yahooIntervalFromDomain(tc.interval)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestClientSourceName(t *testing.T) {
	client := NewClient(NewDefaultConfig())
	assert.Equal(t, "yahoo_finance", client.SourceName())
}

func TestClientFetchHistoricalData(t *testing.T) {
	// 테스트 서버 설정
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 요청 검증
		assert.Equal(t, "/v8/finance/chart/AAPL", r.URL.Path)

		// 쿼리 파라미터 검증
		query := r.URL.Query()
		assert.Equal(t, "AAPL", query.Get("symbol"))
		assert.NotEmpty(t, query.Get("period1"))
		assert.NotEmpty(t, query.Get("period2"))
		assert.Equal(t, "1d", query.Get("interval"))

		// 모의 응답 반환
		mockChartResponse := YahooChartResponse{
			Chart: struct {
				Result []struct {
					Meta struct {
						Currency           string  `json:"currency"`
						Symbol             string  `json:"symbol"`
						ExchangeName       string  `json:"exchangeName"`
						RegularMarketPrice float64 `json:"regularMarketPrice"`
						PreviousClose      float64 `json:"previousClose"`
						Timezone           string  `json:"timezone"`
					} `json:"meta"`
					Timestamp  []int64 `json:"timestamp"`
					Indicators struct {
						Quote []struct {
							Open   []float64 `json:"open"`
							High   []float64 `json:"high"`
							Low    []float64 `json:"low"`
							Close  []float64 `json:"close"`
							Volume []int64   `json:"volume"`
						} `json:"quote"`
						AdjClose []struct {
							AdjClose []float64 `json:"adjclose"`
						} `json:"adjclose,omitempty"`
					} `json:"indicators"`
				} `json:"result"`
				Error interface{} `json:"error"`
			}{
				Result: []struct {
					Meta struct {
						Currency           string  `json:"currency"`
						Symbol             string  `json:"symbol"`
						ExchangeName       string  `json:"exchangeName"`
						RegularMarketPrice float64 `json:"regularMarketPrice"`
						PreviousClose      float64 `json:"previousClose"`
						Timezone           string  `json:"timezone"`
					} `json:"meta"`
					Timestamp  []int64 `json:"timestamp"`
					Indicators struct {
						Quote []struct {
							Open   []float64 `json:"open"`
							High   []float64 `json:"high"`
							Low    []float64 `json:"low"`
							Close  []float64 `json:"close"`
							Volume []int64   `json:"volume"`
						} `json:"quote"`
						AdjClose []struct {
							AdjClose []float64 `json:"adjclose"`
						} `json:"adjclose,omitempty"`
					} `json:"indicators"`
				}{
					{
						Meta: struct {
							Currency           string  `json:"currency"`
							Symbol             string  `json:"symbol"`
							ExchangeName       string  `json:"exchangeName"`
							RegularMarketPrice float64 `json:"regularMarketPrice"`
							PreviousClose      float64 `json:"previousClose"`
							Timezone           string  `json:"timezone"`
						}{
							Currency:     "USD",
							Symbol:       "AAPL",
							ExchangeName: "NASDAQ",
						},
						Timestamp: []int64{1609459200, 1609545600, 1609632000},
						Indicators: struct {
							Quote []struct {
								Open   []float64 `json:"open"`
								High   []float64 `json:"high"`
								Low    []float64 `json:"low"`
								Close  []float64 `json:"close"`
								Volume []int64   `json:"volume"`
							} `json:"quote"`
							AdjClose []struct {
								AdjClose []float64 `json:"adjclose"`
							} `json:"adjclose,omitempty"`
						}{
							Quote: []struct {
								Open   []float64 `json:"open"`
								High   []float64 `json:"high"`
								Low    []float64 `json:"low"`
								Close  []float64 `json:"close"`
								Volume []int64   `json:"volume"`
							}{
								{
									Open:   []float64{133.52, 134.08, 135.83},
									High:   []float64{134.74, 135.99, 136.69},
									Low:    []float64{131.72, 132.43, 133.51},
									Close:  []float64{132.69, 135.55, 136.01},
									Volume: []int64{100000000, 98000000, 95000000},
								},
							},
							AdjClose: []struct {
								AdjClose []float64 `json:"adjclose"`
							}{
								{
									AdjClose: []float64{132.69, 135.55, 136.01},
								},
							},
						},
					},
				},
				Error: nil,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockChartResponse)
	}))
	defer server.Close()

	// 테스트 클라이언트 설정
	config := NewDefaultConfig()
	config.BaseURL = server.URL + "/v8/finance"
	client := NewClient(config)

	// 테스트 요청
	ctx := context.Background()
	now := time.Now()
	startTime := now.Add(-7 * 24 * time.Hour)
	endTime := now

	request := source.HistoricalDataRequest{
		Symbol:    "AAPL",
		AssetType: source.AssetTypeStock,
		Interval:  source.IntervalDaily,
		StartTime: startTime,
		EndTime:   endTime,
	}

	// 함수 호출 및 검증
	response, err := client.FetchHistoricalData(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, response)

	// 응답 검증
	assert.Equal(t, "AAPL", response.Symbol)
	assert.Equal(t, source.AssetTypeStock, response.AssetType)
	assert.Equal(t, source.IntervalDaily, response.Interval)
	assert.Len(t, response.Data, 3)

	// 첫 번째 데이터 포인트 검증
	firstDataPoint := response.Data[0]
	assert.Equal(t, float64(133.52), firstDataPoint.Open)
	assert.Equal(t, float64(134.74), firstDataPoint.High)
	assert.Equal(t, float64(131.72), firstDataPoint.Low)
	assert.Equal(t, float64(132.69), firstDataPoint.Close)
	assert.Equal(t, int64(100000000), firstDataPoint.Volume)
	assert.Equal(t, float64(132.69), firstDataPoint.AdjustedClose)

	// 캐싱 테스트
	cachedResponse, err := client.FetchHistoricalData(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, response, cachedResponse) // 동일한 응답이어야 함
}

func TestClientFetchRealTimeData(t *testing.T) {
	// 테스트 서버 설정
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 요청 검증
		assert.Equal(t, "/v8/finance/quote", r.URL.Path)

		// 쿼리 파라미터 검증
		query := r.URL.Query()
		assert.Equal(t, "AAPL", query.Get("symbols"))

		// 모의 응답 반환
		mockQuoteResponse := YahooQuoteResponse{
			QuoteResponse: struct {
				Result []struct {
					Symbol                     string  `json:"symbol"`
					Language                   string  `json:"language"`
					Region                     string  `json:"region"`
					QuoteType                  string  `json:"quoteType"`
					Currency                   string  `json:"currency"`
					MarketState                string  `json:"marketState"`
					RegularMarketPrice         float64 `json:"regularMarketPrice"`
					RegularMarketChange        float64 `json:"regularMarketChange"`
					RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
					RegularMarketVolume        int64   `json:"regularMarketVolume"`
					RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
					RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
					RegularMarketTime          int64   `json:"regularMarketTime"`
					MarketCap                  float64 `json:"marketCap,omitempty"`
				} `json:"result"`
				Error interface{} `json:"error"`
			}{
				Result: []struct {
					Symbol                     string  `json:"symbol"`
					Language                   string  `json:"language"`
					Region                     string  `json:"region"`
					QuoteType                  string  `json:"quoteType"`
					Currency                   string  `json:"currency"`
					MarketState                string  `json:"marketState"`
					RegularMarketPrice         float64 `json:"regularMarketPrice"`
					RegularMarketChange        float64 `json:"regularMarketChange"`
					RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
					RegularMarketVolume        int64   `json:"regularMarketVolume"`
					RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
					RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
					RegularMarketTime          int64   `json:"regularMarketTime"`
					MarketCap                  float64 `json:"marketCap,omitempty"`
				}{
					{
						Symbol:                     "AAPL",
						Language:                   "en-US",
						Region:                     "US",
						QuoteType:                  "EQUITY",
						Currency:                   "USD",
						MarketState:                "REGULAR",
						RegularMarketPrice:         150.10,
						RegularMarketChange:        2.5,
						RegularMarketChangePercent: 1.69,
						RegularMarketVolume:        87654321,
						RegularMarketDayHigh:       151.20,
						RegularMarketDayLow:        148.50,
						RegularMarketTime:          time.Now().Unix(),
						MarketCap:                  2500000000000,
					},
				},
				Error: nil,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockQuoteResponse)
	}))
	defer server.Close()

	// 테스트 클라이언트 설정
	config := NewDefaultConfig()
	config.BaseURL = server.URL + "/v8/finance"
	client := NewClient(config)

	// 테스트 요청
	ctx := context.Background()
	request := source.RealTimeDataRequest{
		Symbol:    "AAPL",
		AssetType: source.AssetTypeStock,
	}

	// 함수 호출 및 검증
	response, err := client.FetchRealTimeData(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, response)

	// 응답 검증
	assert.Equal(t, "AAPL", response.Symbol)
	assert.Equal(t, source.AssetTypeStock, response.AssetType)
	assert.Equal(t, float64(150.10), response.CurrentPrice)
	assert.Equal(t, float64(2.5), response.Change)
	assert.Equal(t, float64(1.69), response.ChangePercent)
	assert.Equal(t, int64(87654321), response.Volume)
	assert.Equal(t, float64(151.20), response.High24h)
	assert.Equal(t, float64(148.50), response.Low24h)
	assert.Equal(t, float64(2500000000000), response.MarketCap)
}

func TestCacheOperations(t *testing.T) {
	client := NewClient(NewDefaultConfig())

	// 캐시 추가
	testData := "test data"
	client.addToCache("test_key", testData, 10*time.Minute)

	// 캐시 조회
	data, found := client.getFromCache("test_key")
	assert.True(t, found)
	assert.Equal(t, testData, data)

	// 존재하지 않는 키 조회
	_, found = client.getFromCache("nonexistent_key")
	assert.False(t, found)

	// 만료된 캐시 조회
	client.addToCache("expired_key", testData, -1*time.Minute) // 이미 만료된 시간
	_, found = client.getFromCache("expired_key")
	assert.False(t, found)
}

func TestRateLimitHandling(t *testing.T) {
	config := NewDefaultConfig()
	config.RateLimitPerMin = 2
	client := NewClient(config)

	// 첫 번째 요청: 성공해야 함
	err := client.checkRateLimit()
	assert.NoError(t, err)
	client.updateRequestCount()

	// 두 번째 요청: 성공해야 함
	err = client.checkRateLimit()
	assert.NoError(t, err)
	client.updateRequestCount()

	// 세 번째 요청: 속도 제한 초과로 실패해야 함
	err = client.checkRateLimit()
	assert.Error(t, err)

	// 에러 타입 검증
	rateLimitErr, ok := err.(*source.RateLimitError)
	assert.True(t, ok)
	assert.Equal(t, "yahoo_finance", rateLimitErr.SourceName())
	assert.Equal(t, "RATE_LIMIT_EXCEEDED", rateLimitErr.ErrorCode())
	assert.True(t, rateLimitErr.GetRetryAfter() > 0)
}
