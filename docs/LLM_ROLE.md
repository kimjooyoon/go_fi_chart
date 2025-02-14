# LLM 개발 프로세스 가이드

## 1. DDD 이벤트 스토밍
### 1.1 도메인 이벤트 식별
- 모든 비즈니스 이벤트를 시간 순서대로 나열
- 각 이벤트의 트리거와 결과 명확히 정의
- 이벤트 간의 인과 관계 파악

### 1.2 커맨드 식별
- 각 이벤트를 발생시키는 커맨드 정의
- 커맨드의 실행 주체 명확화
- 커맨드 실행 조건 정의

### 1.3 애그리게잇 식별
- 이벤트와 커맨드를 중심으로 애그리게잇 경계 설정
- 일관성 경계 정의
- 트랜잭션 범위 설정

## 2. 풍부한 도메인 모델링
### 2.1 도메인 객체 설계
- Value Object와 Entity 구분
- 불변성 보장
- 도메인 규칙 캡슐화

### 2.2 도메인 서비스 정의
- 도메인 객체 간 협력 관계 설계
- 트랜잭션 스크립트 지양
- 도메인 로직의 응집도 극대화

### 2.3 바운디드 컨텍스트 정의
- 컨텍스트 간 경계 명확화
- 컨텍스트 간 통신 규약 정의
- 공유 커널 최소화

## 3. TDD 프로세스
### 3.1 테스트 작성 규칙
- 실패하는 테스트 먼저 작성
- 테스트는 하나의 동작만 검증
- Given-When-Then 패턴 준수
- 테스트 케이스 네이밍 규칙: should_동작_when_조건

### 3.2 테스트 사이클
1. 실패하는 테스트 작성
2. 가장 단순한 구현으로 테스트 통과
3. 리팩토링
4. 테스트 통과 확인
5. 반복

### 3.3 테스트 커버리지
- 도메인 로직 100% 커버리지 목표
- 엣지 케이스 포함
- 실패 케이스 반드시 포함

## 4. 도메인 모델 완성
### 4.1 검증 항목
- 모든 비즈니스 규칙 구현 확인
- 불변식 검증
- 도메인 이벤트 발행 확인

### 4.2 리팩토링 기준
- 중복 제거
- 응집도 향상
- 결합도 감소

## 5. 유스케이스 설계
### 5.1 유스케이스 작성 규칙
- 도메인 모델의 오케스트레이션에만 집중
- 비즈니스 로직 포함 금지
- 트랜잭션 경계 명확화

### 5.2 유스케이스 구조
```go
type UseCase interface {
    Execute(ctx context.Context, command Command) (Result, error)
}
```

### 5.3 유스케이스 테스트
- 도메인 모델 호출 순서 검증
- 트랜잭션 경계 검증
- 에러 처리 검증

## 6. 레포지토리 패턴 구현
### 6.1 레포지토리 인터페이스 설계
- 각 애그리게잇 루트마다 별도의 레포지토리 인터페이스 정의
- 기본 인터페이스 구조:
```go
// Repository 기본 인터페이스
type Repository[T Entity, ID comparable] interface {
    // 기본 CRUD 작업
    Save(ctx context.Context, entity T) error
    FindByID(ctx context.Context, id ID) (T, error)
    Update(ctx context.Context, entity T) error
    Delete(ctx context.Context, id ID) error
    
    // 검색 작업
    FindAll(ctx context.Context, criteria SearchCriteria) ([]T, error)
    FindOne(ctx context.Context, criteria SearchCriteria) (T, error)
    
    // 트랜잭션 관리
    WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// SearchCriteria 검색 조건 인터페이스
type SearchCriteria interface {
    ToQuery() (string, []interface{})
}
```

### 6.2 구현 규칙
1. 책임 범위
   - 순수한 영속성 처리만 담당
   - 도메인 로직 포함 금지
   - 트랜잭션 경계 관리

2. 에러 처리
   - 명확한 에러 타입 정의
   - 데이터베이스 에러를 도메인 에러로 변환
   - 낙관적 락킹 처리

3. 성능 최적화
   - 적절한 인덱스 사용
   - N+1 문제 방지
   - 배치 처리 지원

