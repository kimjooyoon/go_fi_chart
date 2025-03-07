package validation

import (
	"fmt"
	"time"

	"github.com/kimjooyoon/go_fi_chart/services/datacollection/internal/domain/source"
)

// ValidationResult는 검증 결과를 나타냅니다.
type ValidationResult struct {
	// IsValid는 검증 통과 여부를 나타냅니다.
	IsValid bool

	// Errors는 검증 실패 시 발생한 에러 목록입니다.
	Errors []ValidationError
}

// ValidationError는 검증 오류 정보를 나타냅니다.
type ValidationError struct {
	// Field는 오류가 발생한 필드명입니다.
	Field string

	// Code는 오류 코드입니다.
	Code string

	// Message는 오류 메시지입니다.
	Message string

	// Value는 오류가 발생한 값입니다.
	Value interface{}
}

// ValidationSeverity는 검증 오류의 심각도를 나타냅니다.
type ValidationSeverity string

const (
	// SeverityError는 치명적인 오류를 나타냅니다. 데이터를 사용할 수 없습니다.
	SeverityError ValidationSeverity = "error"

	// SeverityWarning은 경고를 나타냅니다. 데이터를 사용할 수 있지만 주의가 필요합니다.
	SeverityWarning ValidationSeverity = "warning"

	// SeverityInfo는 정보성 메시지를 나타냅니다. 데이터를 정상적으로 사용할 수 있습니다.
	SeverityInfo ValidationSeverity = "info"
)

// Validator는 데이터 유효성 검증을 위한 인터페이스를 정의합니다.
type Validator interface {
	// ValidateHistoricalData는 과거 가격 데이터의 유효성을 검증합니다.
	ValidateHistoricalData(response *source.HistoricalDataResponse) (*ValidationResult, error)

	// ValidateRealTimeData는 실시간 가격 데이터의 유효성을 검증합니다.
	ValidateRealTimeData(response *source.RealTimeDataResponse) (*ValidationResult, error)

	// ValidateMetadata는 자산 메타데이터의 유효성을 검증합니다.
	ValidateMetadata(response *source.MetadataResponse) (*ValidationResult, error)

	// AddRule은 검증 규칙을 추가합니다.
	AddRule(rule ValidationRule)
}

// ValidationRule은 데이터 유효성 검증 규칙을 정의하는 인터페이스입니다.
type ValidationRule interface {
	// GetName은 규칙의 이름을 반환합니다.
	GetName() string

	// GetDescription은 규칙에 대한 설명을 반환합니다.
	GetDescription() string

	// GetSeverity는 규칙 위반 시 심각도를 반환합니다.
	GetSeverity() ValidationSeverity

	// ValidateHistoricalData는 과거 가격 데이터에 대한 규칙을 검증합니다.
	ValidateHistoricalData(response *source.HistoricalDataResponse) (*ValidationResult, error)

	// ValidateRealTimeData는 실시간 가격 데이터에 대한 규칙을 검증합니다.
	ValidateRealTimeData(response *source.RealTimeDataResponse) (*ValidationResult, error)

	// ValidateMetadata는 자산 메타데이터에 대한 규칙을 검증합니다.
	ValidateMetadata(response *source.MetadataResponse) (*ValidationResult, error)
}

// ValidationConfig는 유효성 검증 설정을 정의합니다.
type ValidationConfig struct {
	// MaxPriceThreshold는 허용 가능한 최대 가격입니다.
	MaxPriceThreshold float64

	// MinPriceThreshold는 허용 가능한 최소 가격입니다.
	MinPriceThreshold float64

	// MaxVolumeThreshold는 허용 가능한 최대 거래량입니다.
	MaxVolumeThreshold int64

	// MaxDataAge는 실시간 데이터의 최대 허용 나이(기간)입니다.
	MaxDataAge time.Duration

	// RequiredMetadataFields는 필수 메타데이터 필드 목록입니다.
	RequiredMetadataFields []string

	// StrictValidation은 엄격한 검증 모드 활성화 여부입니다.
	StrictValidation bool
}

