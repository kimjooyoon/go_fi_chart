# 이벤트 기반 아키텍처 (EDA)

## 개요

시스템은 이벤트 기반 아키텍처를 채택하여 느슨한 결합과 확장성을 확보합니다. 모든 상태 변경은 이벤트로 발행되며, 결과적 일관성 모델을 따릅니다.

## 이벤트 흐름

### 1. 이벤트 발행
- 모든 상태 변경은 이벤트로 발행
- 이벤트는 불변(immutable)하며 순서 보장
- 이벤트 스키마 버전 관리

### 2. 이벤트 구조
```json
{
"eventId": "uuid",
"eventType": "AssetCreated",
"version": "1.0",
"timestamp": "2024-02-16T12:00:00Z",
"aggregateId": "asset-123",
"data": {
// 이벤트 타입별 페이로드
},
"metadata": {
"correlationId": "uuid",
"causationId": "uuid",
"userId": "user-123"
}
}
```

### 3. 이벤트 저장
- MongoDB 이벤트 스토어 사용
- 이벤트 로그 영구 보관
- 이벤트 재생(replay) 지원

## CQRS 구현

### 1. Command 측
- MongoDB: 이벤트 저장소
- 명령 처리 및 이벤트 발행
- 동시성 제어 및 버전 관리

### 2. Query 측
- PostgreSQL: 읽기 전용 뷰
- 이벤트 구독 및 뷰 갱신
- 성능 최적화된 쿼리 모델

## 결과적 일관성

### 1. 일관성 모델
- Write와 Read 모델 간 지연 허용
- 최신성(freshness) 메트릭 모니터링
- 비동기 뷰 갱신

### 2. 보상 처리
- 실패한 이벤트 처리 재시도
- 보상 이벤트 발행
- 불일치 감지 및 복구

## 도메인 이벤트

### Asset 도메인
```
- AssetCreated
- AssetUpdated
- AssetDeleted
- AssetValueChanged
- GoalAssigned
- GoalAchieved
```

### Portfolio 도메인
```
- PortfolioCreated
- PortfolioUpdated
- PortfolioDeleted
- AssetAllocated
- RebalancingRequested
- RebalancingCompleted
```

### Transaction 도메인
```
- TransactionCreated
- TransactionExecuted
- TransactionFailed
- TransactionCancelled
```

## 이벤트 처리

### 1. 이벤트 핸들러
```go
type EventHandler interface {
HandleEvent(ctx context.Context, event Event) error
}
```

### 2. 재시도 정책
- 지수 백오프(exponential backoff)
- 데드레터 큐(DLQ) 사용
- 최대 재시도 횟수 설정

### 3. 멱등성 보장
- 이벤트 ID 기반 중복 처리 방지
- 버전 기반 충돌 감지
- 트랜잭션 로그 유지

## 모니터링

### 1. 이벤트 메트릭
- 이벤트 처리 지연시간
- 실패율 및 재시도 횟수
- 이벤트 큐 크기

### 2. 일관성 메트릭
- 읽기/쓰기 모델 간 지연시간
- 뷰 갱신 성공률
- 불일치 감지 횟수

## 장애 처리

### 1. 서비스 분리
- 이벤트 발행과 처리의 분리
- 실패 격리(failure isolation)
- 부분적 기능 저하 허용

### 2. 복구 전략
- 이벤트 재생을 통한 뷰 재구성
- 스냅샷 기반 빠른 복구
- 수동 개입 절차 정의 