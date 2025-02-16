# 바운디드 컨텍스트

## 개요

Go Fi Chart는 다음과 같은 바운디드 컨텍스트로 구성됩니다:

## 1. 자산 관리 서비스 (Asset Management)

### 책임
- 자산 CRUD 작업
- 포트폴리오 관리
- 거래 처리 및 기록

### 도메인 모델
- **Asset**: 자산 정보
- ID, 이름, 유형, 가치 등
- **Portfolio**: 포트폴리오 구성
- 자산 목록, 할당 비율 등
- **Transaction**: 거래 기록
- 거래 유형, 금액, 시간 등

### API 엔드포인트
- `/api/v1/assets`
- `/api/v1/portfolios`
- `/api/v1/transactions`

## 2. 분석 서비스 (Analysis)

### 책임
- 시계열 데이터 분석
- 포트폴리오 최적화
- 리스크 분석

### 도메인 모델
- **TimeSeriesData**: 시계열 데이터
- 시간, 값, 메타데이터
- **PortfolioAnalysis**: 포트폴리오 분석
- 성과 지표, 최적화 결과
- **RiskMetrics**: 리스크 지표
- VaR, 변동성, 상관관계

### API 엔드포인트
- `/api/v1/analysis/timeseries`
- `/api/v1/analysis/portfolio`
- `/api/v1/analysis/risk`

## 3. 모니터링 서비스 (Monitoring)

### 책임
- 시스템 상태 모니터링
- 메트릭 수집
- 알림 관리

### 도메인 모델
- **Metric**: 메트릭 정보
- 이름, 값, 타입, 레이블
- **Alert**: 알림 정보
- 수준, 메시지, 상태
- **HealthCheck**: 상태 체크
- 서비스 상태, 에러 정보

### API 엔드포인트
- `/api/v1/metrics`
- `/api/v1/alerts`
- `/api/v1/health`

## 4. 데이터 수집 서비스 (Data Collection)

### 책임
- 실시간 데이터 수집
- ETL 프로세스 관리
- 데이터 품질 관리

### 도메인 모델
- **DataSource**: 데이터 소스
- 소스 정보, 연결 설정
- **DataPipeline**: 데이터 파이프라인
- 처리 단계, 상태
- **DataQuality**: 데이터 품질
- 검증 규칙, 품질 지표

### API 엔드포인트
- `/api/v1/datasources`
- `/api/v1/pipelines`
- `/api/v1/quality`

## 5. 시각화 서비스 (Visualization)

### 책임
- 차트 생성
- 대시보드 관리
- 데이터 탐색 인터페이스

### 도메인 모델
- **Chart**: 차트 정보
- 타입, 데이터, 설정
- **Dashboard**: 대시보드
- 레이아웃, 위젯
- **DataExplorer**: 데이터 탐색
- 쿼리, 필터, 뷰

### API 엔드포인트
- `/api/v1/charts`
- `/api/v1/dashboards`
- `/api/v1/explorer`

## 6. 게이미피케이션 서비스 (Gamification)

### 책임
- 사용자 참여 관리
- 보상 시스템
- 진행 상황 추적

### 도메인 모델
- **Profile**: 사용자 프로필
- 레벨, 경험치, 뱃지
- **Reward**: 보상
- 타입, 조건, 가치
- **Progress**: 진행 상황
- 목표, 달성도, 스트릭

### API 엔드포인트
- `/api/v1/profiles`
- `/api/v1/rewards`
- `/api/v1/progress`

## 서비스 간 통신

### 이벤트 기반 통신
- **Event Bus**: 도메인 이벤트 발행/구독
- **Message Queue**: 비동기 작업 처리

### 동기 통신
- **HTTP/gRPC**: 서비스 간 직접 통신
- **API Gateway**: 클라이언트 요청 라우팅

## 데이터 일관성

### SAGA 패턴
- 분산 트랜잭션 관리
- 보상 트랜잭션 구현

### 이벤트 소싱
- 상태 변경 이벤트 기록
- 이벤트 스토어 구현 