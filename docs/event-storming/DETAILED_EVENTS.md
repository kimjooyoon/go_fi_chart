# 도메인별 이벤트 스토밍 상세

## 1. 자산 관리 도메인

### 도메인 이벤트
1. **자산 관련**
- `AssetCreated`: 새로운 자산 생성됨
- `AssetUpdated`: 자산 정보 업데이트됨
- `AssetDeleted`: 자산 삭제됨
- `AssetPriceChanged`: 자산 가격 변경됨

2. **포트폴리오 관련**
- `PortfolioCreated`: 새로운 포트폴리오 생성됨
- `PortfolioUpdated`: 포트폴리오 정보 업데이트됨
- `PortfolioDeleted`: 포트폴리오 삭제됨
- `PortfolioRebalanced`: 포트폴리오 재조정됨
- `AssetAllocationChanged`: 자산 할당 비율 변경됨

3. **거래 관련**
- `TransactionCreated`: 새로운 거래 기록됨
- `TransactionExecuted`: 거래 실행 완료됨
- `TransactionFailed`: 거래 실행 실패함

### 커맨드
1. **자산 관리**
- `CreateAsset`
- `UpdateAsset`
- `DeleteAsset`
- `UpdateAssetPrice`

2. **포트폴리오 관리**
- `CreatePortfolio`
- `UpdatePortfolio`
- `DeletePortfolio`
- `RebalancePortfolio`
- `UpdateAssetAllocation`

3. **거래 처리**
- `CreateTransaction`
- `ExecuteTransaction`
- `CancelTransaction`

### 정책
1. **자산 가격 업데이트**
- When: `AssetPriceChanged`
- Then: 관련 포트폴리오 가치 재계산

2. **포트폴리오 재조정**
- When: `PortfolioRebalanced`
- Then: 필요한 거래 목록 생성

3. **거래 실행**
- When: `TransactionCreated`
- Then: 거래 실행 및 포트폴리오 업데이트

## 2. 분석 도메인

### 도메인 이벤트
1. **시계열 데이터**
- `TimeSeriesDataReceived`: 새로운 시계열 데이터 수신됨
- `TimeSeriesDataProcessed`: 시계열 데이터 처리됨
- `TimeSeriesAnalysisCompleted`: 시계열 분석 완료됨

2. **포트폴리오 분석**
- `PortfolioAnalysisStarted`: 포트폴리오 분석 시작됨
- `PortfolioAnalysisCompleted`: 포트폴리오 분석 완료됨
- `RiskAnalysisCompleted`: 리스크 분석 완료됨

3. **리스크 분석**
- `RiskMetricsCalculated`: 리스크 지표 계산됨
- `RiskLevelChanged`: 리스크 수준 변경됨
- `RiskAlertTriggered`: 리스크 알림 발생됨

### 커맨드
1. **시계열 분석**
- `ProcessTimeSeriesData`
- `AnalyzeTimeSeries`
- `GenerateTimeSeriesReport`

2. **포트폴리오 분석**
- `AnalyzePortfolio`
- `CalculateReturns`
- `OptimizePortfolio`

3. **리스크 분석**
- `CalculateRiskMetrics`
- `AssessRiskLevel`
- `GenerateRiskReport`

### 정책
1. **시계열 데이터 처리**
- When: `TimeSeriesDataReceived`
- Then: 데이터 정규화 및 분석 시작

2. **포트폴리오 최적화**
- When: `PortfolioAnalysisCompleted`
- Then: 최적화 제안 생성

3. **리스크 관리**
- When: `RiskLevelChanged`
- Then: 필요시 알림 생성

## 3. 모니터링 도메인

### 도메인 이벤트
1. **메트릭 관련**
- `MetricCollected`: 메트릭 수집됨
- `MetricThresholdExceeded`: 메트릭 임계값 초과됨
- `MetricNormalized`: 메트릭 정상화됨

2. **알림 관련**
- `AlertCreated`: 새로운 알림 생성됨
- `AlertAcknowledged`: 알림 확인됨
- `AlertResolved`: 알림 해결됨
- `AlertEscalated`: 알림 에스컬레이션됨

