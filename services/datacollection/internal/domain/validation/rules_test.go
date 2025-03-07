package validation

import (
	"testing"
	"time"

	"github.com/kimjooyoon/go_fi_chart/services/datacollection/internal/domain/source"
	"github.com/stretchr/testify/assert"
)

func TestTimeIntervalRule(t *testing.T) {
	// 규칙 생성
	rule := NewTimeIntervalRule(24*time.Hour, SeverityError)

	assert.Equal(t, "TimeIntervalRule", rule.GetName())
	assert.Equal(t, SeverityError, rule.GetSeverity())

	// 정상 간격 데이터 생성
	now := time.Now()
	normalData := []source.PriceData{
		{
			Timestamp:     now,
			Open:          100.0,
			High:          105.0,
			Low:           98.0,
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
		{
			Timestamp:     now.Add(12 * time.Hour), // 12시간 간격 (24시간 미만)
			Open:          102.0,
			High:          107.0,
			Low:           101.0,
			Close:         106.0,
			AdjustedClose: 106.0,
			Volume:        1200000,
		},
	}

	response := &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      normalData,
	}

	// 검증 수행
	result, err := rule.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)

	// 비정상 간격 데이터 생성
	abnormalData := []source.PriceData{
		{
			Timestamp:     now,
			Open:          100.0,
			High:          105.0,
			Low:           98.0,
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
		{
			Timestamp:     now.Add(36 * time.Hour), // 36시간 간격 (24시간 초과)
			Open:          102.0,
			High:          107.0,
			Low:           101.0,
			Close:         106.0,
			AdjustedClose: 106.0,
			Volume:        1200000,
		},
	}

	response = &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      abnormalData,
	}

	// 검증 수행
	result, err = rule.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.NotEmpty(t, result.Errors)
	assert.Equal(t, "time_gap", result.Errors[0].Code)

	// 시간 역전 데이터 생성
	reversedData := []source.PriceData{
		{
			Timestamp:     now.Add(12 * time.Hour), // 미래 시점
			Open:          100.0,
			High:          105.0,
			Low:           98.0,
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
		{
			Timestamp:     now, // 과거 시점
			Open:          102.0,
			High:          107.0,
			Low:           101.0,
			Close:         106.0,
			AdjustedClose: 106.0,
			Volume:        1200000,
		},
	}

	response = &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      reversedData,
	}

	// 검증 수행
	result, err = rule.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.NotEmpty(t, result.Errors)
	assert.Equal(t, "time_reversal", result.Errors[0].Code)

	// 실시간 데이터 및 메타데이터 검증은 항상 통과해야 함
	realTimeResponse := &source.RealTimeDataResponse{
		Symbol:        "AAPL",
		AssetType:     "stock",
		CurrentPrice:  100.0,
		Timestamp:     now,
		Change:        2.0,
		ChangePercent: 2.0,
		Volume:        1000000,
	}

	result, err = rule.ValidateRealTimeData(realTimeResponse)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)

	metadataResponse := &source.MetadataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Name:      "Apple Inc.",
	}

	result, err = rule.ValidateMetadata(metadataResponse)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)
}

