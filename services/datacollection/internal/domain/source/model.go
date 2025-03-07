package source

import (
	"time"
)

// AssetType은 금융 자산 유형을 정의합니다.
type AssetType string

const (
	AssetTypeStock      AssetType = "stock"       // 주식
	AssetTypeBond       AssetType = "bond"        // 채권
	AssetTypeETF        AssetType = "etf"         // ETF
	AssetTypeCrypto     AssetType = "crypto"      // 암호화폐
	AssetTypeCommodity  AssetType = "commodity"   // 원자재
	AssetTypeFuture     AssetType = "future"      // 선물
	AssetTypeForex      AssetType = "forex"       // 외환
	AssetTypeIndex      AssetType = "index"       // 지수
	AssetTypeMutualFund AssetType = "mutual_fund" // 뮤추얼 펀드
)

// Interval은 데이터 수집 간격을 정의합니다.
type Interval string

const (
	Interval1Min      Interval = "1m"  // 1분
	Interval5Min      Interval = "5m"  // 5분
	Interval15Min     Interval = "15m" // 15분
	Interval30Min     Interval = "30m" // 30분
	Interval1Hour     Interval = "1h"  // 1시간
	Interval4Hour     Interval = "4h"  // 4시간
	IntervalDaily     Interval = "1d"  // 일별
	IntervalWeekly    Interval = "1wk" // 주별
	IntervalMonthly   Interval = "1mo" // 월별
	IntervalQuarterly Interval = "3mo" // 분기별
	IntervalYearly    Interval = "1y"  // 연별
)

// PriceData는 시간에 따른 가격 데이터를 나타냅니다.
type PriceData struct {
	Timestamp     time.Time // 타임스탬프
	Open          float64   // 시가
	High          float64   // 고가
	Low           float64   // 저가
	Close         float64   // 종가
	Volume        int64     // 거래량
	AdjustedClose float64   // 수정 종가
}

// HistoricalDataRequest는 과거 데이터 요청을 위한 구조체입니다.
type HistoricalDataRequest struct {
	Symbol    string    // 자산 심볼
	AssetType AssetType // 자산 유형
	Interval  Interval  // 데이터 간격
	StartTime time.Time // 시작 시간
	EndTime   time.Time // 종료 시간
}

// HistoricalDataResponse는 과거 데이터 응답을 위한 구조체입니다.
type HistoricalDataResponse struct {
	Symbol    string      // 자산 심볼
	AssetType AssetType   // 자산 유형
	Interval  Interval    // 데이터 간격
	Data      []PriceData // 가격 데이터 배열
}

// RealTimeDataRequest는 실시간 데이터 요청을 위한 구조체입니다.
type RealTimeDataRequest struct {
	Symbol    string    // 자산 심볼
	AssetType AssetType // 자산 유형
}

// RealTimeDataResponse는 실시간 데이터 응답을 위한 구조체입니다.
type RealTimeDataResponse struct {
	Symbol        string    // 자산 심볼
	AssetType     AssetType // 자산 유형
	CurrentPrice  float64   // 현재 가격
	Timestamp     time.Time // 가격 업데이트 시간
	Change        float64   // 변화량
	ChangePercent float64   // 변화율(%)
	Volume        int64     // 거래량
	MarketCap     float64   // 시가총액
	High24h       float64   // 24시간 최고가
	Low24h        float64   // 24시간 최저가
}

// MetadataRequest는 메타데이터 요청을 위한 구조체입니다.
type MetadataRequest struct {
	Symbol    string    // 자산 심볼
	AssetType AssetType // 자산 유형
}

// MetadataResponse는 자산 메타데이터 응답을 위한 구조체입니다.
type MetadataResponse struct {
	Symbol      string    // 자산 심볼
	AssetType   AssetType // 자산 유형
	Name        string    // 자산 이름
	Exchange    string    // 거래소
	Currency    string    // 통화
	Country     string    // 국가
	Description string    // 설명
	Sector      string    // 섹터
	Industry    string    // 산업
	Website     string    // 웹사이트
	LogoURL     string    // 로고 URL
	LastUpdated time.Time // 마지막 업데이트 시간
}
