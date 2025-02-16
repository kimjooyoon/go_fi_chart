# 알림 시스템

## 개요

모니터링 시스템의 알림 부분은 다음과 같은 기능을 제공합니다:

1. 알림 생성 및 관리
2. 알림 수준 분류
3. 알림 전파 및 처리

## 아키텍처

### 알림 인터페이스

```go
type Alert struct {
    ID        string
    Level     AlertLevel
    Source    string
    Message   string
    Timestamp time.Time
    Metadata  map[string]string
}
```

알림은 다음과 같은 정보를 포함합니다:

- 고유 식별자 (ID)
- 심각도 수준 (Level)
- 발생 소스 (Source)
- 알림 메시지 (Message)
- 발생 시간 (Timestamp)
- 추가 정보 (Metadata)

### 알림 수준

```go
type AlertLevel string

const (
    LevelInfo     AlertLevel = "INFO"
    LevelWarning  AlertLevel = "WARNING"
    LevelError    AlertLevel = "ERROR"
    LevelCritical AlertLevel = "CRITICAL"
)
```

알림은 심각도에 따라 4단계로 분류됩니다:

- INFO: 정보성 알림
- WARNING: 주의가 필요한 상황
- ERROR: 오류 발생
- CRITICAL: 긴급 대응 필요

### Notifier 인터페이스

```go
type Notifier interface {
    Notify(ctx context.Context, alert Alert) error
}
```

알림 처리기는 이 인터페이스를 구현해야 합니다:

- 알림 전달 로직 캡슐화
- 컨텍스트 기반 처리
- 에러 처리 지원

## 구현체

### SimpleNotifier

기본적인 알림 처리기 구현체입니다:

- 여러 알림 핸들러 지원
- 동시성 제어
- 이벤트 발행 기능

## 이벤트 시스템

알림 시스템은 이벤트 기반으로 동작합니다:

- `TypeAlertTriggered` 이벤트 발행
- Publisher를 통한 이벤트 전파
- 비동기 처리 지원

## 테스트

각 컴포넌트는 다음과 같은 테스트를 포함합니다:

1. 단위 테스트
    - 알림 생성 및 전달 검증
    - 이벤트 발행 검증
    - 에러 처리 검증

2. 동시성 테스트
    - 스레드 안전성 검증
    - 핸들러 관리 검증

## 사용 예시

### 알림 생성 및 전송

```go
alert := Alert{
    ID:      "alert-1",
    Level:   LevelWarning,
    Source:  "system",
    Message: "High CPU usage detected",
    Metadata: map[string]string{
        "cpu": "85%",
    },
}

notifier := NewSimpleNotifier(publisher)
notifier.Notify(ctx, alert)
```

### 알림 핸들러 등록

```go
handler := &CustomHandler{}
notifier.AddHandler(handler)
```

## 확장 가이드

새로운 알림 처리기 추가:

1. `Notifier` 인터페이스 구현
2. 필요한 경우 새로운 알림 수준 추가
3. 테스트 코드 작성
4. 문서화 