# 메트릭 도메인

## 개요

메트릭 도메인은 시스템의 다양한 측정값을 수집하고 관리하는 책임을 가집니다.

## 도메인 모델

### 메트릭 (Metric)

메트릭은 시스템의 특정 측정값을 나타냅니다:

```go
type Metric interface {
Name() string        // 메트릭의 이름
Type() Type         // 메트릭의 타입
Value() Value       // 메트릭의 값
Description() string // 메트릭의 설명
}
```

### 메트릭 타입 (Type)

메트릭은 네 가지 기본 타입을 지원합니다:

- `Counter`: 누적되는 값 (예: 요청 수)
- `Gauge`: 현재 상태 값 (예: 메모리 사용량)
- `Histogram`: 값의 분포 (예: 응답 시간 분포)
- `Summary`: 값의 요약 통계 (예: 응답 시간 백분위)

### 메트릭 값 (Value)

메트릭 값은 다음 정보를 포함합니다:

```go
type Value struct {
Raw       float64           // 실제 측정값
Labels    map[string]string // 메트릭 레이블
Timestamp time.Time         // 측정 시간
}
```

## 컬렉터 (Collector)

메트릭 수집을 담당하는 인터페이스입니다:

```go
type Collector interface {
Collect(ctx context.Context) ([]Metric, error)
AddMetric(metric Metric) error
Reset()
}
```

### 구현체

1. BaseCollector
- 모든 컬렉터의 기본 구현 제공
- 스레드 안전한 메트릭 관리
- 이벤트 발행 기능

2. SimpleCollector
- BaseCollector를 확장한 기본 구현체
- 일반적인 메트릭 수집에 사용

## 이벤트

메트릭 도메인은 다음 이벤트를 발행합니다:

- `metric.collected`: 메트릭 수집 완료 시 발행

## 사용 예시

### 기본 메트릭 생성

```go
value := domain.NewValue(42.0, map[string]string{
"service": "api",
"method": "GET",
})

metric := domain.NewBaseMetric(
"http_requests_total",
domain.TypeCounter,
value,
"Total number of HTTP requests",
)
```

### 메트릭 수집

```go
collector := collectors.NewSimpleCollector(publisher)
collector.AddMetric(metric)
metrics, err := collector.Collect(ctx)
```

## 확장 가이드

1. 새로운 메트릭 타입 추가:
- `Type` 상수에 새로운 타입 추가
- 필요한 경우 `Value` 구조체 확장

2. 커스텀 컬렉터 구현:
- `BaseCollector` 임베딩
- 필요한 추가 기능 구현

3. 새로운 이벤트 추가:
- 이벤트 타입 상수 정의
- Publisher 인터페이스 구현