# 시스템 아키텍처 개요

## 전략적 설계 (Strategic Design)

### 도메인 분석
Go Fi Chart는 금융 자산 관리 도메인을 다음과 같은 하위 도메인으로 구분합니다:

1. **핵심 도메인 (Core Domain)**
- 자산 관리 (Asset)
- 포트폴리오 관리 (Portfolio)
- 거래 관리 (Transaction)

2. **지원 도메인 (Supporting Domain)**
- 모니터링 및 알림
- 성과 분석
- 리스크 관리

3. **일반 도메인 (Generic Domain)**
- 사용자 인증/인가
- 로깅
- 메트릭 수집

### 바운디드 컨텍스트 (Bounded Context)

각 도메인은 명확한 경계를 가진 바운디드 컨텍스트로 구현됩니다:

1. **Asset Context** (포트: 8080)
- 책임:
- 자산 생명주기 관리
- 가치 평가
- 성과 측정
- 집계 루트 (Aggregate Root):
- Asset
- Performance
- Goal

2. **Portfolio Context** (포트: 8081)
- 책임:
- 포트폴리오 구성 관리
- 자산 배분 전략
- 리밸런싱
- 집계 루트:
- Portfolio
- Strategy
- Allocation

3. **Transaction Context** (포트: 8082)
- 책임:
- 거래 실행
- 거래 이력 관리
- 정산
- 집계 루트:
- Transaction
- Settlement
- Balance

4. **Monitoring Context** (포트: 8083)
- 책임:
- 시스템 상태 모니터링
- 메트릭 수집
- 알림 관리
- 집계 루트:
- Metric
- Alert
- HealthCheck

### 컨텍스트 매핑 (Context Mapping)

컨텍스트 간 관계:

1. **Asset ↔ Portfolio**
- 관계: Partnership (협력)
- 통신: 동기식 API
- 데이터: 자산 정보, 가치 평가

2. **Portfolio ↔ Transaction**
- 관계: Customer-Supplier
- 통신: 이벤트 기반 (계획)
- 데이터: 거래 요청, 실행 결과

3. **Transaction ↔ Asset**
- 관계: Conformist
- 통신: 동기식 API
- 데이터: 자산 상태 변경

## 전술적 설계 (Tactical Design)

### 도메인 모델 패턴

1. **값 객체 (Value Objects)**
- Money: 화폐 값
- Percentage: 비율
- TimeRange: 기간
- AssetType: 자산 유형

2. **엔티티 (Entities)**
- Asset
- Portfolio
- Transaction
- Alert

3. **집계 (Aggregates)**
- 트랜잭션 일관성 경계 정의
- 불변성 규칙 적용
- 동시성 제어

4. **도메인 이벤트 (Domain Events)**
- AssetCreated
- PortfolioRebalanced
- TransactionExecuted
- AlertTriggered

5. **도메인 서비스 (Domain Services)**
- PortfolioBalancingService
- AssetValuationService
- RiskAssessmentService

## 기술 아키텍처

### 구현 기술

1. **프레임워크 & 라이브러리**
- Go 1.24.0: 주 개발 언어
- Chi: 경량 웹 프레임워크
- Prometheus & Grafana: 모니터링

2. **영속성**
- 현재: 인메모리 저장소
- 계획:
- PostgreSQL: 트랜잭션 데이터
- Event Store: 도메인 이벤트
- Redis: 캐싱

3. **통신**
- 현재: REST API
- 계획: 이벤트 기반 통신

### 배포 아키텍처

1. **컨테이너화**
- Docker: 서비스 컨테이너화
- Kubernetes: 오케스트레이션 (계획)

2. **확장성**
- 수평적 확장: 서비스별 독립 스케일링
- 자동 스케일링: 부하 기반 (계획)

3. **모니터링**
- 메트릭: Prometheus
- 대시보드: Grafana
- 로깅: 구조화된 JSON

### 보안

1. **인증 & 인가**
- JWT 기반 인증 (계획)
- RBAC 기반 권한 관리 (계획)

2. **데이터 보안**
- HTTPS
- 암호화된 저장소
- 감사 로깅

## 진화 전략

1. **현재 구현**
- 기본 도메인 모델
- 인메모리 저장소
- REST API 통신
- 기본 모니터링

2. **단기 계획**
- 이벤트 소싱 도입
- 영구 저장소 마이그레이션
- 보안 강화

3. **장기 계획**
- CQRS 패턴 적용
- 메시지 큐 도입
- 분산 트랜잭션 지원