# 바운디드 컨텍스트

## 개요

Go Fi Chart는 다음과 같은 바운디드 컨텍스트로 구성됩니다:

## 1. 자산 관리 컨텍스트 (Asset Management)

### 책임
- 자산 생명주기 관리
- 포트폴리오 구성 관리
- 거래 처리 및 기록

### Aggregates
```typescript
// Asset Aggregate
aggregate Asset {
// Root Entity
Asset {
id: AggregateId
name: string
type: AssetType
currentValue: Money
priceHistory: PricePoint[]
metadata: Map
<string, string>
}

// Value Objects
AssetType: 'stock' | 'bond' | 'cash'
Money {
amount: Decimal
currency: string
}
PricePoint {
timestamp: DateTime
value: Money
}
}

// Portfolio Aggregate
aggregate Portfolio {
// Root Entity
Portfolio {
id: AggregateId
name: string
allocations: AssetAllocation[]
totalValue: Money
metadata: Map
<string, string>
}

// Entities
AssetAllocation {
assetId: string
ratio: Percentage
targetValue: Money
currentValue: Money
}
}

// Transaction Aggregate
aggregate Transaction {
// Root Entity
Transaction {
id: AggregateId
type: TransactionType
assetId: string
quantity: Decimal
price: Money
executedAt: DateTime
status: TransactionStatus
}
}
```

### 도메인 서비스
```go
type PortfolioService interface {
Rebalance(ctx context.Context, portfolio Portfolio) error
CalculateValue(ctx context.Context, portfolio Portfolio) (Money, error)
}

type TransactionService interface {
ExecuteTransaction(ctx context.Context, transaction Transaction) error
ValidateTransaction(ctx context.Context, transaction Transaction) error
}
```

## 2. 분석 컨텍스트 (Analysis)

### 책임
- 시계열 데이터 분석
- 포트폴리오 성과 분석
- 리스크 평가

### Aggregates
```typescript
// Analysis Aggregate
aggregate Analysis {
// Root Entity
PortfolioAnalysis {
id: AggregateId
portfolioId: string
period: AnalysisPeriod
metrics: AnalysisMetrics
riskMetrics: RiskMetrics
status: AnalysisStatus
}
}

// TimeSeries Aggregate
aggregate TimeSeries {
// Root Entity
TimeSeriesData {
id: AggregateId
assetId: string
dataPoints: DataPoint[]
metadata: Map
<string, string>
}
}
```

### 도메인 서비스
```go
type AnalysisService interface {
AnalyzePortfolio(ctx context.Context, portfolioId string, period AnalysisPeriod) error
CalculateRiskMetrics(ctx context.Context, analysisId string) error
}

type TimeSeriesService interface {
ProcessTimeSeriesData(ctx context.Context, data TimeSeriesData) error
AnalyzeTimeSeries(ctx context.Context, timeSeriesId string) error
}
```

## 3. 모니터링 컨텍스트 (Monitoring)

### 책임
- 시스템 메트릭 수집
- 알림 관리
- 상태 모니터링

### Aggregates
```typescript
// Metric Aggregate
aggregate Metric {
// Root Entity
Metric {
id: AggregateId
name: string
type: MetricType
value: MetricValue
labels: Map
<string, string>
}
}

// Alert Aggregate
aggregate Alert {
// Root Entity
Alert {
id: AggregateId
level: AlertLevel
source: string
message: string
status: AlertStatus
metadata: Map
<string, string>
}
}
```

### 도메인 서비스
```go
type MetricService interface {
CollectMetrics(ctx context.Context) error
ProcessMetrics(ctx context.Context, metrics []Metric) error
}

type AlertService interface {
ProcessAlert(ctx context.Context, alert Alert) error
NotifyAlert(ctx context.Context, alert Alert) error
}
```

## 컨텍스트 간 통신

### 1. 이벤트 기반 통신
```go
// Domain Events
type DomainEvent interface {
AggregateID() string
EventType() string
Version() int
Timestamp() time.Time
Payload() interface{}
}

// Event Types
const (
AssetCreated = "asset.created"
AssetUpdated = "asset.updated"
PortfolioRebalanced = "portfolio.rebalanced"
TransactionExecuted = "transaction.executed"
AnalysisCompleted = "analysis.completed"
AlertTriggered = "alert.triggered"
)

// Event Store
type EventStore interface {
SaveEvents(ctx context.Context, aggregateId string, events []DomainEvent) error
GetEvents(ctx context.Context, aggregateId string) ([]DomainEvent, error)
}
```