// StandardValidator는 Validator 인터페이스의 기본 구현체입니다.
type StandardValidator struct {
	config ValidationConfig
	rules  []ValidationRule
}

// NewStandardValidator는 새로운 StandardValidator 인스턴스를 생성합니다.
func NewStandardValidator(config ValidationConfig) *StandardValidator {
	// 기본 설정 적용
	if config.MaxPriceThreshold == 0 {
		config.MaxPriceThreshold = 1000000.0 // 백만 달러
	}

	if config.MinPriceThreshold == 0 {
		config.MinPriceThreshold = 0.000001 // 마이크로 센트
	}

	if config.MaxVolumeThreshold == 0 {
		config.MaxVolumeThreshold = 1000000000000 // 1조
	}

	if config.MaxDataAge == 0 {
		config.MaxDataAge = 5 * time.Minute // 5분
	}

	if len(config.RequiredMetadataFields) == 0 {
		config.RequiredMetadataFields = []string{
			"Symbol",
			"Name",
			"AssetType",
			"Exchange",
		}
	}

	return &StandardValidator{
		config: config,
		rules:  []ValidationRule{},
	}
}

// ValidateHistoricalData는 과거 가격 데이터의 유효성을 검증합니다.
func (v *StandardValidator) ValidateHistoricalData(response *source.HistoricalDataResponse) (*ValidationResult, error) {
	if response == nil {
		return nil, fmt.Errorf("input response is nil")
	}

	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// 기본 검증
	v.validateHistoricalDataBasic(response, result)

	// 규칙 기반 검증
	for _, rule := range v.rules {
		ruleResult, err := rule.ValidateHistoricalData(response)
		if err != nil {
			return nil, fmt.Errorf("rule '%s' validation error: %w", rule.GetName(), err)
		}

		if ruleResult != nil && len(ruleResult.Errors) > 0 {
			result.Errors = append(result.Errors, ruleResult.Errors...)

			// 엄격한 검증 모드에서 오류 발생 시 즉시 false 설정
			if v.config.StrictValidation && rule.GetSeverity() == SeverityError {
				result.IsValid = false
			}
		}
	}

	// 오류가 하나라도 있으면 유효하지 않음
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result, nil
}

// ValidateRealTimeData는 실시간 가격 데이터의 유효성을 검증합니다.
func (v *StandardValidator) ValidateRealTimeData(response *source.RealTimeDataResponse) (*ValidationResult, error) {
	if response == nil {
		return nil, fmt.Errorf("input response is nil")
	}

	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// 기본 검증
	v.validateRealTimeDataBasic(response, result)

	// 규칙 기반 검증
	for _, rule := range v.rules {
		ruleResult, err := rule.ValidateRealTimeData(response)
		if err != nil {
			return nil, fmt.Errorf("rule '%s' validation error: %w", rule.GetName(), err)
		}

		if ruleResult != nil && len(ruleResult.Errors) > 0 {
			result.Errors = append(result.Errors, ruleResult.Errors...)

			// 엄격한 검증 모드에서 오류 발생 시 즉시 false 설정
			if v.config.StrictValidation && rule.GetSeverity() == SeverityError {
				result.IsValid = false
			}
		}
	}

	// 오류가 하나라도 있으면 유효하지 않음
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result, nil
}

// ValidateMetadata는 자산 메타데이터의 유효성을 검증합니다.
func (v *StandardValidator) ValidateMetadata(response *source.MetadataResponse) (*ValidationResult, error) {
	if response == nil {
		return nil, fmt.Errorf("input response is nil")
	}

	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// 기본 검증
	v.validateMetadataBasic(response, result)

	// 규칙 기반 검증
	for _, rule := range v.rules {
		ruleResult, err := rule.ValidateMetadata(response)
		if err != nil {
			return nil, fmt.Errorf("rule '%s' validation error: %w", rule.GetName(), err)
		}

		if ruleResult != nil && len(ruleResult.Errors) > 0 {
			result.Errors = append(result.Errors, ruleResult.Errors...)

			// 엄격한 검증 모드에서 오류 발생 시 즉시 false 설정
			if v.config.StrictValidation && rule.GetSeverity() == SeverityError {
				result.IsValid = false
			}
		}
	}

	// 오류가 하나라도 있으면 유효하지 않음
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	return result, nil
}