### 6.3 구현 예시
```go
// AssetRepository 예시
type AssetRepository interface {
    Repository[Asset, string]
    
    // 추가적인 도메인 특화 메서드
    FindByUserID(ctx context.Context, userID string) ([]Asset, error)
    FindByType(ctx context.Context, assetType Type) ([]Asset, error)
    UpdateAmount(ctx context.Context, id string, amount float64) error
}

// TransactionRepository 예시
type TransactionRepository interface {
    Repository[Transaction, string]
    
    // 추가적인 도메인 특화 메서드
    FindByAssetID(ctx context.Context, assetID string) ([]Transaction, error)
    FindByDateRange(ctx context.Context, start, end time.Time) ([]Transaction, error)
    GetTotalAmount(ctx context.Context, assetID string) (float64, error)
}
```

### 6.4 테스트 전략
1. 테스트 계층
   - 단위 테스트: 인메모리 구현체 사용
   - 통합 테스트: 실제 데이터베이스 사용
   - 성능 테스트: 대량 데이터 처리

2. 테스트 데이터 관리
```go
// 테스트 픽스처 예시
type TestFixture struct {
    Assets       []Asset
    Transactions []Transaction
}

func NewTestFixture() *TestFixture {
    return &TestFixture{
        Assets: []Asset{
            // 테스트 데이터
        },
    }
}
```

3. 테스트 케이스 구성
   - 정상 케이스
   - 경계 조건
   - 동시성 처리
   - 에러 처리

## 7. 테스트 자동화
### 7.1 테스트 레벨
- 단위 테스트: 도메인 모델
- 통합 테스트: 레포지토리
- 시스템 테스트: 유스케이스

### 7.2 테스트 데이터
- 테스트 픽스처 관리
- 테스트 데이터 팩토리
- 테스트 데이터 클리너

### 7.3 Mock 전략
```go
type MockRepository struct {
    data map[string]interface{}
}
```

### 7.4 CI/CD 통합
- 모든 PR에 대한 테스트 실행
- 커버리지 리포트 생성
- 성능 테스트 포함

## 8. 품질 기준
### 8.1 코드 품질
- 정적 분석 도구 사용
- 코드 리뷰 체크리스트
- 성능 벤치마크

### 8.2 문서화
- 아키텍처 결정 기록 (ADR)
- API 문서 자동화
- 도메인 용어집 관리

### 8.3 모니터링
- 메트릭 수집
- 로그 집계
- 알림 설정

## 9. 코드 품질 관리

### Go 코드 관리 규칙
1. 코드 수정 후 필수 검증 단계
   - golangci-lint 실행으로 코드 품질 검증
   ```bash
   golangci-lint run ./...
   ```
   - 발견된 문제점 즉시 수정
   - 수정 후 재검증 진행

2. Go 코드 작성 규칙
   - 표준 Go 코드 컨벤션 준수
   - 모든 exported 식별자에 대한 문서화 주석 필수
   - 에러 처리 철저히 구현
   - 테스트 코드 작성 필수

3. 디버깅 절차
   - lint 에러 발생 시 우선 수정
   - 컴파일 에러 해결
   - 런타임 에러 처리 및 로깅 구현
   - 성능 이슈 확인

### 일반적인 코드 관리 규칙
1. 코드 구조화
   - Clean Architecture 원칙 준수
   - 명확한 책임 분리
   - 재사용 가능한 컴포넌트 설계

2. 문서화
   - 코드 변경사항 문서화
   - API 문서 자동화
   - 아키텍처 결정사항 기록

## 10. 개발 프로세스

### 작업 진행 순서
1. 요구사항 분석
2. 설계 검토
3. 구현
4. 코드 품질 검증 (lint)
5. 테스트
6. 문서화
7. 리뷰 요청

### 커뮤니케이션
1. 명확한 설명 제공
2. 문제 해결 과정 공유
3. 대안 제시시 근거 설명

## 11. 도구 사용

### 필수 도구
1. golangci-lint: 코드 품질 검증
2. go test: 단위 테스트
3. go mod: 의존성 관리

### 선택적 도구
1. dlv: 디버깅
2. go-swagger: API 문서화
3. goimports: 임포트 정리 