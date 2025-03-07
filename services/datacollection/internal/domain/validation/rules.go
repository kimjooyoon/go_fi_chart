package validation

import (
	"fmt"
	"math"
	"time"

	"github.com/kimjooyoon/go_fi_chart/services/datacollection/internal/domain/source"
)

// 시간 간격 검증 규칙
type TimeIntervalRule struct {
	name        string
	description string
	severity    ValidationSeverity
	maxGap      time.Duration
}

// NewTimeIntervalRule은 시간 간격 검증 규칙을 생성합니다.
func NewTimeIntervalRule(maxGap time.Duration, severity ValidationSeverity) *TimeIntervalRule {
	return &TimeIntervalRule{
		name:        "TimeIntervalRule",
		description: "시간 간격이 일정 기준을 초과하는지 검증",
		severity:    severity,
		maxGap:      maxGap,
	}
}

func (r *TimeIntervalRule) GetName() string {
	return r.name
}

func (r *TimeIntervalRule) GetDescription() string {
	return r.description
}

func (r *TimeIntervalRule) GetSeverity() ValidationSeverity {
	return r.severity
}

func (r *TimeIntervalRule) ValidateHistoricalData(response *source.HistoricalDataResponse) (*ValidationResult, error) {
	if response == nil || len(response.Data) < 2 {
		return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
	}

	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// 데이터가 시간순으로 정렬되어 있다고 가정
	for i := 1; i < len(response.Data); i++ {
		current := response.Data[i].Timestamp
		previous := response.Data[i-1].Timestamp

		// 이전 데이터가 현재 데이터보다 미래인 경우 (시간 역전)
		if previous.After(current) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d].Timestamp", i),
				Code:    "time_reversal",
				Message: fmt.Sprintf("Time reversal detected: %v is before %v", current, previous),
				Value:   map[string]time.Time{"Current": current, "Previous": previous},
			})
			continue
		}

		// 시간 간격 확인
		gap := current.Sub(previous)
		if gap > r.maxGap {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d-%d].Timestamp", i-1, i),
				Code:    "time_gap",
				Message: fmt.Sprintf("Time gap too large: %v (max allowed: %v)", gap, r.maxGap),
				Value:   map[string]interface{}{"Gap": gap.String(), "MaxAllowed": r.maxGap.String()},
			})
		}
	}

	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result, nil
}

func (r *TimeIntervalRule) ValidateRealTimeData(response *source.RealTimeDataResponse) (*ValidationResult, error) {
	// 실시간 데이터에는 시간 간격 개념이 없으므로 항상 유효
	return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
}

func (r *TimeIntervalRule) ValidateMetadata(response *source.MetadataResponse) (*ValidationResult, error) {
	// 메타데이터에는 시간 간격 개념이 없으므로 항상 유효
	return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
}

// 가격 변동성 검증 규칙
type PriceVolatilityRule struct {
	name             string
	description      string
	severity         ValidationSeverity
	maxVolatility    float64
	lookbackPeriod   int
	volatilityMetric string // "range", "stddev" 등
}

// NewPriceVolatilityRule은 가격 변동성 검증 규칙을 생성합니다.
func NewPriceVolatilityRule(maxVolatility float64, lookbackPeriod int, volatilityMetric string, severity ValidationSeverity) *PriceVolatilityRule {
	return &PriceVolatilityRule{
		name:             "PriceVolatilityRule",
		description:      "가격 변동성이 일정 기준을 초과하는지 검증",
		severity:         severity,
		maxVolatility:    maxVolatility,
		lookbackPeriod:   lookbackPeriod,
		volatilityMetric: volatilityMetric,
	}
}

func (r *PriceVolatilityRule) GetName() string {
	return r.name
}

func (r *PriceVolatilityRule) GetDescription() string {
	return r.description
}

func (r *PriceVolatilityRule) GetSeverity() ValidationSeverity {
	return r.severity
}