3. **상태 체크**
- `HealthCheckExecuted`: 상태 체크 실행됨
- `ServiceStatusChanged`: 서비스 상태 변경됨
- `SystemRecovered`: 시스템 복구됨

### 커맨드
1. **메트릭 관리**
- `CollectMetric`
- `SetMetricThreshold`
- `NormalizeMetric`

2. **알림 관리**
- `CreateAlert`
- `AcknowledgeAlert`
- `ResolveAlert`
- `EscalateAlert`

3. **상태 관리**
- `ExecuteHealthCheck`
- `UpdateServiceStatus`
- `TriggerRecovery`

### 정책
1. **메트릭 모니터링**
- When: `MetricThresholdExceeded`
- Then: 알림 생성

2. **알림 처리**
- When: `AlertCreated`
- Then: 알림 전파 및 필요시 에스컬레이션

3. **상태 관리**
- When: `ServiceStatusChanged`
- Then: 필요한 복구 작업 시작

## 4. 데이터 수집 도메인

### 도메인 이벤트
1. **데이터 소스**
- `DataSourceConnected`: 데이터 소스 연결됨
- `DataSourceDisconnected`: 데이터 소스 연결 해제됨
- `DataSourceError`: 데이터 소스 오류 발생

2. **파이프라인**
- `DataCollectionStarted`: 데이터 수집 시작됨
- `DataTransformationCompleted`: 데이터 변환 완료됨
- `DataLoadingCompleted`: 데이터 적재 완료됨
- `PipelineStatusChanged`: 파이프라인 상태 변경됨

### 커맨드
1. **데이터 소스 관리**
- `ConnectDataSource`
- `DisconnectDataSource`
- `ValidateDataSource`

2. **파이프라인 관리**
- `StartPipeline`
- `PausePipeline`
- `ResumePipeline`
- `StopPipeline`

### 정책
1. **데이터 소스 관리**
- When: `DataSourceError`
- Then: 재연결 시도 및 알림 생성

2. **파이프라인 관리**
- When: `PipelineStatusChanged`
- Then: 필요한 복구 작업 시작

## 5. 시각화 도메인

### 도메인 이벤트
1. **차트 관련**
- `ChartCreated`: 새로운 차트 생성됨
- `ChartUpdated`: 차트 업데이트됨
- `ChartDataRefreshed`: 차트 데이터 갱신됨

2. **대시보드 관련**
- `DashboardCreated`: 새로운 대시보드 생성됨
- `DashboardUpdated`: 대시보드 업데이트됨
- `WidgetAdded`: 위젯 추가됨
- `WidgetRemoved`: 위젯 제거됨

### 커맨드
1. **차트 관리**
- `CreateChart`
- `UpdateChart`
- `RefreshChartData`

2. **대시보드 관리**
- `CreateDashboard`
- `UpdateDashboard`
- `AddWidget`
- `RemoveWidget`

### 정책
1. **차트 관리**
- When: `ChartDataRefreshed`
- Then: 관련 대시보드 업데이트

2. **대시보드 관리**
- When: `WidgetAdded` or `WidgetRemoved`
- Then: 대시보드 레이아웃 재계산

## 6. 게이미피케이션 도메인

### 도메인 이벤트
1. **프로필 관련**
- `ProfileCreated`: 새로운 프로필 생성됨
- `LeveledUp`: 레벨 상승됨
- `ExperienceGained`: 경험치 획득됨
- `BadgeEarned`: 뱃지 획득됨

2. **보상 관련**
- `RewardGranted`: 보상 지급됨
- `AchievementUnlocked`: 업적 달성됨
- `DailyStreakUpdated`: 일일 연속 기록 업데이트됨

### 커맨드
1. **프로필 관리**
- `CreateProfile`
- `UpdateLevel`
- `AddExperience`
- `GrantBadge`

2. **보상 관리**
- `GrantReward`
- `UnlockAchievement`
- `UpdateDailyStreak`

### 정책
1. **레벨 관리**
- When: `ExperienceGained`
- Then: 레벨업 조건 확인 및 처리

2. **보상 관리**
- When: `AchievementUnlocked`
- Then: 관련 보상 지급