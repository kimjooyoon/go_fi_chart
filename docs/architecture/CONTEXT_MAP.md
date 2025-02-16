# 컨텍스트 맵

## 바운디드 컨텍스트 구조

### 1. 핵심 도메인 (Core Domain)

#### 자산 관리 컨텍스트
- **애그리게잇**:
- Asset (루트)
- Portfolio (루트)
- Transaction (루트)
- **저장소**: Document Store (MongoDB)
- **이벤트**: AssetCreated, PortfolioRebalanced, TransactionExecuted

#### 분석 엔진 컨텍스트
- **애그리게잇**:
- Analysis (루트)
- TimeSeries (루트)
- RiskMetrics (루트)
- **저장소**: Time Series DB (InfluxDB)
- **이벤트**: AnalysisCompleted, RiskLevelChanged

### 2. 지원 도메인 (Supporting Domain)

#### 데이터 수집 컨텍스트
- **애그리게잇**:
- DataSource (루트)
- Pipeline (루트)
- **저장소**: Document Store (MongoDB)
- **이벤트**: DataCollected, PipelineExecuted

#### 모니터링 컨텍스트
- **애그리게잇**:
- Metric (루트)
- Alert (루트)
- HealthCheck (루트)
- **저장소**: Time Series DB (InfluxDB)
- **이벤트**: MetricCollected, AlertTriggered

### 3. 일반 도메인 (Generic Domain)

#### 시각화 컨텍스트
- **애그리게잇**:
- Chart (루트)
- Dashboard (루트)
- **저장소**: Document Store (MongoDB)
- **이벤트**: ChartCreated, DashboardUpdated

#### 게이미피케이션 컨텍스트
- **애그리게잇**:
- Profile (루트)
- Achievement (루트)
- **저장소**: Document Store (MongoDB)
- **이벤트**: ProfileLeveledUp, AchievementUnlocked

## 컨텍스트 간 관계

### 1. 파트너십 (Partnership)
- **자산 관리 ↔ 분석 엔진**
- 이벤트 기반 통신
- 실시간 데이터 동기화
- 일관성 보장을 위한 이벤트 소싱

### 2. 공유 커널 (Shared Kernel)
- **자산 관리 ↔ 데이터 수집**
- 공유 도메인 이벤트
- 공통 값 객체 (Value Objects)
- 시장 데이터 스키마

### 3. 고객-공급자 (Customer-Supplier)
- **데이터 수집 → 자산 관리**
- 시장 데이터 스트림
- 실시간 가격 정보
- 품질 보증 계약

### 4. 준수자 (Conformist)
- **게이미피케이션 → 자산 관리**
- 이벤트 구독
- 읽기 전용 모델
- 단방향 의존성

## 통신 패턴

### 1. 이벤트 기반 통신
- Apache Kafka 사용
- 이벤트 스토어로 EventStoreDB 사용
- 도메인 이벤트 버전 관리
- 이벤트 스키마 관리

### 2. 명령 처리
- 비동기 명령 패턴
- 명령 유효성 검증
- 실패 처리 및 보상 트랜잭션

### 3. 쿼리 처리
- CQRS 패턴 적용
- 읽기 전용 모델
- 캐시 전략 (Redis)

## 데이터 저장소 전략

### 1. 애그리게잇 저장소
- MongoDB 사용
- 문서 기반 저장
- 애그리게잇 단위 트랜잭션
- 낙관적 동시성 제어

### 2. 이벤트 저장소
- EventStoreDB 사용
- 이벤트 소싱
- 스냅샷 관리
- 이벤트 버저닝

### 3. 시계열 데이터
- InfluxDB 사용
- 메트릭 저장
- 시계열 분석
- 데이터 보존 정책

### 4. 캐시 계층
- Redis 사용
- 읽기 모델 캐시
- 세션 데이터
- 실시간 집계

## 구현 우선순위

### 1단계: 핵심 도메인
- 자산 관리 컨텍스트
- 이벤트 소싱 인프라
- 기본 CQRS 구현

### 2단계: 지원 도메인
- 데이터 수집 컨텍스트
- 모니터링 시스템
- 이벤트 기반 통신

### 3단계: 일반 도메인
- 시각화 서비스
- 게이미피케이션
- 고급 분석 기능