// AddRule은 검증 규칙을 추가합니다.
func (v *StandardValidator) AddRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

// 과거 가격 데이터 기본 검증
func (v *StandardValidator) validateHistoricalDataBasic(response *source.HistoricalDataResponse, result *ValidationResult) {
	// 심볼 검증
	if response.Symbol == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Symbol",
			Code:    "required",
			Message: "Symbol is required",
			Value:   response.Symbol,
		})
	}

	// 자산 유형 검증
	if response.AssetType == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "AssetType",
			Code:    "required",
			Message: "AssetType is required",
			Value:   response.AssetType,
		})
	}

	// 데이터 존재 여부 검증
	if len(response.Data) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Data",
			Code:    "empty",
			Message: "Historical data is empty",
			Value:   response.Data,
		})
	}

	// 각 가격 데이터 검증
	for i, data := range response.Data {
		// 타임스탬프 검증
		if data.Timestamp.IsZero() {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d].Timestamp", i),
				Code:    "invalid",
				Message: "Timestamp is zero",
				Value:   data.Timestamp,
			})
		}

		// 가격 검증
		if data.Open < v.config.MinPriceThreshold || data.Open > v.config.MaxPriceThreshold {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d].Open", i),
				Code:    "range",
				Message: fmt.Sprintf("Open price must be between %f and %f", v.config.MinPriceThreshold, v.config.MaxPriceThreshold),
				Value:   data.Open,
			})
		}

		if data.High < v.config.MinPriceThreshold || data.High > v.config.MaxPriceThreshold {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d].High", i),
				Code:    "range",
				Message: fmt.Sprintf("High price must be between %f and %f", v.config.MinPriceThreshold, v.config.MaxPriceThreshold),
				Value:   data.High,
			})
		}

		if data.Low < v.config.MinPriceThreshold || data.Low > v.config.MaxPriceThreshold {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d].Low", i),
				Code:    "range",
				Message: fmt.Sprintf("Low price must be between %f and %f", v.config.MinPriceThreshold, v.config.MaxPriceThreshold),
				Value:   data.Low,
			})
		}

		if data.Close < v.config.MinPriceThreshold || data.Close > v.config.MaxPriceThreshold {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d].Close", i),
				Code:    "range",
				Message: fmt.Sprintf("Close price must be between %f and %f", v.config.MinPriceThreshold, v.config.MaxPriceThreshold),
				Value:   data.Close,
			})
		}

		// 거래량 검증
		if data.Volume < 0 || data.Volume > v.config.MaxVolumeThreshold {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d].Volume", i),
				Code:    "range",
				Message: fmt.Sprintf("Volume must be between 0 and %d", v.config.MaxVolumeThreshold),
				Value:   data.Volume,
			})
		}

		// OHLC 논리 검증
		if data.Low > data.High {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d].Low/High", i),
				Code:    "logic",
				Message: "Low price cannot be greater than High price",
				Value:   map[string]float64{"Low": data.Low, "High": data.High},
			})
		}

		if data.Open < data.Low || data.Open > data.High {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d].Open", i),
				Code:    "logic",
				Message: "Open price must be between Low and High prices",
				Value:   map[string]float64{"Open": data.Open, "Low": data.Low, "High": data.High},
			})
		}

		if data.Close < data.Low || data.Close > data.High {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("Data[%d].Close", i),
				Code:    "logic",
				Message: "Close price must be between Low and High prices",
				Value:   map[string]float64{"Close": data.Close, "Low": data.Low, "High": data.High},
			})
		}
	}
}

