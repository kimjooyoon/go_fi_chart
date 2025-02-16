# Monitoring 서비스

## 개요
Monitoring 서비스는 시스템 전반의 상태와 메트릭을 수집하고 모니터링하는 마이크로서비스입니다. 이벤트 기반 아키텍처를 사용하여 실시간 메트릭을 수집하고, 결과적 일관성 모델을 통해 다른 서비스와 통합됩니다.

## 핵심 기능
- 시스템 메트릭 수집
- 서비스 상태 모니터링
- 이벤트 처리 모니터링
- 실시간 알림 관리
- 분산 추적
- 비즈니스 인사이트 분석

## 기술 스택
- Go 1.24.0
- Prometheus & Grafana
- OpenTelemetry
- MongoDB (이벤트 저장소)
- PostgreSQL (메트릭 저장소)
- Apache Kafka (이벤트 스트림)

## 모니터링 영역

### 시스템 메트릭
```yaml
metrics:
infrastructure:
- cpu_usage
- memory_usage
- disk_io
- network_traffic

application:
- request_rate
- error_rate
- response_time
- concurrent_users
```

### 이벤트 메트릭
```yaml
event_metrics:
processing:
- event_processing_time
- event_queue_length
- failed_events
- retry_count

consistency:
- event_lag
- sync_delay
- consistency_violations
```

### 비즈니스 메트릭
```yaml
business_metrics:
asset:
- creation_rate
- modification_rate
- valuation_changes

portfolio:
- rebalancing_frequency
- performance_metrics
- goal_achievement

transaction:
- volume
- success_rate
- processing_time
```

## 도메인 이벤트
- MetricCollected
- ThresholdExceeded
- AlertCreated
- AlertResolved
- PerformanceAnomaly
- ServiceHealthChanged

## 알림 시스템

### 알림 우선순위
1. CRITICAL: 즉시 조치 필요
2. HIGH: 1시간 이내 조치
3. MEDIUM: 24시간 이내 조치
4. LOW: 계획된 유지보수

### 알림 채널
- Slack
- Email
- SMS
- PagerDuty

## 대시보드

### 시스템 현황
- 서비스 건강도
- 리소스 사용률
- 에러율 추이
- 성능 지표

### 이벤트 모니터링
- 이벤트 처리율
- 이벤트 지연시간
- 실패율 및 재시도
- 일관성 메트릭

### 비즈니스 인사이트
- 자산 운영 현황
- 포트폴리오 성과
- 거래 성공률
- 목표 달성률

## 운영 관리

### 백업 및 복구
- 메트릭 데이터 백업
- 알림 이력 보관
- 시스템 복구 절차

### 확장성 관리
- 수평적 확장
- 데이터 보관 정책
- 성능 최적화

### 보안
- 메트릭 데이터 암호화
- 접근 제어
- 감사 로그

## 통합

### 서비스 연동
- Asset 서비스
- Portfolio 서비스
- Transaction 서비스

### 외부 시스템
- 클라우드 모니터링
- APM 도구
- 로그 분석 시스템

## 설정

### 환경 변수
```bash
# 필수 환경 변수
MONGODB_URI=mongodb://...     # MongoDB 연결
POSTGRES_URI=postgres://...   # PostgreSQL 연결
KAFKA_BROKERS=localhost:9092  # Kafka 브로커
PROMETHEUS_PORT=9090          # Prometheus 포트
GRAFANA_PORT=3000            # Grafana 포트
```

### 메트릭 설정
```yaml
collection:
interval: 15s
batch_size: 1000
buffer_size: 10000

retention:
metrics: 15d
events: 30d
alerts: 90d
```

## 개발 환경

### 로컬 설정
```bash
# 서비스 실행
make run-monitoring

# 테스트 실행
make test-monitoring

# 메트릭 확인
curl localhost:9090/metrics
```

### 테스트 데이터
- 메트릭 생성기
- 이벤트 시뮬레이터
- 부하 테스트 도구

## 장애 대응

### 모니터링 실패
1. 메트릭 수집 중단
2. 알림 전송 실패
3. 데이터 저장 오류

### 복구 절차
1. 서비스 상태 확인
2. 데이터 정합성 검증
3. 시스템 재구동
4. 알림 재전송 