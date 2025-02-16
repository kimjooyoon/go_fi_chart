# 유비쿼터스 언어 정의

## 공통 개념

### Aggregate
- 정의: 하나의 단위로 취급되는 도메인 객체들의 클러스터
- 특징:
- 트랜잭션 일관성 경계
- 하나의 루트 엔티티
- 독립적인 비즈니스 규칙

### Entity
- 정의: 고유한 식별자를 가진 도메인 객체
- 특징:
- 수명주기 존재
- 상태 변경 가능
- ID로 동일성 판단

### Value Object
- 정의: 속성으로만 정의되는 불변 객체
- 특징:
- 식별자 없음
- 불변성
- 속성으로 동등성 판단

### Domain Event
- 정의: 도메인에서 발생한 의미 있는 변화
- 특징:
- 과거 시제로 명명
- 발생 시점 포함
- 변경 불가능

## 자산 관리 도메인

### Asset Aggregate
```typescript
aggregate Asset {
// Root Entity
Asset {
id: AggregateId        // 자산 고유 식별자
name: string           // 자산 이름
type: AssetType        // 자산 유형
currentValue: Money    // 현재 가치
priceHistory: PricePoint[] // 가격 이력
metadata: Map
<string, string> // 추가 정보
}

// Value Objects
AssetType: 'stock' | 'bond' | 'cash' // 자산 유형
Money {
amount: Decimal    // 금액
currency: string   // 통화
}
PricePoint {
timestamp: DateTime // 시점
value: Money       // 가격
}
}
```

#### 불변식
- 자산 가치는 0 이상이어야 함
- 가격 이력은 시간순으로 정렬되어야 함
- 동일 시점에 중복된 가격 없음

### Portfolio Aggregate
```typescript
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
```

#### 불변식
- 자산 할당 비율의 합은 100%
- 각 할당 비율은 0% 초과
- 목표 가치와 현재 가치는 0 이상

### Transaction Aggregate
```typescript
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

#### 불변식
- 거래량은 0보다 커야 함
- 거래 가격은 0보다 커야 함
- 실행된 거래는 변경 불가

## 분석 도메인

### Analysis Aggregate
```typescript
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
}
```

#### 불변식
- 분석 기간의 시작은 종료보다 이전
- 변동성은 0 이상
- 상관계수는 -1에서 1 사이

### TimeSeries Aggregate
```typescript
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

#### 불변식
- 데이터 포인트는 시간순 정렬
- 신뢰도는 0에서 1 사이
- 동일 시점 중복 데이터 없음

## 모니터링 도메인

### Metric Aggregate
```typescript
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
```

#### 불변식
- 메트릭 이름은 유효한 식별자
- Counter는 감소할 수 없음
- 타임스탬프는 현재 이전

### Alert Aggregate
```typescript
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

#### 불변식
- 알림 메시지는 비어있지 않음
- 상태 변경은 정해진 흐름 준수
- Critical 알림은 즉시 처리

## 도메인 이벤트

### 자산 관리 이벤트
```typescript
// 자산 생성됨
event AssetCreated {
assetId: string
name: string
type: AssetType
initialValue: Money
timestamp: DateTime
}

// 포트폴리오 재조정됨
event PortfolioRebalanced {
portfolioId: string
newAllocations: AssetAllocation[]
timestamp: DateTime
}

// 거래 실행됨
event TransactionExecuted {
transactionId: string
assetId: string
type: TransactionType
quantity: Decimal
price: Money
timestamp: DateTime
}
```

### 분석 이벤트
```typescript
// 분석 완료됨
event AnalysisCompleted {
analysisId: string
portfolioId: string
metrics: AnalysisMetrics
timestamp: DateTime
}

// 리스크 경고 발생됨
event RiskAlertRaised {
portfolioId: string
riskType: string
level: AlertLevel
timestamp: DateTime
}
```

### 모니터링 이벤트
```typescript
// 메트릭 수집됨
event MetricCollected {
metricId: string
name: string
value: MetricValue
timestamp: DateTime
}

// 알림 발생됨
event AlertTriggered {
alertId: string
level: AlertLevel
message: string
timestamp: DateTime
}
```

## 공통 규칙

### 식별자
- UUID v4 사용
- 컨텍스트 내 유일성 보장
- 생성 후 변경 불가

### 시간 처리
- UTC 기준 저장
- ISO-8601 형식 사용
- 밀리초 단위까지 기록

### 금액 처리
- Decimal 타입 사용
- 소수점 8자리까지 허용
- 통화 단위 필수 기록

### 이벤트 처리
- 이벤트는 불변
- 순서 보장 필요
- 멱등성 고려 