// 실시간 가격 데이터 기본 검증
func (v *StandardValidator) validateRealTimeDataBasic(response *source.RealTimeDataResponse, result *ValidationResult) {
	// 심볼 검증
	if response.Symbol == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Symbol",
			Code:    "required",
			Message: "Symbol is required",
			Value:   response.Symbol,
		})
	}

	// 자산 유형 검증
	if response.AssetType == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "AssetType",
			Code:    "required",
			Message: "AssetType is required",
			Value:   response.AssetType,
		})
	}

	// 타임스탬프 검증
	if response.Timestamp.IsZero() {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Timestamp",
			Code:    "invalid",
			Message: "Timestamp is zero",
			Value:   response.Timestamp,
		})
	}

	// 데이터 최신성 검증
	dataAge := time.Since(response.Timestamp)
	if dataAge > v.config.MaxDataAge {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Timestamp",
			Code:    "stale",
			Message: fmt.Sprintf("Data is too old (age: %v, max allowed: %v)", dataAge, v.config.MaxDataAge),
			Value:   response.Timestamp,
		})
	}

	// 가격 검증
	if response.CurrentPrice < v.config.MinPriceThreshold || response.CurrentPrice > v.config.MaxPriceThreshold {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "CurrentPrice",
			Code:    "range",
			Message: fmt.Sprintf("Current price must be between %f and %f", v.config.MinPriceThreshold, v.config.MaxPriceThreshold),
			Value:   response.CurrentPrice,
		})
	}

	// 거래량 검증
	if response.Volume < 0 || response.Volume > v.config.MaxVolumeThreshold {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Volume",
			Code:    "range",
			Message: fmt.Sprintf("Volume must be between 0 and %d", v.config.MaxVolumeThreshold),
			Value:   response.Volume,
		})
	}

	// 24시간 고가/저가 검증
	if response.Low24h > response.High24h {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Low24h/High24h",
			Code:    "logic",
			Message: "24h Low price cannot be greater than 24h High price",
			Value:   map[string]float64{"Low24h": response.Low24h, "High24h": response.High24h},
		})
	}

	// 현재 가격과 24시간 고가/저가 논리 검증
	if response.CurrentPrice < response.Low24h || response.CurrentPrice > response.High24h {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "CurrentPrice",
			Code:    "logic",
			Message: "Current price should be between 24h Low and High prices",
			Value:   map[string]float64{"CurrentPrice": response.CurrentPrice, "Low24h": response.Low24h, "High24h": response.High24h},
		})
	}
}

// 자산 메타데이터 기본 검증
func (v *StandardValidator) validateMetadataBasic(response *source.MetadataResponse, result *ValidationResult) {
	// 필수 필드 검증
	for _, field := range v.config.RequiredMetadataFields {
		switch field {
		case "Symbol":
			if response.Symbol == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Symbol",
					Code:    "required",
					Message: "Symbol is required",
					Value:   response.Symbol,
				})
			}
		case "Name":
			if response.Name == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Name",
					Code:    "required",
					Message: "Name is required",
					Value:   response.Name,
				})
			}
		case "AssetType":
			if response.AssetType == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "AssetType",
					Code:    "required",
					Message: "AssetType is required",
					Value:   response.AssetType,
				})
			}
		case "Exchange":
			if response.Exchange == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Exchange",
					Code:    "required",
					Message: "Exchange is required",
					Value:   response.Exchange,
				})
			}
		case "Currency":
			if response.Currency == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Currency",
					Code:    "required",
					Message: "Currency is required",
					Value:   response.Currency,
				})
			}
		case "Country":
			if response.Country == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Country",
					Code:    "required",
					Message: "Country is required",
					Value:   response.Country,
				})
			}
		case "Sector":
			if response.Sector == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Sector",
					Code:    "required",
					Message: "Sector is required",
					Value:   response.Sector,
				})
			}
		case "Industry":
			if response.Industry == "" {
				result.Errors = append(result.Errors, ValidationError{
					Field:   "Industry",
					Code:    "required",
					Message: "Industry is required",
					Value:   response.Industry,
				})
			}
		}
	}

	// 마지막 업데이트 타임스탬프 검증
	if response.LastUpdated.IsZero() {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "LastUpdated",
			Code:    "invalid",
			Message: "LastUpdated timestamp is zero",
			Value:   response.LastUpdated,
		})
	}

	// URL 형식 검증
	if response.Website != "" {
		// URL 형식 검증 로직 추가 가능
	}

	if response.LogoURL != "" {
		// URL 형식 검증 로직 추가 가능
	}
}
