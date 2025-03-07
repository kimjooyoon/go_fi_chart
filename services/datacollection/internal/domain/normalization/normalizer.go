package normalization

import (
	"fmt"
	"math"
	"time"

	"github.com/kimjooyoon/go_fi_chart/services/datacollection/internal/domain/source"
)

// Normalizer는 데이터 정규화를 위한 인터페이스를 정의합니다.
type Normalizer interface {
	// NormalizeHistoricalData는 과거 가격 데이터를 정규화합니다.
	NormalizeHistoricalData(response *source.HistoricalDataResponse) (*source.HistoricalDataResponse, error)

	// NormalizeRealTimeData는 실시간 가격 데이터를 정규화합니다.
	NormalizeRealTimeData(response *source.RealTimeDataResponse) (*source.RealTimeDataResponse, error)

	// NormalizeMetadata는 자산 메타데이터를 정규화합니다.
	NormalizeMetadata(response *source.MetadataResponse) (*source.MetadataResponse, error)
}

// NormalizationConfig는 정규화 설정을 정의합니다.
type NormalizationConfig struct {
	// DefaultTimezone은 시간대 정규화를 위한 기본 타임존입니다.
	DefaultTimezone *time.Location

	// DefaultCurrency는 통화 정규화를 위한 기본 통화입니다.
	DefaultCurrency string

	// InterpolationMethod는 누락된 데이터 보간 방법입니다.
	// "linear", "previous", "zero", "none" 중 하나의 값을 가질 수 있습니다.
	InterpolationMethod string

	// PriceScaleFactor는 가격 데이터에 적용할 스케일 팩터입니다.
	PriceScaleFactor float64

	// VolumeScaleFactor는 거래량 데이터에 적용할 스케일 팩터입니다.
	VolumeScaleFactor float64

	// MaxGapRatio는 허용되는 최대 갭 비율입니다.
	// 이 값을 초과하는 갭은 보간되지 않고 에러로 처리됩니다.
	MaxGapRatio float64
}

// StandardNormalizer는 Normalizer 인터페이스의 기본 구현체입니다.
type StandardNormalizer struct {
	config NormalizationConfig
}

// NewStandardNormalizer는 새로운 StandardNormalizer 인스턴스를 생성합니다.
func NewStandardNormalizer(config NormalizationConfig) *StandardNormalizer {
	// 기본 설정 적용
	if config.DefaultTimezone == nil {
		config.DefaultTimezone = time.UTC
	}

	if config.DefaultCurrency == "" {
		config.DefaultCurrency = "USD"
	}

	if config.InterpolationMethod == "" {
		config.InterpolationMethod = "linear"
	}

	if config.PriceScaleFactor == 0 {
		config.PriceScaleFactor = 1.0
	}

	if config.VolumeScaleFactor == 0 {
		config.VolumeScaleFactor = 1.0
	}

	if config.MaxGapRatio == 0 {
		config.MaxGapRatio = 0.1 // 10% 이상의 갭은 보간하지 않음
	}

	return &StandardNormalizer{
		config: config,
	}
}

// NormalizeHistoricalData는 과거 가격 데이터를 정규화합니다.
func (n *StandardNormalizer) NormalizeHistoricalData(response *source.HistoricalDataResponse) (*source.HistoricalDataResponse, error) {
	if response == nil {
		return nil, fmt.Errorf("input response is nil")
	}

	// 원본 데이터 복사
	normalizedResponse := &source.HistoricalDataResponse{
		Symbol:    response.Symbol,
		AssetType: response.AssetType,
		Interval:  response.Interval,
		Data:      make([]source.PriceData, len(response.Data)),
	}

	// 데이터 정규화
	for i, data := range response.Data {
		// 시간대 정규화
		normalizedTime := data.Timestamp.In(n.config.DefaultTimezone)

		// 가격 데이터 정규화
		normalizedResponse.Data[i] = source.PriceData{
			Timestamp:     normalizedTime,
			Open:          n.normalizePrice(data.Open),
			High:          n.normalizePrice(data.High),
			Low:           n.normalizePrice(data.Low),
			Close:         n.normalizePrice(data.Close),
			AdjustedClose: n.normalizePrice(data.AdjustedClose),
			Volume:        n.normalizeVolume(data.Volume),
		}
	}

	// 정렬 확인
	n.ensureChronologicalOrder(normalizedResponse.Data)

	// 누락된 데이터 처리
	err := n.interpolateMissingData(normalizedResponse)
	if err != nil {
		return nil, fmt.Errorf("interpolation error: %w", err)
	}

	return normalizedResponse, nil
}

