# Go Fi Chart

[![CI](https://github.com/kimjooyoon/go_fi_chart/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/kimjooyoon/go_fi_chart/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kimjooyoon/go_fi_chart/branch/main/graph/badge.svg)](https://codecov.io/gh/kimjooyoon/go_fi_chart)
[![Go Report Card](https://goreportcard.com/badge/github.com/kimjooyoon/go_fi_chart)](https://goreportcard.com/report/github.com/kimjooyoon/go_fi_chart)

금융 데이터 분석 및 시각화 플랫폼

## 프로젝트 목표

Go Fi Chart는 다음과 같은 목표를 가진 금융 데이터 플랫폼입니다:

1. **데이터 통합**
  - 다양한 소스의 금융 데이터 수집
  - 실시간 데이터 처리
  - 데이터 정규화 및 검증

2. **분석 자동화**
  - 자동화된 데이터 분석
  - 패턴 인식 및 예측
  - 커스텀 분석 파이프라인

3. **시각화 도구**
  - 실시간 차트 생성
  - 대화형 데이터 탐색
  - 맞춤형 대시보드

## 핵심 기능

### 1. 데이터 파이프라인

- 실시간 데이터 수집
- ETL 프로세스 자동화
- 데이터 품질 관리

### 2. 분석 엔진

- 시계열 데이터 분석
- 포트폴리오 최적화
- 리스크 분석

### 3. 시각화 도구

- 실시간 차트 렌더링
- 대화형 데이터 탐색
- 맞춤형 대시보드

### 4. [모니터링 시스템](docs/monitoring/README.md)

- [메트릭 수집](docs/monitoring/METRICS.md)
- [알림 관리](docs/monitoring/ALERTS.md)
- [상태 모니터링](docs/monitoring/HEALTH.md)

## 아키텍처

프로젝트는 도메인 주도 설계(DDD)와 이벤트 기반 아키텍처를 채택했습니다.
자세한 내용은 [바운디드 컨텍스트](docs/architecture/BOUNDED_CONTEXTS.md) 문서를 참조하세요.

### 서비스 구성

1. **[자산 관리 서비스](services/asset)**
- 자산 CRUD
- 포트폴리오 관리
- 거래 처리

2. **[분석 서비스](services/analysis)**
- 시계열 데이터 분석
- 포트폴리오 최적화
- 리스크 분석

3. **[모니터링 서비스](services/monitoring)**
- 시스템 상태 관리
- 메트릭 수집
- 알림 처리

4. **[데이터 수집 서비스](services/datacollection)**
- 실시간 데이터 수집
- ETL 프로세스
- 데이터 품질 관리

5. **[시각화 서비스](services/visualization)**
- 차트 생성
- 대시보드 관리
- 데이터 탐색

6. **[게이미피케이션 서비스](services/gamification)**
- 사용자 참여 관리
- 보상 시스템
- 진행 상황 추적

### 공통 컴포넌트

- **[이벤트 시스템](internal/domain/event)**
  - 이벤트 정의
  - 이벤트 처리
  - 이벤트 저장

- **[API Gateway](internal/api)**
- 요청 라우팅
- 인증/인가
- 속도 제한

- **[설정 관리](internal/config)**
- 환경 설정
- 서비스 설정
- 보안 설정

## 시작하기

### 필수 조건

- Go 1.21 이상
- Make
- Docker (선택사항)

### 설치
```bash
# 저장소 클론
git clone https://github.com/username/go_fi_chart.git
cd go_fi_chart

# 의존성 설치
make init

# 테스트 실행
make test

# 서비스 실행
make run
```

## 개발 가이드

### 코드 품질

```bash
# 전체 린트
make lint

# 특정 서비스 린트
make lint-monitoring
```

### 테스트
```bash
# 전체 테스트
make test

# 특정 서비스 테스트
make test-monitoring

# 커버리지 리포트
make coverage
```

### 보안
```bash
# 보안 검사
make security
```

## 문서

- [모니터링 시스템](docs/monitoring/README.md)
  - [메트릭 수집](docs/monitoring/METRICS.md)
  - [알림 관리](docs/monitoring/ALERTS.md)
  - [상태 체크](docs/monitoring/HEALTH.md)
- [이벤트 스토밍](docs/event-storming/README.md)
- [작업 현황](docs/DONE.md)
- [향후 계획](docs/TODO.md)
- [문제 해결](docs/PROBLEMS.md)
- [LLM 역할](docs/LLM_ROLE.md)

## 기여하기

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 라이선스

MIT License 