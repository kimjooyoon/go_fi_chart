# API 설계

## 개요

시스템은 외부 통신을 위한 GraphQL API와 서비스 간 통신을 위한 gRPC를 사용합니다. 모든 상태 변경은 이벤트로 발행되며, 결과적 일관성 모델을 따릅니다.

## API 계층

### 1. 외부 API (GraphQL)
- 클라이언트 애플리케이션 통신
- 쿼리와 뮤테이션 분리
- 구독을 통한 실시간 업데이트
- 스키마 예시:
```graphql
type Asset {
id: ID!
userId: ID!
type: AssetType!
name: String!
amount: Money!
performance: Performance
goals: [Goal!]
createdAt: DateTime!
updatedAt: DateTime!
}

type Money {
amount: Float!
currency: String!
}

input CreateAssetInput {
userId: ID!
type: AssetType!
name: String!
amount: Float!
currency: String!
}

type Mutation {
createAsset(input: CreateAssetInput!): Asset!
updateAsset(id: ID!, input: UpdateAssetInput!): Asset!
deleteAsset(id: ID!): Boolean!
}

type Query {
asset(id: ID!): Asset
assets(userId: ID!): [Asset!]!
assetsByType(type: AssetType!): [Asset!]!
}

type Subscription {
assetUpdated(userId: ID!): Asset!
goalAchieved(assetId: ID!): Goal!
}
```

### 2. 서비스 간 통신 (gRPC)
- 내부 서비스 간 통신
- 프로토콜 버퍼 정의
- 양방향 스트리밍 지원
- 프로토 정의 예시:
```protobuf
syntax = "proto3";

package asset;

service AssetService {
rpc ValidateAsset(ValidateAssetRequest) returns (ValidateAssetResponse);
rpc GetAssetValue(GetAssetValueRequest) returns (GetAssetValueResponse);
rpc WatchAssetChanges(WatchAssetRequest) returns (stream AssetEvent);
}

message Asset {
string id = 1;
string user_id = 2;
AssetType type = 3;
string name = 4;
Money amount = 5;
// ...
}

message Money {
double amount = 1;
string currency = 2;
}
```

### 3. 이벤트 스트림
- 도메인 이벤트 발행/구독
- MongoDB 이벤트 저장소
- 이벤트 스키마 예시:
```json
{
"eventId": "uuid",
"eventType": "AssetCreated",
"version": "1.0",
"timestamp": "2024-02-16T12:00:00Z",
"aggregateId": "asset-123",
"data": {
"userId": "user-123",
"assetType": "STOCK",
"name": "AAPL",
"amount": {
"value": 1000.00,
"currency": "USD"
}
}
}
```

## 공통 설계 원칙

### 1. 버전 관리
- GraphQL 스키마 버전 관리
- 프로토콜 버퍼 버전 관리
- 이벤트 스키마 버전 관리

### 2. 인증/인가
- JWT 기반 인증
- GraphQL 지시어를 통한 권한 제어
- gRPC 인터셉터를 통한 서비스 인증

### 3. 에러 처리
```graphql
type Error {
code: String!
message: String!
path: [String!]
extensions: JSONObject
}
```

### 4. 페이지네이션
- 커서 기반 페이지네이션
```graphql
type PageInfo {
hasNextPage: Boolean!
endCursor: String
}

type AssetConnection {
edges: [AssetEdge!]!
pageInfo: PageInfo!
}

type AssetEdge {
node: Asset!
cursor: String!
}
```

### 5. 필터링 및 검색
```graphql
input AssetFilter {
types: [AssetType!]
minAmount: Float
maxAmount: Float
dateRange: DateRange
}

type Query {
searchAssets(filter: AssetFilter!): AssetConnection!
}
```

## 도메인별 API

### Asset 서비스
1. GraphQL API
```graphql
type Query {
asset(id: ID!): Asset
assets(userId: ID!): [Asset!]!
assetPerformance(id: ID!): Performance!
}

type Mutation {
createAsset(input: CreateAssetInput!): Asset!
updateAsset(id: ID!, input: UpdateAssetInput!): Asset!
deleteAsset(id: ID!): Boolean!
assignGoal(assetId: ID!, input: GoalInput!): Goal!
}

type Subscription {
assetUpdated(userId: ID!): Asset!
goalProgress(goalId: ID!): GoalProgress!
}
```

2. gRPC 서비스
```protobuf
service AssetService {
rpc ValidateAsset(ValidateAssetRequest) returns (ValidateAssetResponse);
rpc CalculateValue(CalculateValueRequest) returns (CalculateValueResponse);
rpc WatchAssetChanges(WatchAssetRequest) returns (stream AssetEvent);
}
```

### Portfolio 서비스
1. GraphQL API
```graphql
type Query {
portfolio(id: ID!): Portfolio
portfolios(userId: ID!): [Portfolio!]!
portfolioPerformance(id: ID!): Performance!
}

type Mutation {
createPortfolio(input: CreatePortfolioInput!): Portfolio!
updatePortfolio(id: ID!, input: UpdatePortfolioInput!): Portfolio!
deletePortfolio(id: ID!): Boolean!
rebalancePortfolio(id: ID!): RebalanceResult!
}

type Subscription {
portfolioUpdated(userId: ID!): Portfolio!
rebalanceProgress(portfolioId: ID!): RebalanceProgress!
}
```

### Transaction 서비스
1. GraphQL API
```graphql
type Query {
transaction(id: ID!): Transaction
transactions(filter: TransactionFilter!): TransactionConnection!
transactionSummary(assetId: ID!): TransactionSummary!
}

type Mutation {
createTransaction(input: CreateTransactionInput!): Transaction!
cancelTransaction(id: ID!): Transaction!
}

type Subscription {
transactionStatus(id: ID!): TransactionStatus!
}
```

## API 문서화

### 1. GraphQL 스키마
- 스키마 정의
- 타입 설명
- 예제 쿼리

### 2. gRPC 프로토콜
- 프로토콜 버퍼 정의
- 서비스 설명
- 스트리밍 예제

### 3. 이벤트 스키마
- 이벤트 타입 정의
- 페이로드 스키마
- 버전 관리 정책 