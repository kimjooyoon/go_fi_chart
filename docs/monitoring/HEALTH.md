# 상태 체크 시스템

## 개요

모니터링 시스템의 상태 체크 부분은 다음과 같은 기능을 제공합니다:

1. 시스템 상태 모니터링
2. 상태 정보 관리
3. 상태 변경 추적

## 아키텍처

### 상태 체크 인터페이스

```go
type Checker interface {
Check(ctx context.Context) error
Status() Status
LastError() error
}
```

상태 체크기는 이 인터페이스를 구현해야 합니다:

- 상태 확인 기능
- 현재 상태 조회
- 마지막 에러 정보 제공

### 상태 정보

```go
type Status string

const (
StatusUp   Status = "UP"
StatusDown Status = "DOWN"
)
```

시스템 상태는 두 가지로 분류됩니다:

- UP: 정상 동작 중
- DOWN: 문제 발생

## 구현체

### SimpleChecker

기본적인 상태 체크기 구현체입니다:

- 상태 정보 관리
- 동시성 제어
- 에러 추적

## 동작 방식

1. 상태 체크
    - 주기적인 상태 확인
    - 컨텍스트 기반 실행
    - 에러 발생 시 상태 변경

2. 상태 관리
    - 스레드 안전한 상태 업데이트
    - 에러 정보 저장
    - 상태 변경 이력 관리

## 테스트

각 컴포넌트는 다음과 같은 테스트를 포함합니다:

1. 단위 테스트
    - 상태 체크 검증
    - 상태 변경 검증
    - 에러 처리 검증

2. 동시성 테스트
    - 스레드 안전성 검증
    - 상태 관리 검증

## 사용 예시

### 기본 상태 체크

```go
checker := NewSimpleChecker()
err := checker.Check(ctx)
if err != nil {
log.Printf("Health check failed: %v", err)
}

status := checker.Status()
if status == StatusDown {
log.Printf("System is down: %v", checker.LastError())
}
```

### 커스텀 체크 로직

```go
type CustomChecker struct {
*SimpleChecker
}

func (c *CustomChecker) Check(ctx context.Context) error {
// 커스텀 체크 로직 구현
if err := someCheck(); err != nil {
c.UpdateStatus(StatusDown, err)
return err
}
c.UpdateStatus(StatusUp, nil)
return nil
}
```

## 확장 가이드

새로운 상태 체크기 추가:

1. `Checker` 인터페이스 구현
2. 필요한 경우 새로운 상태 추가
3. 테스트 코드 작성
4. 문서화 