func TestPriceVolatilityRule(t *testing.T) {
	// 규칙 생성 (최대 변동성 20%)
	rule := NewPriceVolatilityRule(0.2, 3, "range", SeverityWarning)

	assert.Equal(t, "PriceVolatilityRule", rule.GetName())
	assert.Equal(t, SeverityWarning, rule.GetSeverity())

	// 정상 변동성 데이터 생성
	now := time.Now()
	normalData := []source.PriceData{
		{
			Timestamp:     now.Add(-48 * time.Hour),
			Open:          100.0,
			High:          105.0,
			Low:           98.0,
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
		{
			Timestamp:     now.Add(-24 * time.Hour),
			Open:          102.0,
			High:          106.0,
			Low:           100.0,
			Close:         104.0,
			AdjustedClose: 104.0,
			Volume:        1100000,
		},
		{
			Timestamp:     now,
			Open:          104.0,
			High:          108.0,
			Low:           102.0,
			Close:         106.0,
			AdjustedClose: 106.0,
			Volume:        1200000,
		},
		{
			Timestamp:     now.Add(24 * time.Hour),
			Open:          106.0,
			High:          110.0,
			Low:           104.0,
			Close:         108.0,
			AdjustedClose: 108.0,
			Volume:        1300000,
		},
	}

	response := &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      normalData,
	}

	// 검증 수행
	result, err := rule.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)

	// 높은 변동성 데이터 생성
	highVolatilityData := []source.PriceData{
		{
			Timestamp:     now.Add(-48 * time.Hour),
			Open:          100.0,
			High:          105.0,
			Low:           98.0,
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
		{
			Timestamp:     now.Add(-24 * time.Hour),
			Open:          102.0,
			High:          106.0,
			Low:           100.0,
			Close:         104.0,
			AdjustedClose: 104.0,
			Volume:        1100000,
		},
		{
			Timestamp:     now,
			Open:          104.0,
			High:          108.0,
			Low:           102.0,
			Close:         106.0,
			AdjustedClose: 106.0,
			Volume:        1200000,
		},
		{
			Timestamp:     now.Add(24 * time.Hour),
			Open:          106.0,
			High:          140.0, // 매우 높은 가격 (변동성 증가)
			Low:           90.0,  // 매우 낮은 가격 (변동성 증가)
			Close:         130.0, // 급등 (변동성 증가)
			AdjustedClose: 130.0,
			Volume:        1300000,
		},
	}

	response = &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      highVolatilityData,
	}

	// 검증 수행
	result, err = rule.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.NotEmpty(t, result.Errors)
	assert.Equal(t, "high_volatility", result.Errors[0].Code)

	// 실시간 데이터 검증 테스트
	realTimeResponse := &source.RealTimeDataResponse{
		Symbol:        "AAPL",
		AssetType:     "stock",
		CurrentPrice:  100.0,
		Timestamp:     now,
		Change:        2.0,
		ChangePercent: 2.0,
		Volume:        1000000,
		High24h:       120.0, // 높은 변동성
		Low24h:        80.0,  // 높은 변동성
	}

	result, err = rule.ValidateRealTimeData(realTimeResponse)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.NotEmpty(t, result.Errors)
	assert.Equal(t, "high_volatility", result.Errors[0].Code)

	// 정상 실시간 데이터
	normalRealTimeResponse := &source.RealTimeDataResponse{
		Symbol:        "AAPL",
		AssetType:     "stock",
		CurrentPrice:  100.0,
		Timestamp:     now,
		Change:        2.0,
		ChangePercent: 2.0,
		Volume:        1000000,
		High24h:       105.0, // 정상 변동성
		Low24h:        95.0,  // 정상 변동성
	}

	result, err = rule.ValidateRealTimeData(normalRealTimeResponse)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)
}

