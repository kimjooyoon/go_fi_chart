package normalization

import (
	"testing"
	"time"

	"github.com/kimjooyoon/go_fi_chart/services/datacollection/internal/domain/source"
	"github.com/stretchr/testify/assert"
)

func TestNewStandardNormalizer(t *testing.T) {
	// 기본 설정으로 생성
	normalizer := NewStandardNormalizer(NormalizationConfig{})

	// 내부 config 필드에 직접 접근할 수 없으므로 기능 테스트로 검증
	// 기본 설정값이 적용되었는지 간접적으로 확인

	// 간단한 응답 생성
	resp := &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      []source.PriceData{},
	}

	// 오류 없이 처리되는지 확인
	normalized, err := normalizer.NormalizeHistoricalData(resp)
	assert.NoError(t, err)
	assert.NotNil(t, normalized)
}

func TestNormalizeHistoricalData(t *testing.T) {
	// 테스트를 위한 타임존 설정
	loc, _ := time.LoadLocation("America/New_York")
	utc := time.UTC

	// normalizer 생성
	normalizer := NewStandardNormalizer(NormalizationConfig{
		DefaultTimezone:     utc,
		DefaultCurrency:     "USD",
		PriceScaleFactor:    1.0,
		VolumeScaleFactor:   1.0,
		InterpolationMethod: "linear",
	})

	// 테스트 데이터 생성
	now := time.Now()
	testData := []source.PriceData{
		{
			Timestamp:     now.In(loc),
			Open:          100.0,
			High:          105.0,
			Low:           98.0,
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
		{
			Timestamp:     now.Add(24 * time.Hour).In(loc),
			Open:          102.0,
			High:          107.0,
			Low:           101.0,
			Close:         106.0,
			AdjustedClose: 106.0,
			Volume:        1200000,
		},
		{
			Timestamp:     now.Add(48 * time.Hour).In(loc),
			Open:          106.0,
			High:          110.0,
			Low:           104.0,
			Close:         108.0,
			AdjustedClose: 108.0,
			Volume:        1500000,
		},
	}

	response := &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      testData,
	}

	// 정규화 수행
	normalized, err := normalizer.NormalizeHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, normalized)
	assert.Equal(t, "AAPL", normalized.Symbol)
	assert.Equal(t, "stock", normalized.AssetType)
	assert.Equal(t, "daily", normalized.Interval)
	assert.Len(t, normalized.Data, len(testData))

	// 타임존 변환 확인
	for i, data := range normalized.Data {
		assert.Equal(t, testData[i].Timestamp.In(utc), data.Timestamp)
		assert.Equal(t, testData[i].Open, data.Open)
		assert.Equal(t, testData[i].High, data.High)
		assert.Equal(t, testData[i].Low, data.Low)
		assert.Equal(t, testData[i].Close, data.Close)
		assert.Equal(t, testData[i].AdjustedClose, data.AdjustedClose)
		assert.Equal(t, testData[i].Volume, data.Volume)
	}

	// 누락된 데이터 처리 테스트
	testDataWithGaps := []source.PriceData{
		{
			Timestamp:     now.In(loc),
			Open:          100.0,
			High:          105.0,
			Low:           98.0,
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
		{
			Timestamp:     now.Add(24 * time.Hour).In(loc),
			Open:          0, // 누락된 데이터
			High:          0, // 누락된 데이터
			Low:           0, // 누락된 데이터
			Close:         0, // 누락된 데이터
			AdjustedClose: 0, // 누락된 데이터
			Volume:        0, // 누락된 데이터
		},
	}

	response = &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      testDataWithGaps,
	}

	// 정규화 설정에 보간 활성화
	normalizerWithInterpolation := NewStandardNormalizer(NormalizationConfig{
		DefaultTimezone:     utc,
		InterpolationMethod: "previous", // 이전 값으로 보간
		MaxGapRatio:         0.5,        // 큰 값으로 설정하여 모든 갭 허용
	})

	normalized, err = normalizerWithInterpolation.NormalizeHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, normalized)

	// 보간 결과 확인 - 두 번째 데이터는 첫 번째 데이터로 보간되어야 함
	assert.Equal(t, testDataWithGaps[0].Open, normalized.Data[1].Open)
	assert.Equal(t, testDataWithGaps[0].High, normalized.Data[1].High)
	assert.Equal(t, testDataWithGaps[0].Low, normalized.Data[1].Low)
	assert.Equal(t, testDataWithGaps[0].Close, normalized.Data[1].Close)
	assert.Equal(t, testDataWithGaps[0].AdjustedClose, normalized.Data[1].AdjustedClose)
	assert.Equal(t, testDataWithGaps[0].Volume, normalized.Data[1].Volume)
}

