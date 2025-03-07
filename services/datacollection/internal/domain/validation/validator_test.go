package validation

import (
	"testing"
	"time"

	"github.com/kimjooyoon/go_fi_chart/services/datacollection/internal/domain/source"
	"github.com/stretchr/testify/assert"
)

func TestNewStandardValidator(t *testing.T) {
	// 기본 설정으로 생성
	validator := NewStandardValidator(ValidationConfig{})

	// 내부 config 필드에 직접 접근할 수 없으므로 기능 테스트로 검증
	// 기본 설정값이 적용되었는지 간접적으로 확인

	// 간단한 응답 생성
	resp := &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      []source.PriceData{},
	}

	// 오류 없이 처리되는지 확인 (빈 데이터는 검증 실패)
	result, err := validator.ValidateHistoricalData(resp)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestValidateHistoricalData(t *testing.T) {
	// validator 생성
	validator := NewStandardValidator(ValidationConfig{
		MaxPriceThreshold:  1000.0,
		MinPriceThreshold:  1.0,
		MaxVolumeThreshold: 2000000,
		StrictValidation:   true,
	})

	// 유효한 데이터 생성
	now := time.Now()
	validData := []source.PriceData{
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
			Timestamp:     now.Add(24 * time.Hour),
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
		Data:      validData,
	}

	// 유효성 검증 수행
	result, err := validator.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)

	// 유효하지 않은 데이터 생성 (가격 범위 초과)
	invalidData := []source.PriceData{
		{
			Timestamp:     now,
			Open:          100.0,
			High:          1500.0, // MaxPriceThreshold 초과
			Low:           98.0,
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
	}

	response = &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      invalidData,
	}

	// 유효성 검증 수행
	result, err = validator.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Field, "High")

	// 유효하지 않은 데이터 생성 (OHLC 논리 오류)
	invalidLogicData := []source.PriceData{
		{
			Timestamp:     now,
			Open:          100.0,
			High:          105.0,
			Low:           110.0, // Low가 High보다 큼
			Close:         102.0,
			AdjustedClose: 102.0,
			Volume:        1000000,
		},
	}

	response = &source.HistoricalDataResponse{
		Symbol:    "AAPL",
		AssetType: "stock",
		Interval:  "daily",
		Data:      invalidLogicData,
	}

	// 유효성 검증 수행
	result, err = validator.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Field, "Low/High")
}

func TestValidateRealTimeData(t *testing.T) {
	// validator 생성
	validator := NewStandardValidator(ValidationConfig{
		MaxPriceThreshold:  60000.0,
		MinPriceThreshold:  1.0,
		MaxVolumeThreshold: 2000000000,
		MaxDataAge:         10 * time.Minute,
		StrictValidation:   true,
	})

	// 유효한 데이터 생성
	now := time.Now()
	response := &source.RealTimeDataResponse{
		Symbol:        "BTC",
		AssetType:     "crypto",
		CurrentPrice:  50000.0,
		Timestamp:     now,
		Change:        1000.0,
		ChangePercent: 2.0,
		Volume:        1500000000,
		MarketCap:     900000000000,
		High24h:       52000.0,
		Low24h:        49000.0,
	}

	// 유효성 검증 수행
	result, err := validator.ValidateRealTimeData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)

	// 유효하지 않은 데이터 생성 (오래된 데이터)
	response = &source.RealTimeDataResponse{
		Symbol:        "BTC",
		AssetType:     "crypto",
		CurrentPrice:  50000.0,
		Timestamp:     now.Add(-15 * time.Minute), // MaxDataAge 초과
		Change:        1000.0,
		ChangePercent: 2.0,
		Volume:        1500000000,
		MarketCap:     900000000000,
		High24h:       52000.0,
		Low24h:        49000.0,
	}

	// 유효성 검증 수행
	result, err = validator.ValidateRealTimeData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Field, "Timestamp")
	assert.Contains(t, result.Errors[0].Code, "stale")

	// 유효하지 않은 데이터 생성 (논리 오류)
	response = &source.RealTimeDataResponse{
		Symbol:        "BTC",
		AssetType:     "crypto",
		CurrentPrice:  50000.0,
		Timestamp:     now,
		Change:        1000.0,
		ChangePercent: 2.0,
		Volume:        1500000000,
		MarketCap:     900000000000,
		High24h:       52000.0,
		Low24h:        55000.0, // Low24h가 High24h보다 큼
	}

	// 유효성 검증 수행
	result, err = validator.ValidateRealTimeData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.NotEmpty(t, result.Errors)
	assert.Contains(t, result.Errors[0].Field, "Low24h/High24h")
}

func TestValidateMetadata(t *testing.T) {
	// validator 생성
	validator := NewStandardValidator(ValidationConfig{
		RequiredMetadataFields: []string{
			"Symbol",
			"Name",
			"AssetType",
			"Exchange",
		},
		StrictValidation: true,
	})

	// 유효한 데이터 생성
	now := time.Now()
	response := &source.MetadataResponse{
		Symbol:      "MSFT",
		AssetType:   "stock",
		Name:        "Microsoft Corporation",
		Exchange:    "NASDAQ",
		Currency:    "USD",
		Country:     "United States",
		Description: "Technology company",
		Sector:      "Technology",
		Industry:    "Software",
		Website:     "https://www.microsoft.com",
		LogoURL:     "https://logo.com/msft.png",
		LastUpdated: now,
	}

	// 유효성 검증 수행
	result, err := validator.ValidateMetadata(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Errors)

	// 유효하지 않은 데이터 생성 (필수 필드 누락)
	response = &source.MetadataResponse{
		Symbol:      "MSFT",
		AssetType:   "stock",
		Name:        "", // 필수 필드 누락
		Exchange:    "", // 필수 필드 누락
		Currency:    "USD",
		Country:     "United States",
		Description: "Technology company",
		Sector:      "Technology",
		Industry:    "Software",
		Website:     "https://www.microsoft.com",
		LogoURL:     "https://logo.com/msft.png",
		LastUpdated: now,
	}

	// 유효성 검증 수행
	result, err = validator.ValidateMetadata(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)
	assert.NotEmpty(t, result.Errors)
	assert.Equal(t, 2, len(result.Errors))
}

func TestAddRule(t *testing.T) {
	// validator 생성
	validator := NewStandardValidator(ValidationConfig{
		StrictValidation: true,
	})

	// 룰 생성
	rule := NewTimeIntervalRule(24*time.Hour, SeverityError)

	// 룰 추가
	validator.AddRule(rule)

	// 테스트 데이터 생성 (시간 간격 초과)
	now := time.Now()
	largeGapData := []source.PriceData{
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
			Timestamp:     now.Add(48 * time.Hour), // 48시간 간격 (24시간 초과)
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
		Data:      largeGapData,
	}

	// 유효성 검증 수행
	result, err := validator.ValidateHistoricalData(response)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsValid)

	// 내장 검증 외에 추가한 룰에 의한 오류도 포함되어야 함
	foundRuleError := false
	for _, err := range result.Errors {
		if err.Code == "time_gap" {
			foundRuleError = true
			break
		}
	}
	assert.True(t, foundRuleError, "Rule validation error should be present")
}
