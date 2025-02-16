# Portfolio 서비스

## 개요
Portfolio 서비스는 사용자의 투자 포트폴리오를 관리하는 마이크로서비스입니다. 이벤트 기반 아키텍처를 사용하여 포트폴리오 상태 변경을 추적하고, 결과적 일관성 모델을 통해 다른 서비스와 통합됩니다.

## 주요 기능
- 포트폴리오 생성, 조회, 수정, 삭제
- 자산 배분 관리
- 포트폴리오 성과 분석
- 이벤트 기반 상태 관리
- 실시간 리밸런싱 알림

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

### gRPC 서비스
```protobuf
service PortfolioService {
rpc ValidatePortfolio(ValidatePortfolioRequest) returns (ValidatePortfolioResponse);
rpc CalculatePerformance(CalculatePerformanceRequest) returns (CalculatePerformanceResponse);
rpc WatchPortfolioChanges(WatchPortfolioRequest) returns (stream PortfolioEvent);
}
```

## 도메인 이벤트
- PortfolioCreated
- PortfolioUpdated
- PortfolioDeleted
- AssetAllocated
- RebalancingRequested
- RebalancingCompleted

## 의존성
- Asset 서비스: 포트폴리오 구성 자산 정보 조회
- Transaction 서비스: 포트폴리오 거래 내역 연동
- Monitoring 서비스: 포트폴리오 상태 모니터링

## 설정
### 환경 변수
- `PORT`: 서비스 포트 (기본값: 8081)
- `MONGODB_URI`: MongoDB 연결 문자열
- `POSTGRES_URI`: PostgreSQL 연결 문자열
- `GRPC_PORT`: gRPC 서버 포트 (기본값: 9081)
- `ASSET_SERVICE_URL`: Asset 서비스 gRPC 주소

## 로컬 개발 환경 설정
```bash
# 서비스 실행
go run cmd/server/main.go

# 테스트 실행
go test ./...
```

## 도메인 모델

### Portfolio (포트폴리오)
```go
type Portfolio struct {
ID        string
UserID    string
Name      string
Assets    []PortfolioAsset
CreatedAt time.Time
UpdatedAt time.Time
}
```

### PortfolioAsset (값 객체)
```go
type PortfolioAsset struct {
AssetID string
Weight  Percentage
}
```

### Percentage (값 객체)
```go
type Percentage struct {
Value float64 // 0-100
}
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

## 리밸런싱 프로세스

### 트리거 조건
- 자산 가치 변동
- 정기 리밸런싱 일정
- 수동 요청

### 실행 단계
1. 리밸런싱 필요성 분석
2. 리밸런싱 계획 수립
3. 거래 주문 생성
4. 실행 상태 모니터링
5. 결과 보고 