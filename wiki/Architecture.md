# 아키텍처

이 문서는 Go Finance Chart 프로젝트의 아키텍처와 주요 컴포넌트에 대해 설명합니다. 시스템의 핵심 구성 요소와 상호작용 방식을 이해하는 데 도움이 됩니다.

## 시스템 개요

Go Finance Chart는 클린 아키텍처 원칙을 따르는 계층화된 구조로 설계되었습니다. 각 계층은 특정 책임을 가지며, 의존성은 내부 계층(도메인)을 향하도록 설계되었습니다.

```
┌───────────────────────────────────────────────────────────────┐
│                       사용자 인터페이스                         │
└───────────────────────────────┬───────────────────────────────┘
                                │
┌───────────────────────────────┼───────────────────────────────┐
│                       애플리케이션 서비스                       │
└───────────────────────────────┼───────────────────────────────┘
                                │
┌───────────────────────────────┼───────────────────────────────┐
│                          도메인 계층                           │
└───────────────────────────────┼───────────────────────────────┘
                                │
┌───────────────────────────────┼───────────────────────────────┐
│                         인프라 계층                            │
└───────────────────────────────┴───────────────────────────────┘
```

## 주요 컴포넌트

### 1. 데이터 소스 인터페이스 (완료됨)

데이터 소스 인터페이스는 외부 데이터 제공자와의 일관된 통신을 위한 추상화 계층입니다.

#### 핵심 인터페이스

```go
type DataSource interface {
    // 데이터 가져오기 메서드
    Fetch(query DataQuery) (DataResponse, error)
    
    // 데이터 소스 상태 확인
    IsAvailable() bool
    
    // 데이터 소스 메타데이터
    GetMetadata() DataSourceMetadata
}

type MarketDataSource interface {
    DataSource
    
    // 가격 이력 조회
    GetPriceHistory(ticker string, period TimePeriod) ([]PricePoint, error)
    
    // 기업 정보 조회
    GetCompanyInfo(ticker string) (CompanyInfo, error)
    
    // 시장 요약 정보
    GetMarketSummary() (MarketSummary, error)
}

type NewsDataSource interface {
    DataSource
    
    // 뉴스 기사 검색
    SearchNews(query string, filters NewsFilters) ([]NewsArticle, error)
    
    // 특정 자산 관련 뉴스
    GetNewsForAsset(ticker string, limit int) ([]NewsArticle, error)
    
    // 최신 뉴스 가져오기
    GetLatestNews(category string, limit int) ([]NewsArticle, error)
}
```

#### 구현체

데이터 소스 인터페이스는 다음과 같은 구현체를 가질 예정입니다:

