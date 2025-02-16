# Asset 서비스

## 개요
Asset 서비스는 사용자의 자산을 관리하는 마이크로서비스입니다. 이벤트 기반 아키텍처를 사용하여 자산의 상태 변경을 추적하고, 결과적 일관성 모델을 통해 다른 서비스와 통합됩니다.

## 주요 기능
- 자산 생성, 조회, 수정, 삭제
- 자산 가치 평가
- 자산 성과 추적
- 이벤트 기반 상태 관리

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
asset(id: ID!): Asset
assets(userId: ID!): [Asset!]!
assetsByType(type: AssetType!): [Asset!]!
}

type Mutation {
createAsset(input: CreateAssetInput!): Asset!
updateAsset(id: ID!, input: UpdateAssetInput!): Asset!
deleteAsset(id: ID!): Boolean!
}

type Subscription {
assetUpdated(userId: ID!): Asset!
goalAchieved(assetId: ID!): Goal!
}
```

### gRPC 서비스
```protobuf
service AssetService {
rpc ValidateAsset(ValidateAssetRequest) returns (ValidateAssetResponse);
rpc GetAssetValue(GetAssetValueRequest) returns (GetAssetValueResponse);
rpc WatchAssetChanges(WatchAssetRequest) returns (stream AssetEvent);
}
```

## 도메인 이벤트
- AssetCreated
- AssetUpdated
- AssetDeleted
- AssetValueChanged
- GoalAssigned
- GoalAchieved

## 의존성
- Portfolio 서비스: 포트폴리오 구성을 위한 자산 정보 제공
- Transaction 서비스: 자산 거래 정보 연동
- Monitoring 서비스: 자산 상태 모니터링

## 설정
### 환경 변수
- `PORT`: 서비스 포트 (기본값: 8080)
- `MONGODB_URI`: MongoDB 연결 문자열
- `POSTGRES_URI`: PostgreSQL 연결 문자열
- `GRPC_PORT`: gRPC 서버 포트 (기본값: 9080)

## 로컬 개발 환경 설정
```bash
# 서비스 실행
go run cmd/server/main.go

# 테스트 실행
go test ./...
```

## 도메인 모델

### Asset (자산)
```go
type Asset struct {
ID           string
UserID       string
Type         Type
Name         string
Amount       Money
Performance  *Performance
Goals        []*Goal
CreatedAt    time.Time
UpdatedAt    time.Time
}
```

### Money (값 객체)
```go
type Money struct {
Amount   float64
Currency string
}
```

### Performance (값 객체)
```go
type Performance struct {
StartValue     Money
CurrentValue   Money
GrowthRate     float64
RiskScore      float64
LastUpdateTime time.Time
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