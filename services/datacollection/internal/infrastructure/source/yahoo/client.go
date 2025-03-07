package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/kimjooyoon/go_fi_chart/services/datacollection/internal/domain/source"
)

// Client는 Yahoo Finance API와 상호작용하는 클라이언트를 구현합니다.
type Client struct {
	config     *Config
	httpClient *http.Client

	// 속도 제한 관련
	rateLimitMutex sync.Mutex
	requestCount   struct {
		minute int
		day    int
		reset  time.Time
	}

	// 캐싱 관련
	cacheMutex sync.RWMutex
	cache      map[string]cacheEntry
}

// cacheEntry는 캐시 항목을 나타냅니다.
type cacheEntry struct {
	data      interface{}
	timestamp time.Time
	expiry    time.Time
}

// YahooChartResponse는 Yahoo Finance 차트 API 응답 구조를 정의합니다.
type YahooChartResponse struct {
	Chart struct {
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
	} `json:"chart"`
}

// YahooQuoteResponse는 Yahoo Finance 실시간 시세 API 응답 구조를 정의합니다.
type YahooQuoteResponse struct {
	QuoteResponse struct {
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
	} `json:"quoteResponse"`
}

// YahooSearchResponse는 Yahoo Finance 검색 API 응답 구조를 정의합니다.
type YahooSearchResponse struct {
	ResultSet struct {
		Query  string `json:"query"`
		Result []struct {
			Symbol       string `json:"symbol"`
			Name         string `json:"name"`
			ExchDisp     string `json:"exchDisp"`
			TypeDisp     string `json:"typeDisp"`
			Exchange     string `json:"exchange"`
			Type         string `json:"type"`
			IndustryDisp string `json:"industryDisp,omitempty"`
			SectorDisp   string `json:"sectorDisp,omitempty"`
		} `json:"Result"`
	} `json:"ResultSet"`
}

// NewClient는 새로운 Yahoo Finance API 클라이언트를 생성합니다.
func NewClient(config *Config) *Client {
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	// 프록시 설정
	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		if err == nil {
			httpClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
		}
	}

	return &Client{
		config:     config,
		httpClient: httpClient,
		cache:      make(map[string]cacheEntry),
		requestCount: struct {
			minute int
			day    int
			reset  time.Time
		}{
			reset: time.Now().Add(time.Minute),
		},
	}
}

// FetchHistoricalData는 주어진 자산과 시간 범위에 대한 과거 가격 데이터를 가져옵니다.
func (c *Client) FetchHistoricalData(ctx context.Context, request source.HistoricalDataRequest) (*source.HistoricalDataResponse, error) {
	// 캐시 키 생성
	cacheKey := fmt.Sprintf("hist:%s:%s:%s:%d:%d",
		request.Symbol,
		request.AssetType,
		request.Interval,
		request.StartTime.Unix(),
		request.EndTime.Unix())

	// 캐시 확인
	if data, ok := c.getFromCache(cacheKey); ok {
		return data.(*source.HistoricalDataResponse), nil
	}

	// 요청 API 경로 설정
	interval := yahooIntervalFromDomain(request.Interval)

	// 요청 URL 파라미터 생성
	params := url.Values{}
	params.Add("symbol", request.Symbol)
	params.Add("period1", strconv.FormatInt(request.StartTime.Unix(), 10))
	params.Add("period2", strconv.FormatInt(request.EndTime.Unix(), 10))
	params.Add("interval", interval)
	params.Add("includePrePost", "false")
	params.Add("events", "div,split")

	// 차트 API 호출
	endpoint := fmt.Sprintf("%s/chart/%s?%s", c.config.BaseURL, request.Symbol, params.Encode())

	// 속도 제한 확인 및 요청
	var yahooResp YahooChartResponse
	err := c.doRequest(ctx, endpoint, &yahooResp)
	if err != nil {
		return nil, err
	}

	// 응답 데이터 변환
	response, err := c.convertHistoricalData(yahooResp, request)
	if err != nil {
		return nil, err
	}

	// 캐시에 저장
	c.addToCache(cacheKey, response, c.config.CacheDuration)

	return response, nil
}