1. **YahooFinanceDataSource**: Yahoo Finance API로부터 시장 데이터를 가져오는 구현체 (#8)
2. **FinancialNewsDataSource**: 다양한 뉴스 API를 통합하는 구현체 (#9)
3. **MockDataSource**: 테스트 및 개발 목적의 가짜 데이터 제공 구현체

### 2. Repository 패턴 (진행 중)

Repository 패턴은 데이터 접근 로직을 추상화하고 일관된 인터페이스를 제공합니다.

#### 제네릭 Repository 인터페이스 (완료됨)

```go
type Repository[T Entity] interface {
    // 기본 CRUD 작업
    FindByID(id string) (T, error)
    FindAll() ([]T, error)
    Save(entity T) error
    Update(entity T) error
    Delete(id string) error
    
    // 페이징 및 필터링
    FindWithFilter(filter Filter) ([]T, error)
    Count(filter Filter) (int, error)
}
```

#### 도메인별 Repository

현재 세 가지 주요 도메인에 대한 Repository를 구현/리팩토링 중입니다:

1. **AssetRepository** (#14): 금융 자산 데이터 접근
   ```go
   type AssetRepository interface {
       Repository[Asset]
       
       // 도메인 특화 메서드
       FindByTicker(ticker string) (Asset, error)
       FindByCategory(category AssetCategory) ([]Asset, error)
       FindByMarket(market string) ([]Asset, error)
   }
   ```

2. **PortfolioRepository** (#15): 투자자 포트폴리오 데이터 접근
   ```go
   type PortfolioRepository interface {
       Repository[Portfolio]
       
       // 도메인 특화 메서드
       FindByUserID(userID string) ([]Portfolio, error)
       GetPortfolioPerformance(portfolioID string, period TimePeriod) (PerformanceData, error)
       AddAssetToPortfolio(portfolioID string, assetID string, quantity float64) error
   }
   ```

3. **TransactionRepository** (#16): 거래 이력 데이터 접근
   ```go
   type TransactionRepository interface {
       Repository[Transaction]
       
       // 도메인 특화 메서드
       FindByPortfolioID(portfolioID string) ([]Transaction, error)
       FindByAssetID(assetID string) ([]Transaction, error)
       FindByDateRange(startDate, endDate time.Time) ([]Transaction, error)
       GetTransactionStatsByType(portfolioID string) (TransactionStats, error)
   }
   ```

### 3. 분석 엔진 (계획 중)

분석 엔진은 시장 데이터와 포트폴리오 정보를 처리하여 인사이트를 제공합니다.

#### 핵심 컴포넌트

1. **포트폴리오 성과 분석기** (#10)
   - 수익률 계산
   - 리스크 지표 평가
   - 섹터별/자산별 분석

2. **백테스팅 엔진** (#5)
   - 과거 데이터 기반 투자 전략 검증
   - 성과 지표 계산
   - 시나리오 분석

3. **시장 동향 분석기** (#12)
   - 기술적 지표 계산
   - 시장 감성 분석
   - 이상치 감지

4. **포트폴리오 최적화 엔진** (#11)
   - 효율적 투자선 계산
   - 리스크-리턴 최적화
   - 자산 배분 제안

### 4. 사용자 인터페이스 (일부 완료)

사용자 인터페이스는 데이터와 분석 결과를 시각화하고 사용자 상호작용을 처리합니다.

#### 기본 UI 구성요소 (완료됨)

- 차트 컴포넌트
- 데이터 테이블
- 필터 및 검색 컴포넌트
- 폼 및 입력 컴포넌트
- 모달 및 팝업
- 반응형 레이아웃 시스템

#### 알림 시스템 (계획됨)

앱 내 알림 시스템(#6)은 다음과 같은 기능을 제공할 예정입니다:

- 실시간 알림 처리
- 알림 우선순위 관리
- 사용자 설정 가능한 알림 규칙
- 다양한 알림 채널 지원 (앱 내, 이메일, 푸시 등)

## 데이터 흐름 다이어그램

다음 다이어그램은 Go Finance Chart 내 주요 데이터 흐름을 보여줍니다:

```
외부 데이터 소스 (Yahoo Finance, 뉴스 API)
       │
       ▼
┌─────────────────┐
│  데이터 소스    │
│  인터페이스     │───┐
└─────────────────┘   │
       │              │
       ▼              │
┌─────────────────┐   │
│  데이터 변환    │   │
│  & 정제         │   │
└─────────────────┘   │
       │              │
       ▼              │
┌─────────────────┐   │
│  Repository     │◀──┘
│  인터페이스     │
└─────────────────┘
       │
       ▼
┌─────────────────┐
│  도메인 서비스  │
│  & 분석 엔진    │
└─────────────────┘
       │
       ▼
┌─────────────────┐
│  사용자         │
│  인터페이스     │
└─────────────────┘
```

## 기술 스택

### 백엔드
- Go 언어 (1.18+, 제네릭 지원)
- 데이터베이스: 추후 결정 예정 (#4)
- HTTP 클라이언트: 표준 라이브러리 또는 서드파티 클라이언트

### 프론트엔드
- 프레임워크/라이브러리: 추후 결정 예정
- 차트 라이브러리: 고성능 금융 차트 라이브러리

### 개발 도구
- 버전 관리: Git
- CI/CD: 추후 결정 예정
- 테스트: Go 표준 테스트 라이브러리

## 향후 아키텍처 고려사항

- **확장성**: 시스템이 증가하는 데이터 볼륨과 사용자 수를 처리할 수 있도록 설계
- **성능**: 대량의 금융 데이터 처리 및 분석을 위한 최적화
- **보안**: 민감한 금융 데이터를 보호하기 위한 보안 조치
- **유지보수성**: 모듈화된 설계로 미래 변경 및 기능 추가가 용이하도록 함
- **테스트 용이성**: 각 컴포넌트가 독립적으로 테스트 가능하도록 설계