# 모니터링 아키텍처

## 개요

이벤트 기반 아키텍처의 모니터링은 다음 세 가지 주요 영역에 초점을 맞춥니다:
1. 이벤트 흐름 모니터링
2. 서비스 상태 모니터링
3. 비즈니스 메트릭 수집

## 모니터링 계층

### 1. 이벤트 모니터링
- 이벤트 처리 지연 시간
- 이벤트 큐 길이
- 실패한 이벤트 수
- 재시도 횟수
- 이벤트 처리율

### 2. 서비스 모니터링
- 서비스 상태
- API 응답 시간
- 에러율
- 리소스 사용량
- 캐시 히트율

### 3. 비즈니스 메트릭
- 자산 생성/수정/삭제 비율
- 포트폴리오 변경 횟수
- 거래 성공/실패율
- 목표 달성률

## 모니터링 구현

### 1. 메트릭 수집
```yaml
metrics:
# 이벤트 메트릭
- name: event_processing_time
type: histogram
labels:
- event_type
- service
- status

- name: event_queue_length
type: gauge
labels:
- queue_name

# 서비스 메트릭
- name: api_response_time
type: histogram
labels:
- endpoint
- method
- status_code

- name: error_count
type: counter
labels:
- service
- error_type

# 비즈니스 메트릭
- name: asset_operations
type: counter
labels:
- operation
- asset_type
- status
```

### 2. 알림 규칙
```yaml
alerts:
# 이벤트 처리 지연
- name: EventProcessingDelay
condition: event_processing_time > 5s
severity: warning
annotations:
description: "이벤트 처리 지연 발생"

# 이벤트 큐 적체
- name: EventQueueBacklog
condition: event_queue_length > 1000
severity: critical
annotations:
description: "이벤트 큐 적체 발생"

# 서비스 에러율 증가
- name: HighErrorRate
condition: error_rate > 5%
severity: critical
annotations:
description: "서비스 에러율 임계치 초과"
```

### 3. 대시보드 구성
```yaml
dashboards:
# 이벤트 모니터링
- name: "이벤트 흐름"
panels:
- title: "이벤트 처리 시간"
metric: event_processing_time
type: heatmap

- title: "이벤트 큐 상태"
metric: event_queue_length
type: graph

# 서비스 상태
- name: "서비스 상태"
panels:
- title: "API 응답 시간"
metric: api_response_time
type: graph

- title: "에러율"
metric: error_count
type: graph

# 비즈니스 메트릭
- name: "비즈니스 현황"
panels:
- title: "자산 운영 현황"
metric: asset_operations
type: graph

- title: "거래 성공률"
metric: transaction_success_rate
type: gauge
```

## 운영 가이드

### 1. 모니터링 우선순위
1. 이벤트 처리 지연
2. 서비스 가용성
3. 비즈니스 메트릭

### 2. 대응 절차
1. 이벤트 처리 지연
- 이벤트 큐 상태 확인
- 처리기 로그 분석
- 리소스 사용량 확인
- 필요시 스케일 아웃

2. 서비스 장애
- 에러 로그 분석
- 의존성 서비스 상태 확인
- 회로 차단기 상태 확인
- 필요시 롤백 또는 장애 조치

3. 비즈니스 이상
- 관련 이벤트 로그 분석
- 데이터 정합성 확인
- 필요시 수동 개입

### 3. 로그 수집
```yaml
logging:
# 이벤트 로그
- category: event
fields:
- event_id
- event_type
- timestamp
- processing_time
- status

# 서비스 로그
- category: service
fields:
- request_id
- method
- path
- status_code
- response_time

# 비즈니스 로그
- category: business
fields:
- operation
- entity_id
- user_id
- result
- metadata
``` 