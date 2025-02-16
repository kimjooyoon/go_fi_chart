# 메트릭 수집 시스템

## 개요

모니터링 시스템의 메트릭 수집 부분은 다음과 같은 컨텍스트로 구분됩니다:

1. 기본 메트릭 (Base Metrics)
2. GitHub 메트릭 (GitHub Metrics)
3. 메트릭 수집기 (Collectors)

## 아키텍처

### 메트릭 인터페이스

```go
type Metric interface {
    Name() string
    Type() Type
    Value() Value
    Description() string
}
```

모든 메트릭은 이 인터페이스를 구현해야 합니다. 이를 통해:

- 메트릭의 이름, 타입, 값, 설명을 일관된 방식으로 제공
- 다양한 타입의 메트릭을 동일한 방식으로 처리 가능
- 확장성 있는 메트릭 시스템 구현 가능

### 컬렉터 인터페이스

```go
type Collector interface {
    Collect(ctx context.Context) ([]Metric, error)
}
```

컬렉터는 메트릭을 수집하고 이벤트를 발행하는 역할을 합니다:

- 메트릭 수집 로직 캡슐화
- 이벤트 기반 아키텍처 지원
- 스레드 안전성 보장

## 구현체

### SimpleCollector

기본적인 메트릭 수집기 구현체입니다:

- 일반적인 메트릭 수집에 사용
- 동시성 제어를 위한 뮤텍스 사용
- 이벤트 발행 기능 내장

### GitHub Collector

GitHub 관련 메트릭을 수집하는 특화된 구현체:

- GitHub 액션 상태 메트릭 수집
- GitHub 액션 실행 시간 메트릭 수집
- 상태와 시간을 별도의 메트릭으로 관리

## 이벤트 시스템

메트릭 수집 시스템은 이벤트 기반으로 동작합니다:

- `TypeMetricCollected` 이벤트 발행
- Publisher를 통한 이벤트 전파
- 비동기 처리 지원

## 테스트

각 컴포넌트는 다음과 같은 테스트를 포함합니다:

1. 단위 테스트
    - 메트릭 생성 및 수집 검증
    - 이벤트 발행 검증
    - 에러 처리 검증

2. 동시성 테스트
    - 스레드 안전성 검증
    - 경쟁 상태 방지 확인

## 사용 예시

### 기본 메트릭 수집

```go
collector := NewSimpleCollector(publisher)
metric := NewBaseMetric("test_metric", TypeGauge, NewValue(42.0), "Test metric")
collector.AddMetric(metric)
metrics, err := collector.Collect(ctx)
```

### GitHub 메트릭 수집

```go
collector := NewCollector(publisher)
collector.AddActionStatusMetric("workflow", ActionStatusSuccess)
collector.AddActionDurationMetric("workflow", 10 * time.Second)
metrics, err := collector.Collect(ctx)
```

## 확장 가이드

새로운 메트릭 타입 추가:

1. `Metric` 인터페이스 구현
2. 필요한 경우 전용 컬렉터 구현
3. 테스트 코드 작성
4. 문서화 