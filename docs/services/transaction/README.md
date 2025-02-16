# Transaction 서비스

## 개요
Transaction 서비스는 자산 거래 내역을 관리하는 마이크로서비스입니다. 이벤트 기반 아키텍처를 사용하여 거래 상태 변경을 추적하고, 결과적 일관성 모델을 통해 다른 서비스와 통합됩니다.

## 주요 기능
- 거래 내역 생성, 조회, 수정, 삭제
- 자산별 거래 내역 관리
- 포트폴리오별 거래 내역 관리
- 이벤트 기반 거래 처리
- 실시간 거래 상태 추적

## 기술 스택
- Go 1.24.0
- GraphQL API (외부 통신)
- gRPC (서비스 간 통신)
- MongoDB (이벤트 저장소)
- PostgreSQL (읽기 모델)

## API 엔드포인트

### GraphQL API
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
transactionExecuted(assetId: ID!): Transaction!
}
```

### gRPC 서비스
```protobuf
service TransactionService {
rpc ValidateTransaction(ValidateTransactionRequest) returns (ValidateTransactionResponse);
rpc ExecuteTransaction(ExecuteTransactionRequest) returns (ExecuteTransactionResponse);
rpc WatchTransactionStatus(WatchTransactionRequest) returns (stream TransactionEvent);
}
```

## 도메인 이벤트
- TransactionCreated
- TransactionExecuted
- TransactionFailed
- TransactionCancelled
- TransactionValidated
- TransactionSettled

## 의존성
- Asset 서비스: 거래 대상 자산 정보 조회
- Portfolio 서비스: 포트폴리오 정보 연동
- Monitoring 서비스: 거래 활동 모니터링

## 설정
### 환경 변수
- `PORT`: 서비스 포트 (기본값: 8082)
- `MONGODB_URI`: MongoDB 연결 문자열
- `POSTGRES_URI`: PostgreSQL 연결 문자열
- `GRPC_PORT`: gRPC 서버 포트 (기본값: 9082)
- `ASSET_SERVICE_URL`: Asset 서비스 gRPC 주소
- `PORTFOLIO_SERVICE_URL`: Portfolio 서비스 gRPC 주소

## 로컬 개발 환경 설정
```bash
# 서비스 실행
go run cmd/server/main.go

# 테스트 실행
go test ./...
```

## 도메인 모델

### Transaction (거래)
```go
type Transaction struct {
ID            uuid.UUID
UserID        uuid.UUID
PortfolioID   uuid.UUID
AssetID       uuid.UUID
Type          TransactionType
Amount        Money
Quantity      float64
ExecutedPrice Money
ExecutedAt    time.Time
Status        TransactionStatus
CreatedAt     time.Time
}
```

### TransactionType (값 객체)
```go
type TransactionType string

const (
Buy  TransactionType = "BUY"
Sell TransactionType = "SELL"
)
```

### TransactionStatus (값 객체)
```go
type TransactionStatus string

const (
Pending   TransactionStatus = "PENDING"
Executed  TransactionStatus = "EXECUTED"
Failed    TransactionStatus = "FAILED"
Cancelled TransactionStatus = "CANCELLED"
)
```

## CQRS 구현

### Command 모델
- MongoDB 이벤트 저장소
- 이벤트 소싱 패턴
- 트랜잭션 일관성

### Query 모델
- PostgreSQL 읽기 전용 뷰
- 이벤트 구독 기반 갱신
- 성능 최적화된 쿼리

## 거래 처리 프로세스

### 거래 생성
1. 거래 요청 검증
2. 자산 상태 확인
3. 거래 생성 이벤트 발행
4. 거래 실행 시작

### 거래 실행
1. 거래 유효성 검증
2. 자산 상태 업데이트
3. 거래 실행 이벤트 발행
4. 결과 통보

### 거래 취소
1. 취소 가능 여부 확인
2. 취소 이벤트 발행
3. 자산 상태 롤백
4. 결과 통보 