// FetchRealTimeData는 주어진 자산에 대한 실시간 가격 데이터를 가져옵니다.
func (c *Client) FetchRealTimeData(ctx context.Context, request source.RealTimeDataRequest) (*source.RealTimeDataResponse, error) {
	// 캐시 키 생성
	cacheKey := fmt.Sprintf("realtime:%s:%s", request.Symbol, request.AssetType)

	// 캐시 확인 - 실시간 데이터는 짧은 시간만 캐싱
	if data, ok := c.getFromCache(cacheKey); ok {
		return data.(*source.RealTimeDataResponse), nil
	}

	// 요청 URL 파라미터 생성
	params := url.Values{}
	params.Add("symbols", request.Symbol)
	params.Add("fields", "regularMarketPrice,regularMarketChange,regularMarketChangePercent,regularMarketVolume,regularMarketDayHigh,regularMarketDayLow,regularMarketTime,marketCap")

	// API 호출
	endpoint := fmt.Sprintf("%s/quote?%s", c.config.BaseURL, params.Encode())
	var yahooResp YahooQuoteResponse
	err := c.doRequest(ctx, endpoint, &yahooResp)
	if err != nil {
		return nil, err
	}

	// 응답 확인
	if len(yahooResp.QuoteResponse.Result) == 0 {
		return nil, source.NewSourceError(c.SourceName(), "NO_DATA", fmt.Sprintf("No data found for symbol: %s", request.Symbol))
	}

	// 응답 데이터 변환
	quoteData := yahooResp.QuoteResponse.Result[0]
	response := &source.RealTimeDataResponse{
		Symbol:        request.Symbol,
		AssetType:     request.AssetType,
		CurrentPrice:  quoteData.RegularMarketPrice,
		Timestamp:     time.Unix(quoteData.RegularMarketTime, 0),
		Change:        quoteData.RegularMarketChange,
		ChangePercent: quoteData.RegularMarketChangePercent,
		Volume:        quoteData.RegularMarketVolume,
		MarketCap:     quoteData.MarketCap,
		High24h:       quoteData.RegularMarketDayHigh,
		Low24h:        quoteData.RegularMarketDayLow,
	}

	// 캐시에 저장 (실시간 데이터는 1분만 캐싱)
	c.addToCache(cacheKey, response, time.Minute)

	return response, nil
}

// GetMetadata는 주어진 자산에 대한 메타데이터를 가져옵니다.
func (c *Client) GetMetadata(ctx context.Context, request source.MetadataRequest) (*source.MetadataResponse, error) {
	// 캐시 키 생성
	cacheKey := fmt.Sprintf("meta:%s:%s", request.Symbol, request.AssetType)

	// 캐시 확인
	if data, ok := c.getFromCache(cacheKey); ok {
		return data.(*source.MetadataResponse), nil
	}

	// 두 가지 API 호출을 통해 메타데이터 수집
	// 1. Quote API를 통한 기본 정보
	params := url.Values{}
	params.Add("symbols", request.Symbol)

	endpoint := fmt.Sprintf("%s/quote?%s", c.config.BaseURL, params.Encode())
	var quoteResp YahooQuoteResponse
	err := c.doRequest(ctx, endpoint, &quoteResp)
	if err != nil {
		return nil, err
	}

	if len(quoteResp.QuoteResponse.Result) == 0 {
		return nil, source.NewSourceError(c.SourceName(), "NO_DATA", fmt.Sprintf("No data found for symbol: %s", request.Symbol))
	}

	quoteData := quoteResp.QuoteResponse.Result[0]

	// 2. 검색 API를 통한 추가 정보
	searchParams := url.Values{}
	searchParams.Add("query", request.Symbol)
	searchParams.Add("region", "US")
	searchParams.Add("lang", "en-US")

	searchEndpoint := fmt.Sprintf("https://query1.finance.yahoo.com/v1/finance/search?%s", searchParams.Encode())
	var searchResp YahooSearchResponse
	err = c.doRequest(ctx, searchEndpoint, &searchResp)

	var sector, industry, name, exchange, website string

	if err == nil && len(searchResp.ResultSet.Result) > 0 {
		for _, result := range searchResp.ResultSet.Result {
			if result.Symbol == request.Symbol {
				sector = result.SectorDisp
				industry = result.IndustryDisp
				name = result.Name
				exchange = result.ExchDisp
				break
			}
		}
	}

	// 검색 API 실패해도 쿼트 데이터로 기본 메타데이터 생성
	if name == "" {
		name = request.Symbol
	}

	if exchange == "" {
		exchange = quoteData.Region
	}

	// 결과 생성
	response := &source.MetadataResponse{
		Symbol:      request.Symbol,
		AssetType:   request.AssetType,
		Name:        name,
		Exchange:    exchange,
		Currency:    quoteData.Currency,
		Country:     quoteData.Region,
		Sector:      sector,
		Industry:    industry,
		Website:     website,
		LastUpdated: time.Now(),
	}

	// 캐시에 저장 (메타데이터는 오래 캐싱)
	c.addToCache(cacheKey, response, 24*time.Hour)

	return response, nil
}