// NormalizeRealTimeData는 실시간 가격 데이터를 정규화합니다.
func (n *StandardNormalizer) NormalizeRealTimeData(response *source.RealTimeDataResponse) (*source.RealTimeDataResponse, error) {
	if response == nil {
		return nil, fmt.Errorf("input response is nil")
	}

	// 시간대 정규화
	normalizedTime := response.Timestamp.In(n.config.DefaultTimezone)

	// 가격 데이터 정규화
	normalizedResponse := &source.RealTimeDataResponse{
		Symbol:        response.Symbol,
		AssetType:     response.AssetType,
		CurrentPrice:  n.normalizePrice(response.CurrentPrice),
		Timestamp:     normalizedTime,
		Change:        n.normalizePrice(response.Change),
		ChangePercent: response.ChangePercent, // 퍼센트는 정규화하지 않음
		Volume:        n.normalizeVolume(response.Volume),
		MarketCap:     n.normalizePrice(response.MarketCap),
		High24h:       n.normalizePrice(response.High24h),
		Low24h:        n.normalizePrice(response.Low24h),
	}

	return normalizedResponse, nil
}

// NormalizeMetadata는 자산 메타데이터를 정규화합니다.
func (n *StandardNormalizer) NormalizeMetadata(response *source.MetadataResponse) (*source.MetadataResponse, error) {
	if response == nil {
		return nil, fmt.Errorf("input response is nil")
	}

	// 메타데이터 복사
	normalizedResponse := &source.MetadataResponse{
		Symbol:      response.Symbol,
		AssetType:   response.AssetType,
		Name:        sanitizeString(response.Name),
		Exchange:    sanitizeString(response.Exchange),
		Currency:    n.normalizeCurrency(response.Currency),
		Country:     sanitizeString(response.Country),
		Description: sanitizeString(response.Description),
		Sector:      sanitizeString(response.Sector),
		Industry:    sanitizeString(response.Industry),
		Website:     sanitizeString(response.Website),
		LogoURL:     sanitizeString(response.LogoURL),
		LastUpdated: response.LastUpdated.In(n.config.DefaultTimezone),
	}

	return normalizedResponse, nil
}

// 가격 정규화 함수
func (n *StandardNormalizer) normalizePrice(price float64) float64 {
	if price == 0 {
		return 0
	}

	// 스케일 적용
	price = price * n.config.PriceScaleFactor

	// NaN 및 Inf 값 처리
	if math.IsNaN(price) || math.IsInf(price, 0) {
		return 0
	}

	return price
}

// 거래량 정규화 함수
func (n *StandardNormalizer) normalizeVolume(volume int64) int64 {
	if volume == 0 {
		return 0
	}

	// 스케일 적용
	normalizedVolume := float64(volume) * n.config.VolumeScaleFactor

	// NaN 및 Inf 값 처리
	if math.IsNaN(normalizedVolume) || math.IsInf(normalizedVolume, 0) {
		return 0
	}

	return int64(normalizedVolume)
}

// 통화 정규화 함수
func (n *StandardNormalizer) normalizeCurrency(currency string) string {
	if currency == "" {
		return n.config.DefaultCurrency
	}
	return currency
}

// 문자열 정리 함수
func sanitizeString(s string) string {
	if s == "" {
		return ""
	}
	// 여기에 필요한 문자열 정리 로직을 추가할 수 있습니다.
	return s
}

// 시간순 정렬 확인 함수
func (n *StandardNormalizer) ensureChronologicalOrder(data []source.PriceData) {
	// 데이터가 이미 시간순으로 정렬되어 있다고 가정합니다.
	// 필요한 경우 여기에 정렬 로직을 추가할 수 있습니다.
}

// 누락된 데이터 보간 함수
func (n *StandardNormalizer) interpolateMissingData(response *source.HistoricalDataResponse) error {
	// 보간이 비활성화된 경우
	if n.config.InterpolationMethod == "none" || len(response.Data) < 2 {
		return nil
	}

	// 여기에 보간 로직 구현
	// 현재는 단순한 예시만 포함합니다.
	// 실제 구현에서는 보간 방법에 따라 더 복잡한 로직이 필요할 수 있습니다.
	for i := 1; i < len(response.Data); i++ {
		current := response.Data[i]
		previous := response.Data[i-1]

		// 갭 크기 계산
		gapRatio := calculateGapRatio(current, previous)

		// 갭이 너무 크면 보간하지 않음
		if gapRatio > n.config.MaxGapRatio {
			continue
		}

		// 0 값 확인 및 보간
		if current.Open == 0 && previous.Open != 0 {
			response.Data[i].Open = previous.Open
		}
		if current.High == 0 && previous.High != 0 {
			response.Data[i].High = previous.High
		}
		if current.Low == 0 && previous.Low != 0 {
			response.Data[i].Low = previous.Low
		}
		if current.Close == 0 && previous.Close != 0 {
			response.Data[i].Close = previous.Close
		}
		if current.AdjustedClose == 0 && previous.AdjustedClose != 0 {
			response.Data[i].AdjustedClose = previous.AdjustedClose
		}
		if current.Volume == 0 && previous.Volume != 0 {
			response.Data[i].Volume = previous.Volume
		}
	}

	return nil
}

// 갭 비율 계산 함수
func calculateGapRatio(current, previous source.PriceData) float64 {
	if previous.Close == 0 {
		return math.MaxFloat64
	}

	// 가격 차이의 절대값을 이전 가격으로 나눔
	return math.Abs(current.Close-previous.Close) / previous.Close
}
