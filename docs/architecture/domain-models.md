# 도메인 모델 상세 설계

## Asset Context (자산 관리)

### 애그리게잇
1. **Asset (자산)**
- 불변식:
- 자산 금액은 음수가 될 수 없음
- 모든 거래는 자산 유형과 일치해야 함
- 자산 식별자는 고유해야 함
- 책임:
- 자산 가치 평가
- 거래 유효성 검증
- 성과 측정
- 이벤트:
- AssetCreated
- AssetUpdated
- AssetDeleted
- AssetValueChanged

2. **Goal (목표)**
- 불변식:
- 목표 금액은 양수여야 함
- 목표 기한은 현재보다 미래여야 함
- 책임:
- 목표 진행률 계산
- 달성 여부 판단
- 이벤트:
- GoalAssigned
- GoalAchieved
- GoalUpdated

### Value Objects
1. **Money**
- 속성: 금액, 통화
- 불변식: 금액은 음수가 될 수 없음
- 연산: 더하기, 빼기, 곱하기, 나누기

2. **TimeRange**
- 속성: 시작일, 종료일
- 불변식: 종료일은 시작일보다 이후여야 함
- 연산: 기간 계산, 포함 여부 확인

### 도메인 서비스
1. **AssetValuationService**
- 책임: 자산 가치 평가
- 정책: 자산 유형별 평가 방법 적용
- 이벤트 구독:
- TransactionExecuted
- MarketPriceUpdated

## Portfolio Context (포트폴리오 관리)

### 애그리게잇
1. **Portfolio (포트폴리오)**
- 불변식:
- 자산 비중의 총합은 100%를 초과할 수 없음
- 포트폴리오는 최소 1개 이상의 자산을 포함해야 함
- 책임:
- 자산 배분 관리
- 리밸런싱 필요성 판단
- 성과 분석
- 이벤트:
- PortfolioCreated
- PortfolioUpdated
- PortfolioDeleted
- AssetAllocated
- RebalancingRequested
- RebalancingCompleted

### Value Objects
1. **Percentage**
- 속성: 값(0-100)
- 불변식: 0-100 사이의 값
- 연산: 더하기, 빼기, 비교

2. **PortfolioAsset**
- 속성: 자산ID, 비중
- 불변식: 비중은 0-100 사이

### 도메인 서비스
1. **PortfolioBalancingService**
- 책임: 포트폴리오 리밸런싱
- 정책: 임계치 기반 리밸런싱 결정
- 이벤트 구독:
- AssetValueChanged
- MarketVolatilityDetected

## Transaction Context (거래 관리)

### 애그리게잇
1. **Transaction (거래)**
- 불변식:
- 거래 금액은 0보다 커야 함
- 거래 날짜는 미래일 수 없음
- 거래 유형은 Buy/Sell 중 하나
- 책임:
- 거래 실행
- 거래 기록
- 자산 상태 변경
- 이벤트:
- TransactionCreated
- TransactionExecuted
- TransactionFailed
- TransactionCancelled

### Value Objects
1. **TransactionType**
- 값: Buy, Sell
- 불변식: 정의된 값만 사용 가능

### 도메인 서비스
1. **TransactionValidationService**
- 책임: 거래 유효성 검증
- 정책: 자산 유형별 거래 규칙 적용
- 이벤트 구독:
- AssetStateChanged
- MarketStatusUpdated

## Monitoring Context (모니터링)

### 애그리게잇
1. **Metric (메트릭)**
- 불변식:
- 메트릭 이름은 고유해야 함
- 타임스탬프는 필수
- 책임:
- 메트릭 수집
- 임계치 확인
- 이벤트:
- MetricCollected
- ThresholdExceeded

2. **Alert (알림)**
- 불변식:
- 알림은 우선순위를 가져야 함
- 중복 알림은 제한됨
- 책임:
- 알림 생성
- 알림 전파
- 이벤트:
- AlertCreated
- AlertTriggered
- AlertResolved

### Value Objects
1. **Threshold**
- 속성: 임계값, 비교연산자
- 불변식: 유효한 비교연산자 사용

### 도메인 서비스
1. **AlertingService**
- 책임: 알림 규칙 평가
- 정책: 알림 우선순위 결정
- 이벤트 구독:
- MetricCollected
- ThresholdExceeded

## 이벤트 스토어

### 이벤트 저장소
- MongoDB 사용
- 이벤트 스키마 버전 관리
- 이벤트 재생 지원

### 이벤트 구독
- 결과적 일관성 모델
- 멱등성 보장
- 재시도 정책

## 읽기 모델 (CQRS)

### Asset 뷰
- 최신 자산 상태
- 성과 지표
- 목표 진행률

### Portfolio 뷰
- 포트폴리오 구성
- 자산 배분 현황
- 리밸런싱 필요성

### Transaction 뷰
- 거래 이력
- 거래 통계
- 성과 분석

### Monitoring 뷰
- 실시간 메트릭
- 알림 상태
- 시스템 건강도 