// SourceName은 데이터 소스 이름을 반환합니다.
func (c *Client) SourceName() string {
	return "yahoo_finance"
}

// doRequest는 HTTP 요청을 실행하고 응답을 처리합니다.
func (c *Client) doRequest(ctx context.Context, endpoint string, v interface{}) error {
	// 속도 제한 확인
	if err := c.checkRateLimit(); err != nil {
		return err
	}

	// HTTP 요청 생성
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return source.NewSourceError(c.SourceName(), "REQUEST_ERROR", fmt.Sprintf("Error creating request: %v", err))
	}

	// 헤더 설정
	req.Header.Set("User-Agent", c.config.UserAgent)
	if c.config.APIKey != "" {
		req.Header.Set("X-API-KEY", c.config.APIKey)
	}

	// 요청 실행
	resp, err := c.httpClient.Do(req)
	c.updateRequestCount() // 요청 카운트 업데이트

	// 요청 실패 처리
	if err != nil {
		return c.handleRequestError(err, endpoint)
	}
	defer resp.Body.Close()

	// 응답 코드 확인
	if resp.StatusCode != http.StatusOK {
		return c.handleStatusError(resp)
	}

	// 응답 데이터 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return source.NewSourceError(c.SourceName(), "RESPONSE_ERROR", fmt.Sprintf("Error reading response: %v", err))
	}

	// JSON 파싱
	err = json.Unmarshal(body, v)
	if err != nil {
		return source.NewParseError(c.SourceName(), fmt.Sprintf("Error parsing JSON: %v", err), string(body))
	}

	return nil
}

// convertHistoricalData는 Yahoo Finance 응답을 도메인 모델로 변환합니다.
func (c *Client) convertHistoricalData(yahooResp YahooChartResponse, request source.HistoricalDataRequest) (*source.HistoricalDataResponse, error) {
	if len(yahooResp.Chart.Result) == 0 {
		return nil, source.NewSourceError(c.SourceName(), "NO_DATA", fmt.Sprintf("No historical data found for symbol: %s", request.Symbol))
	}

	result := yahooResp.Chart.Result[0]
	timestamps := result.Timestamp
	quotes := result.Indicators.Quote[0]

	// 데이터 길이 확인
	dataLen := len(timestamps)
	if dataLen == 0 || len(quotes.Open) != dataLen || len(quotes.High) != dataLen ||
		len(quotes.Low) != dataLen || len(quotes.Close) != dataLen || len(quotes.Volume) != dataLen {
		return nil, source.NewSourceError(c.SourceName(), "DATA_MISMATCH", "Data length mismatch in historical data")
	}

	// AdjClose가 있는지 확인
	hasAdjClose := len(result.Indicators.AdjClose) > 0 && len(result.Indicators.AdjClose[0].AdjClose) == dataLen

	// 데이터 변환
	priceData := make([]source.PriceData, dataLen)
	for i := 0; i < dataLen; i++ {
		// null 값 처리 (Yahoo Finance에서는 null을 0으로 표시할 수 있음)
		open := quotes.Open[i]
		high := quotes.High[i]
		low := quotes.Low[i]
		close := quotes.Close[i]
		adjClose := close // 기본값은 일반 종가

		if hasAdjClose {
			adjClose = result.Indicators.AdjClose[0].AdjClose[i]
		}

		priceData[i] = source.PriceData{
			Timestamp:     time.Unix(timestamps[i], 0),
			Open:          open,
			High:          high,
			Low:           low,
			Close:         close,
			Volume:        quotes.Volume[i],
			AdjustedClose: adjClose,
		}
	}

	return &source.HistoricalDataResponse{
		Symbol:    request.Symbol,
		AssetType: request.AssetType,
		Interval:  request.Interval,
		Data:      priceData,
	}, nil
}