func TestVolumeAnomalyRule(t *testing.T) {
	// 규칙 생성 (기준 거래량 초과 배수 3배)
	rule := NewVolumeAnomalyRule(3.0, 3, SeverityWarning)

	assert.Equal(t, "VolumeAnomalyRule", rule.GetName())
	assert.Equal(t, SeverityWarning, rule.GetSeverity())

	// 정상 거래량 데이터 생성
	now := time.Now()
	normalData := []source.PriceData{
		{
			Timestamp:     now.Add(-72 * time.Hour),
			Open:          100.0,
			High:          105.0,
			Low:           98.0,
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
		{
			Timestamp:     now.Add(-48 * time.Hour),
			Open:          102.0,
			High:          106.0,
			Low:           100.0,
			Close:         104.0,
			AdjustedClose: 104.0,
			Volume:        1100000,
		},
		{
			Timestamp:     now.Add(-24 * time.Hour),
			Open:          104.0,
			High:          108.0,
			Low:           102.0,
			Close:         106.0,
			AdjustedClose: 106.0,
			Volume:        1200000,
		},
		{
			Timestamp:     now,
			Open:          106.0,
			High:          110.0,
			Low:           104.0,
			Close:         108.0,
			AdjustedClose: 108.0,
			Volume:        1300000, // 약간 증가 (3배 미만)
		},
	}

	response := &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      normalData,
	}

	// 검증 수행
	result, err := rule.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)

	// 비정상 거래량 데이터 생성
	abnormalData := []source.PriceData{
		{
			Timestamp:     now.Add(-72 * time.Hour),
			Open:          100.0,
			High:          105.0,
			Low:           98.0,
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
		{
			Timestamp:     now.Add(-48 * time.Hour),
			Open:          102.0,
			High:          106.0,
			Low:           100.0,
			Close:         104.0,
			AdjustedClose: 104.0,
			Volume:        1100000,
		},
		{
			Timestamp:     now.Add(-24 * time.Hour),
			Open:          104.0,
			High:          108.0,
			Low:           102.0,
			Close:         106.0,
			AdjustedClose: 106.0,
			Volume:        1200000,
		},
		{
			Timestamp:     now,
			Open:          106.0,
			High:          110.0,
			Low:           104.0,
			Close:         108.0,
			AdjustedClose: 108.0,
			Volume:        5000000, // 급증 (3배 이상)
		},
	}

	response = &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      abnormalData,
	}

	// 검증 수행
	result, err = rule.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.NotEmpty(t, result.Errors)
	assert.Equal(t, "volume_anomaly", result.Errors[0].Code)
}

func TestMetadataConsistencyRule(t *testing.T) {
	// 규칙 생성
	rule := NewMetadataConsistencyRule(SeverityError)

	assert.Equal(t, "MetadataConsistencyRule", rule.GetName())
	assert.Equal(t, SeverityError, rule.GetSeverity())

	// 정상 주식 데이터 생성
	now := time.Now()
	stockResponse := &source.MetadataResponse{
		Symbol:      "AAPL",
		AssetType:   "stock",
		Name:        "Apple Inc.",
		Exchange:    "NASDAQ",
		Currency:    "USD",
		Country:     "United States",
		Sector:      "Technology",
		Industry:    "Consumer Electronics",
		Description: "Manufacturer of smartphones, computers, and wearable devices.",
		Website:     "https://www.apple.com",
		LogoURL:     "https://logo.com/aapl.png",
		LastUpdated: now,
	}

	// 검증 수행
	result, err := rule.ValidateMetadata(stockResponse)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)

	// 비정상 주식 데이터 생성 (필수 필드 누락)
	invalidStockResponse := &source.MetadataResponse{
		Symbol:      "AAPL",
		AssetType:   "stock",
		Name:        "Apple Inc.",
		Exchange:    "NASDAQ",
		Currency:    "USD",
		Country:     "", // 필수 필드 누락
		Sector:      "", // 필수 필드 누락
		Industry:    "", // 필수 필드 누락
		Description: "Manufacturer of smartphones, computers, and wearable devices.",
		Website:     "https://www.apple.com",
		LogoURL:     "https://logo.com/aapl.png",
		LastUpdated: now,
	}

	// 검증 수행
	result, err = rule.ValidateMetadata(invalidStockResponse)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.Equal(t, 3, len(result.Errors))

	// ETF 데이터 생성
	etfResponse := &source.MetadataResponse{
		Symbol:      "SPY",
		AssetType:   "etf",
		Name:        "SPDR S&P 500 ETF Trust",
		Exchange:    "NYSE",
		Currency:    "USD",
		Description: "ETF tracking the S&P 500 index",
		LastUpdated: now,
	}

	// 검증 수행
	result, err = rule.ValidateMetadata(etfResponse)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)

	// 알 수 없는 자산 유형 검증
	unknownAssetResponse := &source.MetadataResponse{
		Symbol:      "XXX",
		AssetType:   "unknown_type",
		Name:        "Unknown Asset",
		Exchange:    "EXCHANGE",
		Currency:    "USD",
		LastUpdated: now,
	}

	// 검증 수행
	result, err = rule.ValidateMetadata(unknownAssetResponse)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.Equal(t, 1, len(result.Errors))
	assert.Equal(t, "unknown", result.Errors[0].Code)

	// 과거 데이터 및 실시간 데이터 검증은 항상 통과해야 함
	historicalResponse := &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      []source.PriceData{},
	}

	result, err = rule.ValidateHistoricalData(historicalResponse)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)

	realTimeResponse := &source.RealTimeDataResponse{
		Symbol:        "AAPL",
		AssetType:     "stock",
		CurrentPrice:  100.0,
		Timestamp:     now,
		Change:        2.0,
		ChangePercent: 2.0,
		Volume:        1000000,
	}

	result, err = rule.ValidateRealTimeData(realTimeResponse)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)
}

