# 도메인 인터페이스 계약

## 공통 사항

### 이벤트 형식
```json
{
"eventId": "uuid",
"eventType": "DomainEvent",
"aggregateId": "uuid",
"aggregateType": "string",
"version": "integer",
"timestamp": "ISO8601",
"data": {},
"metadata": {
"correlationId": "uuid",
"causationId": "uuid",
"userId": "uuid"
}
}
```

### 명령 형식
```json
{
"commandId": "uuid",
"commandType": "string",
"aggregateId": "uuid",
"data": {},
"metadata": {
"userId": "uuid",
"timestamp": "ISO8601"
}
}
```

### 쿼리 형식
```json
{
"queryId": "uuid",
"queryType": "string",
"parameters": {},
"metadata": {
"userId": "uuid",
"timestamp": "ISO8601"
}
}
```

## 1. 자산 관리 도메인

### 명령 (Commands)
```typescript
// 자산 생성
interface CreateAssetCommand {
commandType: "CreateAsset"
data: {
name: string
type: "stock" | "bond" | "cash"
initialValue: Money
metadata?: Record
<string, string>
}
}

// 포트폴리오 생성
interface CreatePortfolioCommand {
commandType: "CreatePortfolio"
data: {
name: string
allocations: Array<{
assetId: string
ratio: number
}>
}
}

// 거래 실행
interface ExecuteTransactionCommand {
commandType: "ExecuteTransaction"
data: {
assetId: string
type: "buy" | "sell"
quantity: number
price: Money
}
}
```

### 이벤트 (Events)
```typescript
// 자산 생성됨
interface AssetCreatedEvent {
eventType: "AssetCreated"
data: {
assetId: string
name: string
type: "stock" | "bond" | "cash"
initialValue: Money
}
}

// 포트폴리오 재조정됨
interface PortfolioRebalancedEvent {
eventType: "PortfolioRebalanced"
data: {
portfolioId: string
newAllocations: Array<{
assetId: string
ratio: number
targetValue: Money
}>
}
}
```

### 쿼리 (Queries)
```typescript
// 자산 조회
interface GetAssetQuery {
queryType: "GetAsset"
parameters: {
assetId: string
}
}

// 포트폴리오 성과 조회
interface GetPortfolioPerformanceQuery {
queryType: "GetPortfolioPerformance"
parameters: {
portfolioId: string
timeRange: {
start: string // ISO8601
end: string   // ISO8601
}
}
}
```

## 2. 분석 도메인

### 명령 (Commands)
```typescript
// 포트폴리오 분석 시작
interface StartPortfolioAnalysisCommand {
commandType: "StartPortfolioAnalysis"
data: {
portfolioId: string
analysisType: "risk" | "performance" | "optimization"
parameters: Record
<string, any>
}
}
```

### 이벤트 (Events)
```typescript
// 분석 완료됨
interface AnalysisCompletedEvent {
eventType: "AnalysisCompleted"
data: {
analysisId: string
portfolioId: string
results: {
metrics: Record
<string, number>
recommendations: string[]
}
}
}
```

### 쿼리 (Queries)
```typescript
// 분석 결과 조회
interface GetAnalysisResultQuery {
queryType: "GetAnalysisResult"
parameters: {
analysisId: string
}
}
```

## 3. 모니터링 도메인

### 명령 (Commands)
```typescript
// 메트릭 기록
interface RecordMetricCommand {
commandType: "RecordMetric"
data: {
name: string
value: number
labels: Record
<string, string>
}
}

// 알림 생성
interface CreateAlertCommand {
commandType: "CreateAlert"
data: {
level: "INFO" | "WARNING" | "ERROR" | "CRITICAL"
source: string
message: string
}
}
```

### 이벤트 (Events)
```typescript
// 메트릭 수집됨
interface MetricCollectedEvent {
eventType: "MetricCollected"
data: {
metricId: string
name: string
value: number
timestamp: string
}
}

// 알림 발생됨
interface AlertTriggeredEvent {
eventType: "AlertTriggered"
data: {
alertId: string
level: string
message: string
}
}
```

## 메시지 브로커 설정

### Kafka 토픽 구조
```
go-fi-chart.{domain}.{aggregate}.{event}
예: go-fi-chart.asset.portfolio.rebalanced
```

### 이벤트 스토어 스트림
```
{domain}-{aggregate}
예: asset-portfolio
```

## 저장소 인터페이스

### 애그리게잇 저장소
```typescript
interface AggregateRepository
<T> {
    save(aggregate: T): Promise
    <void>
        load(id: string): Promise
        <T>
            exists(id: string): Promise
            <boolean>
                }
                ```

                ### 이벤트 저장소
                ```typescript
                interface EventStore {
                appendToStream(streamId: string, events: DomainEvent[]): Promise
                <void>
                    readFromStream(streamId: string): Promise
                    <DomainEvent
                    []>
                    subscribeToStream(streamId: string, handler: EventHandler): Promise
                    <void>
                        }
                        ```

                        ### 읽기 모델 저장소
                        ```typescript
                        interface ReadModelRepository
                        <T> {
                            save(model: T): Promise
                            <void>
                                find(query: Query): Promise
                                <T
                                []>
                                findOne(query: Query): Promise
                                <T>
                                    }
                                    ```