### 2. 동기 통신 (필요한 경우)
```go
// Context Interfaces
type AssetManagementContext interface {
GetAsset(ctx context.Context, id string) (Asset, error)
GetPortfolio(ctx context.Context, id string) (Portfolio, error)
}

type AnalysisContext interface {
GetAnalysis(ctx context.Context, id string) (Analysis, error)
GetTimeSeries(ctx context.Context, id string) (TimeSeries, error)
}
```

## 데이터 일관성 전략

### 1. Aggregate 내부
- 강한 일관성 (Strong Consistency)
- 트랜잭션 단위로 처리
- 불변식 즉시 적용

### 2. Aggregate 간
- 최종 일관성 (Eventual Consistency)
- 이벤트 기반 동기화
- SAGA 패턴 활용

### 3. 컨텍스트 간
- 느슨한 결합
- 비동기 통신 선호
- Anti-Corruption Layer 사용

## 구현 가이드라인

### 1. Aggregate 설계
```go
// Aggregate Root Interface
type AggregateRoot interface {
ID() string
Version() int
Events() []DomainEvent
ClearEvents()
}

// Base Aggregate Implementation
type BaseAggregate struct {
id      string
version int
events  []DomainEvent
}

func (a *BaseAggregate) ID() string { return a.id }
func (a *BaseAggregate) Version() int { return a.version }
func (a *BaseAggregate) Events() []DomainEvent { return a.events }
func (a *BaseAggregate) ClearEvents() { a.events = nil }
```

### 2. Repository 패턴
```go
// Generic Repository Interface
type Repository[T AggregateRoot] interface {
Save(ctx context.Context, aggregate T) error
FindById(ctx context.Context, id string) (T, error)
Delete(ctx context.Context, id string) error
}

// Event Sourcing Repository
type EventSourcingRepository[T AggregateRoot] interface {
Repository[T]
SaveEvents(ctx context.Context, events []DomainEvent) error
GetEvents(ctx context.Context, aggregateId string) ([]DomainEvent, error)
}
```

### 3. Anti-Corruption Layer
```go
// Example for External Service Integration
type ExternalPricingService interface {
GetPrice(symbol string) (float64, error)
}

// Anti-Corruption Layer
type PricingAdapter struct {
external ExternalPricingService
}

func (a *PricingAdapter) GetAssetPrice(asset Asset) (Money, error) {
price, err := a.external.GetPrice(asset.Symbol())
if err != nil {
return Money{}, err
}
return NewMoney(price, asset.Currency()), nil
}
```

## 테스트 전략

### 1. Aggregate 테스트
```go
func TestPortfolioAggregate(t *testing.T) {
// Given
portfolio := NewPortfolio("Test Portfolio")

// When
portfolio.AddAllocation(AssetID("1"), Percentage(50))
portfolio.AddAllocation(AssetID("2"), Percentage(50))

// Then
assert.Len(t, portfolio.Events(), 2)
assert.Equal(t, "portfolio.allocation.added", portfolio.Events()[0].EventType())
}
```

### 2. 도메인 서비스 테스트
```go
func TestPortfolioService(t *testing.T) {
// Given
service := NewPortfolioService(mockRepo, mockEventStore)
portfolio := NewPortfolio("Test Portfolio")

// When
err := service.Rebalance(context.Background(), portfolio)

// Then
assert.NoError(t, err)
assert.Equal(t, "rebalanced", portfolio.Status())
}
```

### 3. 통합 테스트
```go
func TestAssetManagementContext(t *testing.T) {
// Given
ctx := NewAssetManagementContext(config)

// When
asset := ctx.CreateAsset(NewAssetCommand{...})
portfolio := ctx.CreatePortfolio(NewPortfolioCommand{...})

// Then
assert.NotNil(t, asset)
assert.NotNil(t, portfolio)
assert.NoError(t, ctx.Errors())
}
``` 