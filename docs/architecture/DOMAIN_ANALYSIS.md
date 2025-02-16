# 도메인 분석

## 도메인 분류 및 Aggregate 설계

### 핵심 도메인 (Core Domain)

1. **자산 관리 (Asset Management)**
- 비즈니스 차별화의 핵심
- 높은 복잡도와 전문성 요구

#### Aggregates
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

// Value Objects
Percentage: Decimal // 0-100
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

// Value Objects
TransactionType: 'buy' | 'sell'
TransactionStatus: 'pending' | 'executed' | 'failed'
}
```

#### 도메인 서비스
- PortfolioValuationService
- AssetAllocationService
- TransactionExecutionService

#### 정책
- 포트폴리오 총 할당 비율은 100%를 초과할 수 없음
- 거래는 생성 후 수정 불가
- 자산 가격 변경 시 연관 포트폴리오 재평가

2. **분석 엔진 (Analysis Engine)**
- 핵심 비즈니스 가치 제공
- 복잡한 알고리즘 구현

#### Aggregates
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

// Value Objects
AnalysisPeriod {
start: DateTime
end: DateTime
}

AnalysisMetrics {
returns: Decimal
volatility: Decimal
sharpeRatio: Decimal
}

RiskMetrics {
var: Decimal
beta: Decimal
correlations: Map
<string, Decimal>
}

AnalysisStatus: 'in_progress' | 'completed' | 'failed'
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

// Value Objects
DataPoint {
timestamp: DateTime
value: Decimal
confidence: Decimal
}
}
```

#### 도메인 서비스
- TimeSeriesAnalysisService
- RiskAssessmentService
- PortfolioOptimizationService

### 일반 도메인 (Generic Domain)

1. **데이터 수집 (Data Collection)**
- 표준화된 프로세스 존재
- 재사용 가능한 컴포넌트

#### Aggregates
```typescript
// DataSource Aggregate
aggregate DataSource {
// Root Entity
DataSource {
id: AggregateId
type: SourceType
config: SourceConfig
status: SourceStatus
metadata: Map
<string, string>
}

// Value Objects
SourceType: 'api' | 'file' | 'database'
SourceStatus: 'active' | 'inactive' | 'error'
SourceConfig {
credentials: Credentials
endpoint: string
parameters: Map
<string, string>
}
}

// Pipeline Aggregate
aggregate Pipeline {
// Root Entity
Pipeline {
id: AggregateId
sourceId: string
steps: PipelineStep[]
schedule: Schedule
status: PipelineStatus
}

// Value Objects
PipelineStep {
type: StepType
config: Map
<string, string>
}
StepType: 'extract' | 'transform' | 'load'
PipelineStatus: 'running' | 'paused' | 'failed'
}
```

2. **시각화 (Visualization)**
- 일반적인 차트/대시보드 기능
- 시장 표준 존재

#### Aggregates
```typescript
// Chart Aggregate
aggregate Chart {
// Root Entity
Chart {
id: AggregateId
type: ChartType
dataSource: DataSourceConfig
config: ChartConfig
metadata: Map
<string, string>
}

// Value Objects
ChartType: 'line' | 'bar' | 'candlestick'
ChartConfig {
title: string
axes: AxisConfig[]
style: StyleConfig
}
}

// Dashboard Aggregate
aggregate Dashboard {
// Root Entity
Dashboard {
id: AggregateId
name: string
layout: LayoutConfig
widgets: Widget[]
}

// Entities
Widget {
id: string
type: WidgetType
sourceId: string
position: Position
config: Map
<string, string>
}

// Value Objects
WidgetType: 'chart' | 'metric' | 'alert'
Position {
x: number
y: number
width: number
height: number
}
}
```

### 지원 도메인 (Supporting Domain)

1. **모니터링 (Monitoring)**
- 시스템 운영 지원
- 표준화된 솔루션 존재

#### Aggregates
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

// Value Objects
MetricType: 'counter' | 'gauge' | 'histogram' | 'summary'
MetricValue {
raw: Decimal
timestamp: DateTime
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

// Value Objects
AlertLevel: 'info' | 'warning' | 'error' | 'critical'
AlertStatus: 'new' | 'acknowledged' | 'resolved'
}
```

2. **게이미피케이션 (Gamification)**
- 사용자 참여 촉진
- 부가 기능

#### Aggregates
```typescript
// Profile Aggregate
aggregate Profile {
// Root Entity
Profile {
id: AggregateId
userId: string
level: number
experience: number
badges: Badge[]
achievements: Achievement[]
}

// Value Objects
Badge {
type: BadgeType
earnedAt: DateTime
}
Achievement {
type: AchievementType
progress: number
completedAt: DateTime
}
}

// Reward Aggregate
aggregate Reward {
// Root Entity
Reward {
id: AggregateId
profileId: string
type: RewardType
amount: number
reason: string
status: RewardStatus
}

// Value Objects
RewardType: 'xp' | 'badge' | 'achievement'
RewardStatus: 'pending' | 'granted' | 'expired'
}
```

## 구현 전략

### 1. 영속성 전략
- NoSQL 데이터베이스 사용 (MongoDB)
- 각 Aggregate를 독립된 문서로 저장
- 참조는 ID만 유지
- 이벤트 소싱으로 변경 이력 관리

### 2. Repository 인터페이스
```go
type Repository[T any] interface {
Save(ctx context.Context, aggregate T) error
FindById(ctx context.Context, id string) (T, error)
Delete(ctx context.Context, id string) error
}

type EventStore interface {
SaveEvents(ctx context.Context, aggregateId string, events []DomainEvent) error
GetEvents(ctx context.Context, aggregateId string) ([]DomainEvent, error)
}
```

### 3. 도메인 이벤트
```go
type DomainEvent interface {
AggregateID() string
EventType() string
Version() int
Timestamp() time.Time
Payload() interface{}
}
```

### 4. CQRS 패턴
- 명령과 조회 책임 분리
- 읽기 모델은 목적에 맞게 최적화
- 이벤트 소싱과 연계하여 구현

## 구현 우선순위

1. 1단계 (즉시 구현)
- Asset Aggregate 구현
- Portfolio Aggregate 구현
- 기본 Repository 구현
- 이벤트 스토어 구현

2. 2단계 (2-3개월 내)
- Transaction Aggregate 구현
- Analysis Aggregate 구현
- CQRS 기반 조회 최적화
- 이벤트 핸들러 구현

3. 3단계 (4-6개월 내)
- 나머지 Aggregate 구현
- 도메인 서비스 구현
- 통합 테스트 구현

## 기술 스택

1. **백엔드**
- 언어: Go
- 데이터베이스: MongoDB
- 이벤트 스토어: EventStoreDB
- 메시지 큐: Apache Kafka

2. **프론트엔드**
- React + TypeScript
- 상태 관리: Redux
- 차트: D3.js

3. **인프라**
- Kubernetes
- Docker
- Prometheus + Grafana 