// addToCache는 데이터를 캐시에 추가합니다.
func (c *Client) addToCache(key string, data interface{}, duration time.Duration) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	c.cache[key] = cacheEntry{
		data:      data,
		timestamp: time.Now(),
		expiry:    time.Now().Add(duration),
	}
}

// getFromCache는 캐시에서 데이터를 가져옵니다.
func (c *Client) getFromCache(key string) (interface{}, bool) {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	entry, found := c.cache[key]
	if !found {
		return nil, false
	}

	// 만료 확인
	if time.Now().After(entry.expiry) {
		return nil, false
	}

	return entry.data, true
}

// checkRateLimit은 API 호출 전 속도 제한을 확인합니다.
func (c *Client) checkRateLimit() error {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	now := time.Now()

	// 속도 제한 시간 리셋 확인
	if now.After(c.requestCount.reset) {
		c.requestCount.minute = 0
		c.requestCount.reset = now.Add(time.Minute)

		// 일별 카운트는 매일 자정에 리셋
		if now.Day() != c.requestCount.reset.Day() {
			c.requestCount.day = 0
		}
	}

	// 속도 제한 초과 확인
	if c.requestCount.minute >= c.config.GetRateLimitPerMinute() {
		waitTime := c.requestCount.reset.Sub(now)
		return source.NewRateLimitError(c.SourceName(), waitTime)
	}

	if c.requestCount.day >= c.config.GetRateLimitPerDay() {
		// 다음 날 자정까지 대기
		tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		waitTime := tomorrow.Sub(now)
		return source.NewRateLimitError(c.SourceName(), waitTime)
	}

	return nil
}

// updateRequestCount는 요청 후 속도 제한 카운트를 업데이트합니다.
func (c *Client) updateRequestCount() {
	c.rateLimitMutex.Lock()
	defer c.rateLimitMutex.Unlock()

	c.requestCount.minute++
	c.requestCount.day++
}

// handleRequestError는 HTTP 요청 에러를 처리합니다.
func (c *Client) handleRequestError(err error, endpoint string) error {
	return source.NewNetworkError(c.SourceName(), fmt.Sprintf("Request failed for %s: %v", endpoint, err), true)
}

// handleStatusError는 HTTP 상태 코드 에러를 처리합니다.
func (c *Client) handleStatusError(resp *http.Response) error {
	// 상태 코드별 처리
	switch resp.StatusCode {
	case http.StatusTooManyRequests: // 429: Too Many Requests
		retryAfter := 60 * time.Second // 기본값
		if retryHeader := resp.Header.Get("Retry-After"); retryHeader != "" {
			if seconds, err := strconv.Atoi(retryHeader); err == nil {
				retryAfter = time.Duration(seconds) * time.Second
			}
		}
		return source.NewRateLimitError(c.SourceName(), retryAfter)

	case http.StatusForbidden, http.StatusUnauthorized:
		return source.NewSourceError(c.SourceName(), "AUTHORIZATION_ERROR", fmt.Sprintf("Authorization failed: %d %s", resp.StatusCode, resp.Status))

	case http.StatusNotFound:
		return source.NewSourceError(c.SourceName(), "NOT_FOUND", fmt.Sprintf("Resource not found: %s", resp.Status))

	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return source.NewNetworkError(c.SourceName(), fmt.Sprintf("Server error: %d %s", resp.StatusCode, resp.Status), true)

	default:
		return source.NewSourceError(c.SourceName(), "HTTP_ERROR", fmt.Sprintf("HTTP error: %d %s", resp.StatusCode, resp.Status))
	}
}

// yahooIntervalFromDomain은 도메인 간격을 Yahoo Finance 간격으로 변환합니다.
func yahooIntervalFromDomain(interval source.Interval) string {
	switch interval {
	case source.Interval1Min:
		return "1m"
	case source.Interval5Min:
		return "5m"
	case source.Interval15Min:
		return "15m"
	case source.Interval30Min:
		return "30m"
	case source.Interval1Hour:
		return "1h"
	case source.Interval4Hour:
		return "4h"
	case source.IntervalDaily:
		return "1d"
	case source.IntervalWeekly:
		return "1wk"
	case source.IntervalMonthly:
		return "1mo"
	case source.IntervalQuarterly:
		return "3mo"
	case source.IntervalYearly:
		return "1y"
	default:
		return "1d" // 기본 값
	}
}

// Ensure Client implements DataSource interface
var _ source.DataSource = (*Client)(nil)
