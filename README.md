# Go Fi Chart

금융 데이터 분석 및 차트 시각화 서비스

## 개요

이 프로젝트는 금융 데이터를 수집, 분석하고 차트로 시각화하는 서비스를 제공합니다.

## 기능

- 금융 데이터 수집 및 저장
- 데이터 분석 및 처리
- 차트 시각화
- 모니터링 시스템

## 아키텍처

프로젝트는 다음과 같은 주요 컴포넌트로 구성됩니다:

### 도메인 레이어

- 자산 관리 (`internal/domain/asset`)
- 이벤트 시스템 (`internal/domain/event`)
- 게이미피케이션 (`internal/domain/gamification`)

### 서비스 레이어

- 모니터링 서비스 (`services/monitoring`)
    - [메트릭 수집 시스템](docs/monitoring/METRICS.md)
    - 알림 시스템
    - 상태 체크

### 인프라스트럭처 레이어

- 이벤트 저장소 (`internal/infrastructure/events`)
- API 서버 (`internal/api`)
- 설정 관리 (`internal/config`)

## 설치 및 실행

```bash
# 의존성 설치
go mod download

# 테스트 실행
make test

# 서비스 실행
make run
```

## 개발 가이드

### 테스트

```bash
# 전체 테스트
make test

# 특정 서비스 테스트
make test-monitoring
```

### 린트

```bash
# 전체 린트
make lint

# 특정 서비스 린트
make lint-monitoring
```

### 보안 검사

```bash
# 전체 보안 검사
make security

# 특정 서비스 보안 검사
make security-monitoring
```

## 문서

- [메트릭 수집 시스템](docs/monitoring/METRICS.md)
- [이벤트 스토밍](docs/event-storming/README.md)
- [완료된 작업](docs/DONE.md)
- [할 일](docs/TODO.md)

## 라이선스

MIT License 