func (r *PriceVolatilityRule) ValidateHistoricalData(response *source.HistoricalDataResponse) (*ValidationResult, error) {
	if response == nil || len(response.Data) < r.lookbackPeriod {
		return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
	}

	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// 각 포인트별로 전/후 구간의 변동성 검사
	for i := r.lookbackPeriod; i < len(response.Data); i++ {
		// 룩백 데이터 슬라이스
		lookbackData := response.Data[i-r.lookbackPeriod : i]

		// 변동성 계산
		var volatility float64

		switch r.volatilityMetric {
		case "range":
			volatility = calculatePriceRange(lookbackData)
		case "stddev":
			volatility = calculatePriceStdDev(lookbackData)
		default:
			volatility = calculatePriceRange(lookbackData)
		}

		// 기준 가격으로 마지막 종가 사용
		basePrice := lookbackData[len(lookbackData)-1].Close
		if basePrice == 0 {
			continue // 분모가 0이면 건너뜀
		}

		// 변동성 비율 계산
		volatilityRatio := volatility / basePrice

		// 검증
		if volatilityRatio > r.maxVolatility {
			result.Errors = append(result.Errors, ValidationError{
				Field: fmt.Sprintf("Data[%d-%d]", i-r.lookbackPeriod, i-1),
				Code:  "high_volatility",
				Message: fmt.Sprintf("High price volatility detected: %.2f%% (max allowed: %.2f%%)",
					volatilityRatio*100, r.maxVolatility*100),
				Value: map[string]interface{}{"VolatilityRatio": volatilityRatio, "MaxAllowed": r.maxVolatility},
			})
		}
	}

	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result, nil
}

func (r *PriceVolatilityRule) ValidateRealTimeData(response *source.RealTimeDataResponse) (*ValidationResult, error) {
	// 실시간 데이터에서는 고가/저가 차이로 간단히 변동성 검증
	if response == nil {
		return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
	}

	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// 24시간 고가/저가 차이로 변동성 계산
	if response.Low24h > 0 && response.High24h > 0 {
		priceRange := response.High24h - response.Low24h
		volatilityRatio := priceRange / response.Low24h

		if volatilityRatio > r.maxVolatility {
			result.Errors = append(result.Errors, ValidationError{
				Field: "High24h/Low24h",
				Code:  "high_volatility",
				Message: fmt.Sprintf("High 24h price volatility detected: %.2f%% (max allowed: %.2f%%)",
					volatilityRatio*100, r.maxVolatility*100),
				Value: map[string]interface{}{"VolatilityRatio": volatilityRatio, "MaxAllowed": r.maxVolatility},
			})
			result.IsValid = false
		}
	}

	return result, nil
}

func (r *PriceVolatilityRule) ValidateMetadata(response *source.MetadataResponse) (*ValidationResult, error) {
	// 메타데이터에는 가격 변동성 개념이 없으므로 항상 유효
	return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
}

// 거래량 급증 검증 규칙
type VolumeAnomalyRule struct {
	name           string
	description    string
	severity       ValidationSeverity
	maxMultiplier  float64
	lookbackPeriod int
}

// NewVolumeAnomalyRule은 거래량 급증 검증 규칙을 생성합니다.
func NewVolumeAnomalyRule(maxMultiplier float64, lookbackPeriod int, severity ValidationSeverity) *VolumeAnomalyRule {
	return &VolumeAnomalyRule{
		name:           "VolumeAnomalyRule",
		description:    "거래량이 평균보다 일정 배수 이상 증가했는지 검증",
		severity:       severity,
		maxMultiplier:  maxMultiplier,
		lookbackPeriod: lookbackPeriod,
	}
}

func (r *VolumeAnomalyRule) GetName() string {
	return r.name
}

func (r *VolumeAnomalyRule) GetDescription() string {
	return r.description
}

func (r *VolumeAnomalyRule) GetSeverity() ValidationSeverity {
	return r.severity
}

func (r *VolumeAnomalyRule) ValidateHistoricalData(response *source.HistoricalDataResponse) (*ValidationResult, error) {
	if response == nil || len(response.Data) < r.lookbackPeriod+1 {
		return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
	}

	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// 각 포인트별로 이전 구간 평균 거래량과 비교
	for i := r.lookbackPeriod; i < len(response.Data); i++ {
		currentVolume := response.Data[i].Volume

		// 이전 기간 평균 거래량 계산
		var totalVolume int64
		for j := i - r.lookbackPeriod; j < i; j++ {
			totalVolume += response.Data[j].Volume
		}

		avgVolume := float64(totalVolume) / float64(r.lookbackPeriod)
		if avgVolume == 0 {
			continue // 분모가 0이면 건너뜀
		}

		// 현재 거래량과 평균 거래량의 비율
		multiplier := float64(currentVolume) / avgVolume

		// 검증
		if multiplier > r.maxMultiplier {
			result.Errors = append(result.Errors, ValidationError{
				Field: fmt.Sprintf("Data[%d].Volume", i),
				Code:  "volume_anomaly",
				Message: fmt.Sprintf("Volume anomaly detected: %.2fx increase over average (max allowed: %.2fx)",
					multiplier, r.maxMultiplier),
				Value: map[string]interface{}{"Multiplier": multiplier, "MaxAllowed": r.maxMultiplier},
			})
		}
	}

	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result, nil
}

