# Go Financial Chart

[![CI](https://github.com/kimjooyoon/go_fi_chart/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/kimjooyoon/go_fi_chart/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/kimjooyoon/go_fi_chart/branch/main/graph/badge.svg)](https://codecov.io/gh/kimjooyoon/go_fi_chart)
[![Go Report Card](https://goreportcard.com/badge/github.com/kimjooyoon/go_fi_chart)](https://goreportcard.com/report/github.com/kimjooyoon/go_fi_chart)

## 개요

Go Financial Chart는 도메인 주도 설계와 이벤트 기반 아키텍처를 기반으로 구현된 고성능 금융 자산 관리 플랫폼입니다. CQRS 패턴과 결과적 일관성 모델을 채택하여 확장성과 유연성을 보장합니다.

## 핵심 기능

### 자산 관리
- 다양한 금융 자산 유형 지원
- 실시간 자산 가치 평가
- 목표 기반 자산 관리
- 이벤트 기반 상태 추적

### 포트폴리오 관리
- 동적 자산 배분
- 자동 포트폴리오 리밸런싱
- 실시간 성과 분석
- 목표 달성 모니터링

### 거래 관리
- 이벤트 소싱 기반 거래 처리
- 실시간 거래 상태 추적
- 거래 이력 관리
- 결과적 일관성 보장

### 모니터링
- 실시간 시스템 메트릭
- 분산 추적
- 이벤트 처리 모니터링
- 비즈니스 인사이트 분석

## 아키텍처

### 도메인 주도 설계
- 풍부한 도메인 모델
- 명확한 바운디드 컨텍스트
- 도메인 이벤트 중심 설계
- 유비쿼터스 언어 적용

### 이벤트 기반 아키텍처
- 이벤트 소싱 패턴
- CQRS 구현
- 결과적 일관성 모델
- 비동기 이벤트 처리

### 마이크로서비스
- 자율적 서비스
- 독립적 배포
- 격리된 데이터 관리
- 서비스 간 이벤트 통신

## 기술 스택

### 핵심 기술
- Go 1.24.0
- gRPC & GraphQL
- MongoDB & PostgreSQL
- Apache Kafka

### 인프라스트럭처
- Docker & Kubernetes
- Istio Service Mesh
- Prometheus & Grafana
- ELK Stack

## 프로젝트 구조
```
.
├── cmd/                    # 실행 파일
│   └── server/
│       └── main.go
├── docs/                   # 문서
│   ├── architecture/      # 아키텍처 문서
│   │   ├── API_CONTRACTS.md
│   │   ├── BOUNDED_CONTEXTS.md
│   │   ├── CONTEXT_MAP.md
│   │   ├── domain-models.md
│   │   └── event-driven.md
│   ├── development/       # 개발 가이드
│   │   ├── conventions.md
│   │   └── workflow.md
│   ├── event-storming/    # 이벤트 스토밍 결과
│   │   ├── 1.EVENTS.md
│   │   ├── 2.COMMANDS.md
│   │   └── 3.AGGREGATES.md
│   ├── monitoring/        # 모니터링 문서
│   └── services/          # 서비스 문서
├── internal/              # 내부 패키지
│   ├── api/              # API 핸들러
│   ├── config/           # 설정
│   ├── di/               # 의존성 주입
│   ├── domain/           # 도메인 모델
│   │   ├── asset/       # 자산 도메인
│   │   ├── event/       # 이벤트 정의
│   │   └── gamification/# 게이미피케이션
│   └── infrastructure/   # 인프라스트럭처
│       └── events/       # 이벤트 처리
├── metrics/              # 메트릭 수집
│   └── github/          # GitHub 메트릭
├── pkg/                 # 공개 패키지
│   ├── domain/         # 공유 도메인 모델
│   │   └── valueobjects/# 값 객체
│   └── services/       # 공유 서비스
└── services/           # 마이크로서비스
├── analysis/       # 분석 서비스
├── asset/          # 자산 서비스
├── datacollection/ # 데이터 수집
├── gamification/   # 게이미피케이션
├── monitoring/     # 모니터링 서비스
├── portfolio/      # 포트폴리오 서비스
├── transaction/    # 거래 서비스
└── visualization/  # 시각화 서비스
```

## 시작하기

### 요구사항
- Go 1.24.0
- Docker & Docker Compose
- Make

### 로컬 개발 환경
```bash
# 저장소 클론
git clone https://github.com/yourusername/go_fi_chart.git
cd go_fi_chart

# 개발 환경 설정
make setup-dev

# 서비스 실행
make run
```

### 환경 변수
```bash
# 필수 환경 변수
MONGODB_URI=mongodb://...    # MongoDB 연결 문자열
POSTGRES_URI=postgres://...  # PostgreSQL 연결 문자열
KAFKA_BROKERS=localhost:9092 # Kafka 브로커 주소
```

## 문서

### 아키텍처
- [도메인 모델](docs/architecture/domain-models.md)
- [이벤트 기반 아키텍처](docs/architecture/event-driven.md)
- [바운디드 컨텍스트](docs/architecture/context-mapping.md)

### 서비스
- [Asset 서비스](docs/services/asset/README.md)
- [Portfolio 서비스](docs/services/portfolio/README.md)
- [Transaction 서비스](docs/services/transaction/README.md)
- [Monitoring 서비스](docs/services/monitoring/README.md)

### 개발
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