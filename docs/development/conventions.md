# 개발 컨벤션

## 도메인 모델링

### 1. 애그리게잇 설계
- 트랜잭션 일관성 경계 정의
- 이벤트 발행 책임 할당
- 불변 조건 관리

### 2. 이벤트 모델링
- 도메인 이벤트 명명 규칙: `{엔티티}{상태변경}`
- 이벤트 버전 관리
- 이벤트 스키마 문서화

### 3. CQRS 패턴
- Command 모델: 도메인 중심 설계
- Query 모델: 성능 최적화
- 이벤트 핸들러 구현

## 코드 구조

### 1. 프로젝트 레이아웃
```
service/
├── cmd/
│   └── server/
├── internal/
│   ├── domain/         # 도메인 모델, 이벤트
│   ├── application/    # 유스케이스, 커맨드 핸들러
│   ├── infrastructure/ # 저장소, 메시징
│   └── api/           # API 엔드포인트
└── pkg/               # 공개 패키지
```

### 2. 도메인 패키지 구조
```
domain/
├── model.go           # 도메인 모델
├── events.go         # 도메인 이벤트
├── commands.go       # 커맨드 정의
├── repository.go     # 저장소 인터페이스
└── service.go        # 도메인 서비스
```

### 3. 네이밍 컨벤션
- 이벤트: `{Entity}{Event}Event`
- 커맨드: `{Action}{Entity}Command`
- 핸들러: `{Entity}{Event}Handler`

## 테스트

### 1. 테스트 계층
- 단위 테스트: 도메인 모델, 이벤트
- 통합 테스트: 이벤트 흐름, 저장소
- E2E 테스트: API, 이벤트 처리

### 2. 이벤트 테스트
```go
func TestAssetCreatedEvent_Should_Update_ReadModel(t *testing.T) {
// Given
// When
// Then
}
```

### 3. 테스트 데이터
- 테스트 픽스처 관리
- 이벤트 스트림 모킹
- 저장소 격리

## 에러 처리

### 1. 도메인 에러
```go
type DomainError struct {
Code    string
Message string
Details map[string]interface{}
}
```

### 2. 이벤트 처리 에러
- 재시도 가능 여부 표시
- 컨텍스트 정보 포함
- 에러 로깅 정책

## 로깅

### 1. 이벤트 로깅
- 이벤트 메타데이터 포함
- 상관관계 ID 추적
- 성능 메트릭 수집

### 2. 구조화된 로깅
```json
{
"level": "INFO",
"event_id": "uuid",
"event_type": "AssetCreated",
"aggregate_id": "asset-123",
"correlation_id": "uuid",
"timestamp": "2024-02-16T12:00:00Z"
}
```

## 버전 관리

### 1. 이벤트 버전 관리
- 이벤트 스키마 버전
- 마이그레이션 전략
- 하위 호환성 유지

### 2. API 버전 관리
- API 버전 정책
- 변경 이력 관리
- 클라이언트 마이그레이션

### 3. 브랜치 전략
- main: 프로덕션
- develop: 개발
- feature/*: 기능
- release/*: 릴리스