func (r *VolumeAnomalyRule) ValidateRealTimeData(response *source.RealTimeDataResponse) (*ValidationResult, error) {
	// 실시간 데이터에서는 룩백 데이터가 없으므로 검증 불가
	return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
}

func (r *VolumeAnomalyRule) ValidateMetadata(response *source.MetadataResponse) (*ValidationResult, error) {
	// 메타데이터에는 거래량 개념이 없으므로 항상 유효
	return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
}

// 메타데이터 일관성 검증 규칙
type MetadataConsistencyRule struct {
	name        string
	description string
	severity    ValidationSeverity
	fieldMap    map[string]map[string]bool // AssetType별 필수 필드 맵
}

// NewMetadataConsistencyRule은 메타데이터 일관성 검증 규칙을 생성합니다.
func NewMetadataConsistencyRule(severity ValidationSeverity) *MetadataConsistencyRule {
	fieldMap := map[string]map[string]bool{
		"stock": {
			"Symbol":      true,
			"Name":        true,
			"Exchange":    true,
			"Currency":    true,
			"Country":     true,
			"Sector":      true,
			"Industry":    true,
			"Description": false,
			"Website":     false,
			"LogoURL":     false,
		},
		"etf": {
			"Symbol":      true,
			"Name":        true,
			"Exchange":    true,
			"Currency":    true,
			"Country":     false,
			"Sector":      false,
			"Industry":    false,
			"Description": true,
			"Website":     false,
			"LogoURL":     false,
		},
		"crypto": {
			"Symbol":      true,
			"Name":        true,
			"Exchange":    true,
			"Currency":    true,
			"Country":     false,
			"Sector":      false,
			"Industry":    false,
			"Description": false,
			"Website":     false,
			"LogoURL":     false,
		},
		"forex": {
			"Symbol":      true,
			"Name":        true,
			"Exchange":    false,
			"Currency":    true,
			"Country":     false,
			"Sector":      false,
			"Industry":    false,
			"Description": false,
			"Website":     false,
			"LogoURL":     false,
		},
	}

	return &MetadataConsistencyRule{
		name:        "MetadataConsistencyRule",
		description: "자산 유형에 따른 메타데이터 필드 일관성 검증",
		severity:    severity,
		fieldMap:    fieldMap,
	}
}

func (r *MetadataConsistencyRule) GetName() string {
	return r.name
}

func (r *MetadataConsistencyRule) GetDescription() string {
	return r.description
}

func (r *MetadataConsistencyRule) GetSeverity() ValidationSeverity {
	return r.severity
}

func (r *MetadataConsistencyRule) ValidateHistoricalData(response *source.HistoricalDataResponse) (*ValidationResult, error) {
	// 과거 데이터에는 메타데이터 일관성 검증을 적용하지 않음
	return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
}

func (r *MetadataConsistencyRule) ValidateRealTimeData(response *source.RealTimeDataResponse) (*ValidationResult, error) {
	// 실시간 데이터에는 메타데이터 일관성 검증을 적용하지 않음
	return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
}