func TestNormalizeAssetType(t *testing.T) {
	// 다양한 자산 유형 변환 테스트
	assert.Equal(t, "stock", normalizeAssetType("stock"))
	assert.Equal(t, "stock", normalizeAssetType("stocks"))
	assert.Equal(t, "stock", normalizeAssetType("equity"))
	assert.Equal(t, "stock", normalizeAssetType("equities"))

	assert.Equal(t, "etf", normalizeAssetType("etf"))
	assert.Equal(t, "etf", normalizeAssetType("ETF"))
	assert.Equal(t, "etf", normalizeAssetType("exchange-traded fund"))
	assert.Equal(t, "etf", normalizeAssetType("exchange_traded_fund"))

	assert.Equal(t, "crypto", normalizeAssetType("crypto"))
	assert.Equal(t, "crypto", normalizeAssetType("cryptocurrency"))
	assert.Equal(t, "crypto", normalizeAssetType("digital asset"))
	assert.Equal(t, "crypto", normalizeAssetType("digital_asset"))

	assert.Equal(t, "forex", normalizeAssetType("forex"))
	assert.Equal(t, "forex", normalizeAssetType("fx"))
	assert.Equal(t, "forex", normalizeAssetType("foreign exchange"))

	// 알 수 없는 유형은 그대로 반환
	assert.Equal(t, "unknown", normalizeAssetType("unknown"))
	assert.Equal(t, "custom", normalizeAssetType("custom"))
}

func TestCalculatePriceRange(t *testing.T) {
	// 가격 범위 계산 테스트
	data := []source.PriceData{
		{
			High: 100.0,
			Low:  90.0,
		},
		{
			High: 110.0,
			Low:  95.0,
		},
		{
			High: 105.0,
			Low:  85.0,
		},
	}

	// 최대 고가 - 최소 저가 = 110 - 85 = 25
	assert.Equal(t, 25.0, calculatePriceRange(data))

	// 빈 데이터
	assert.Equal(t, 0.0, calculatePriceRange([]source.PriceData{}))
}

func TestCalculatePriceStdDev(t *testing.T) {
	// 표준편차 계산 테스트
	data := []source.PriceData{
		{
			Close: 100.0,
		},
		{
			Close: 110.0,
		},
		{
			Close: 105.0,
		},
		{
			Close: 95.0,
		},
	}

	// 평균: 102.5, 편차 제곱합: 162.5, 표준편차: sqrt(162.5/3) ≈ 7.36
	stdDev := calculatePriceStdDev(data)
	assert.InDelta(t, 7.36, stdDev, 0.01)

	// 데이터가 충분하지 않을 경우
	assert.Equal(t, 0.0, calculatePriceStdDev([]source.PriceData{{Close: 100.0}}))
}