func TestNormalizeRealTimeData(t *testing.T) {
	// 테스트를 위한 타임존 설정
	loc, _ := time.LoadLocation("Asia/Tokyo")
	utc := time.UTC

	// normalizer 생성
	normalizer := NewStandardNormalizer(NormalizationConfig{
		DefaultTimezone:  utc,
		DefaultCurrency:  "USD",
		PriceScaleFactor: 2.0, // 가격 배율 테스트를 위해 2배로 설정
	})

	// 테스트 데이터 생성
	now := time.Now()
	response := &source.RealTimeDataResponse{
		Symbol:        "BTC",
		AssetType:     "crypto",
		CurrentPrice:  50000.0,
		Timestamp:     now.In(loc),
		Change:        1000.0,
		ChangePercent: 2.0,
		Volume:        1500000,
		MarketCap:     900000000000,
		High24h:       52000.0,
		Low24h:        49000.0,
	}

	// 정규화 수행
	normalized, err := normalizer.NormalizeRealTimeData(response)
	assert.NoError(t, err)
	assert.NotNil(t, normalized)
	assert.Equal(t, "BTC", normalized.Symbol)
	assert.Equal(t, "crypto", normalized.AssetType)

	// 타임존 변환 확인
	assert.Equal(t, now.In(utc), normalized.Timestamp)

	// 가격 정규화 (2배) 확인
	assert.Equal(t, 100000.0, normalized.CurrentPrice)
	assert.Equal(t, 2000.0, normalized.Change)
	assert.Equal(t, 2.0, normalized.ChangePercent) // 퍼센트는 정규화하지 않음
	assert.Equal(t, 1500000, normalized.Volume)    // 거래량에는 PriceScaleFactor가 아닌 VolumeScaleFactor 적용
	assert.Equal(t, 1800000000000.0, normalized.MarketCap)
	assert.Equal(t, 104000.0, normalized.High24h)
	assert.Equal(t, 98000.0, normalized.Low24h)
}

func TestNormalizeMetadata(t *testing.T) {
	// 테스트를 위한 타임존 설정
	loc, _ := time.LoadLocation("Europe/London")
	utc := time.UTC

	// normalizer 생성
	normalizer := NewStandardNormalizer(NormalizationConfig{
		DefaultTimezone: utc,
		DefaultCurrency: "EUR", // 기본 통화를 EUR로 설정
	})

	// 테스트 데이터 생성
	now := time.Now()
	response := &source.MetadataResponse{
		Symbol:      "MSFT",
		AssetType:   "stock",
		Name:        "Microsoft Corporation",
		Exchange:    "NASDAQ",
		Currency:    "", // 비어있는 필드
		Country:     "United States",
		Description: "Technology company",
		Sector:      "Technology",
		Industry:    "Software",
		Website:     "https://www.microsoft.com",
		LogoURL:     "https://logo.com/msft.png",
		LastUpdated: now.In(loc),
	}

	// 정규화 수행
	normalized, err := normalizer.NormalizeMetadata(response)
	assert.NoError(t, err)
	assert.NotNil(t, normalized)
	assert.Equal(t, "MSFT", normalized.Symbol)
	assert.Equal(t, "stock", normalized.AssetType)
	assert.Equal(t, "Microsoft Corporation", normalized.Name)
	assert.Equal(t, "NASDAQ", normalized.Exchange)
	assert.Equal(t, "EUR", normalized.Currency) // 기본값으로 대체
	assert.Equal(t, "United States", normalized.Country)
	assert.Equal(t, "Technology", normalized.Sector)
	assert.Equal(t, "Software", normalized.Industry)
	assert.Equal(t, "https://www.microsoft.com", normalized.Website)
	assert.Equal(t, "https://logo.com/msft.png", normalized.LogoURL)

	// 타임존 변환 확인
	assert.Equal(t, now.In(utc), normalized.LastUpdated)
}
