# Monitoring 서비스

## 개요
Monitoring 서비스는 시스템 전반의 상태와 메트릭을 수집하고 모니터링하는 마이크로서비스입니다.

## 주요 기능
- 시스템 메트릭 수집
- 서비스 상태 모니터링
- Prometheus 메트릭 익스포트
- 알림 관리

## 기술 스택
- Go 1.24.0
- Chi 라우터
- Prometheus Client
- In-memory 저장소 (현재)

## API 엔드포인트
### 메트릭 및 모니터링
- `GET /metrics` - Prometheus 메트릭 엔드포인트
- `GET /health` - 서비스 헬스 체크
- `GET /ready` - 서비스 레디니스 체크

## 수집 메트릭
### 시스템 메트릭
- CPU 사용량
- 메모리 사용량
- 디스크 I/O
- 네트워크 트래픽

### 비즈니스 메트릭
- API 요청 수
- 응답 시간
- 에러율
- 거래 처리량

## 의존성
- Asset 서비스: 자산 상태 모니터링
- Portfolio 서비스: 포트폴리오 상태 모니터링
- Transaction 서비스: 거래 활동 모니터링

## 설정
### 환경 변수
- `PORT`: 서비스 포트 (기본값: 8083)
- `SCRAPE_INTERVAL`: 메트릭 수집 주기 (기본값: 15s)
- `RETENTION_PERIOD`: 메트릭 보관 기간 (기본값: 15d)

## 로컬 개발 환경 설정
```bash
# 서비스 실행
go run cmd/server/main.go

# 테스트 실행
go test ./...
```

## 알림 설정
### 알림 규칙
- 서비스 다운
- 높은 에러율
- 느린 응답 시간
- 리소스 부족

## 대시보드
Grafana 대시보드 템플릿이 제공됩니다:
- 시스템 모니터링
- 서비스 성능
- 비즈니스 메트릭 