func (r *MetadataConsistencyRule) ValidateMetadata(response *source.MetadataResponse) (*ValidationResult, error) {
	if response == nil {
		return &ValidationResult{IsValid: true, Errors: []ValidationError{}}, nil
	}

	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// 자산 유형 확인
	assetTypeStr := string(response.AssetType)
	assetTypeStr = normalizeAssetType(assetTypeStr)

	// 정의된 자산 유형인지 확인
	requiredFields, ok := r.fieldMap[assetTypeStr]
	if !ok {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "AssetType",
			Code:    "unknown",
			Message: fmt.Sprintf("Unknown asset type: %s", assetTypeStr),
			Value:   response.AssetType,
		})
		result.IsValid = false
		return result, nil
	}

	// 자산 유형에 따른 필수 필드 확인
	for field, required := range requiredFields {
		if !required {
			continue
		}

		switch field {
		case "Symbol":
			if response.Symbol == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Symbol",
					Code:    "required_for_asset_type",
					Message: fmt.Sprintf("Symbol is required for asset type: %s", assetTypeStr),
					Value:   response.Symbol,
				})
			}
		case "Name":
			if response.Name == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Name",
					Code:    "required_for_asset_type",
					Message: fmt.Sprintf("Name is required for asset type: %s", assetTypeStr),
					Value:   response.Name,
				})
			}
		case "Exchange":
			if response.Exchange == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Exchange",
					Code:    "required_for_asset_type",
					Message: fmt.Sprintf("Exchange is required for asset type: %s", assetTypeStr),
					Value:   response.Exchange,
				})
			}
		case "Currency":
			if response.Currency == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Currency",
					Code:    "required_for_asset_type",
					Message: fmt.Sprintf("Currency is required for asset type: %s", assetTypeStr),
					Value:   response.Currency,
				})
			}
		case "Country":
			if response.Country == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Country",
					Code:    "required_for_asset_type",
					Message: fmt.Sprintf("Country is required for asset type: %s", assetTypeStr),
					Value:   response.Country,
				})
			}
		case "Sector":
			if response.Sector == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Sector",
					Code:    "required_for_asset_type",
					Message: fmt.Sprintf("Sector is required for asset type: %s", assetTypeStr),
					Value:   response.Sector,
				})
			}
		case "Industry":
			if response.Industry == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Industry",
					Code:    "required_for_asset_type",
					Message: fmt.Sprintf("Industry is required for asset type: %s", assetTypeStr),
					Value:   response.Industry,
				})
			}
		case "Description":
			if response.Description == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Description",
					Code:    "required_for_asset_type",
					Message: fmt.Sprintf("Description is required for asset type: %s", assetTypeStr),
					Value:   response.Description,
				})
			}
		case "Website":
			if response.Website == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Website",
					Code:    "required_for_asset_type",
					Message: fmt.Sprintf("Website is required for asset type: %s", assetTypeStr),
					Value:   response.Website,
				})
			}
		case "LogoURL":
			if response.LogoURL == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "LogoURL",
					Code:    "required_for_asset_type",
					Message: fmt.Sprintf("LogoURL is required for asset type: %s", assetTypeStr),
					Value:   response.LogoURL,
				})
			}
		}
	}

	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result, nil
}

// 자산 유형 정규화 함수
func normalizeAssetType(assetType string) string {
	switch assetType {
	case "stock", "stocks", "equity", "equities":
		return "stock"
	case "etf", "ETF", "exchange-traded fund", "exchange_traded_fund":
		return "etf"
	case "crypto", "cryptocurrency", "digital asset", "digital_asset", "digital currency", "digital_currency":
		return "crypto"
	case "forex", "currency", "fx", "foreign exchange", "foreign_exchange":
		return "forex"
	default:
		return assetType
	}
}

// 헬퍼 함수들

// 가격 범위 계산
func calculatePriceRange(data []source.PriceData) float64 {
	if len(data) == 0 {
		return 0
	}

	var max, min float64 = -math.MaxFloat64, math.MaxFloat64

	for _, price := range data {
		if price.High > max {
			max = price.High
		}
		if price.Low < min || min == math.MaxFloat64 {
			min = price.Low
		}
	}

	if min == math.MaxFloat64 {
		min = 0
	}

	if max == -math.MaxFloat64 {
		max = 0
	}

	return max - min
}

// 종가 표준편차 계산
func calculatePriceStdDev(data []source.PriceData) float64 {
	if len(data) < 2 {
		return 0
	}

	// 평균 계산
	var sum float64
	for _, price := range data {
		sum += price.Close
	}
	mean := sum / float64(len(data))

	// 분산 계산
	var variance float64
	for _, price := range data {
		diff := price.Close - mean
		variance += diff * diff
	}
	variance /= float64(len(data) - 1)

	// 표준편차 반환
	return math.Sqrt(variance)
}
