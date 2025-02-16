# Go Fi Chart

[![CI](https://github.com/kimjooyoon/go_fi_chart/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/kimjooyoon/go_fi_chart/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kimjooyoon/go_fi_chart/branch/main/graph/badge.svg)](https://codecov.io/gh/kimjooyoon/go_fi_chart)
[![Go Report Card](https://goreportcard.com/badge/github.com/kimjooyoon/go_fi_chart)](https://goreportcard.com/report/github.com/kimjooyoon/go_fi_chart)

금융 자산 관리 및 포트폴리오 분석을 위한 도메인 주도 설계(DDD) 기반의 마이크로서비스 플랫폼

## 프로젝트 개요

Go Fi Chart는 금융 자산 관리를 위한 확장 가능한 마이크로서비스 플랫폼입니다. 도메인 주도 설계 원칙을 기반으로 구축되어 복잡한 금융 도메인을 효과적으로 모델링하고 관리합니다.

### 핵심 도메인

- **자산 관리 (Asset Domain)**
- 자산의 생성, 평가, 추적
- 자산 가치 계산 및 성과 측정
- 복잡한 자산 구조 모델링

- **포트폴리오 관리 (Portfolio Domain)**
- 포트폴리오 구성 및 자산 배분
- 리밸런싱 전략
- 성과 분석 및 리스크 관리

- **거래 관리 (Transaction Domain)**
- 거래 실행 및 기록
- 거래 이력 추적
- 정산 및 검증

- **모니터링 (Monitoring Domain)**
- 시스템 상태 추적
- 성능 메트릭 수집
- 이상 징후 감지 및 알림

## 아키텍처

### 바운디드 컨텍스트

각 도메인은 독립된 바운디드 컨텍스트로 구현되어 있으며, 명확한 컨텍스트 경계와 도메인 모델을 가집니다:

1. **Asset 서비스** (포트: 8080)
- 자산 애그리게잇
- 가치 평가 정책
- 성과 측정 도메인 서비스

2. **Portfolio 서비스** (포트: 8081)
- 포트폴리오 애그리게잇
- 자산 배분 정책
- 리밸런싱 도메인 서비스

3. **Transaction 서비스** (포트: 8082)
- 거래 애그리게잇
- 거래 정책
- 정산 도메인 서비스

4. **Monitoring 서비스** (포트: 8083)
- 메트릭 수집기
- 알림 정책
- 상태 모니터링 서비스

### 기술 스택

- **언어 및 프레임워크**
- Go 1.24.0
- Chi 웹 프레임워크
- Event Sourcing (계획)

- **영속성**
- 인메모리 저장소 (현재)
- PostgreSQL (계획)
- Event Store (계획)

- **모니터링**
- Prometheus
- Grafana

## 시작하기

### 요구사항

- Go 1.24.0 이상
- Make
- Docker (선택사항)

### 로컬 개발 환경 설정

```bash
# 저장소 클론
git clone https://github.com/aske/go_fi_chart.git
cd go_fi_chart

# 의존성 설치
go mod download

# 개발 도구 설치
make setup-dev

# 서비스 실행
make run-services

# 테스트 실행
make test
```

### 도커 환경 (선택사항)

```bash
# 이미지 빌드
make docker-build

# 컨테이너 실행
make docker-run
```

## 프로젝트 구조

```
.
├── docs/                    # 문서
│   ├── architecture/       # 아키텍처 문서
│   ├── development/       # 개발 가이드
│   └── event-storming/    # 이벤트 스토밍 결과
├── pkg/                    # 공유 도메인 모델
│   └── domain/
│       ├── events/        # 도메인 이벤트
│       └── valueobjects/  # 값 객체
└── services/              # 마이크로서비스
├── asset/            # 자산 서비스
├── portfolio/        # 포트폴리오 서비스
├── transaction/      # 거래 서비스
└── monitoring/       # 모니터링 서비스
```

## 도메인 모델링

프로젝트는 다음과 같은 DDD 패턴을 적용합니다:

- **애그리게잇**: 트랜잭션 일관성 경계 정의
- **값 객체**: 불변성과 동등성 보장
- **도메인 이벤트**: 도메인 변경사항 추적
- **도메인 서비스**: 복잡한 도메인 로직 처리
- **리포지토리**: 영속성 추상화

## 문서

- [아키텍처 개요](docs/architecture/overview.md)
- [도메인 모델](docs/architecture/domain.md)
- [개발 컨벤션](docs/development/conventions.md)
- [개발 워크플로우](docs/development/workflow.md)

## 기여하기

1. 도메인 모델 이해하기
2. 이슈 생성 또는 기존 이슈 선택
3. 브랜치 생성 (`feature/*`, `bugfix/*`)
4. 변경사항 커밋
5. Pull Request 생성
6. 코드 리뷰 진행
7. 머지

## 라이선스

MIT License

## 연락처

- 이슈 트래커: [GitHub Issues] 