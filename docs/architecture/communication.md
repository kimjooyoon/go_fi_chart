# 서비스 간 통신

## 통신 패턴

### 1. 동기식 통신 (REST API)
현재 구현된 주요 통신 방식입니다.

#### Asset 서비스 → Portfolio 서비스
- 목적: 포트폴리오 구성을 위한 자산 정보 제공
- 엔드포인트: `GET /api/v1/assets/{id}`
- 응답: 자산 상세 정보

#### Portfolio 서비스 → Transaction 서비스
- 목적: 포트폴리오 거래 내역 조회
- 엔드포인트: `GET /api/v1/transactions/portfolio/{portfolioId}`
- 응답: 포트폴리오 거래 내역 목록

#### Transaction 서비스 → Asset 서비스
- 목적: 거래 대상 자산 정보 조회
- 엔드포인트: `GET /api/v1/assets/{id}`
- 응답: 자산 상세 정보

### 2. 비동기식 통신 (이벤트 기반)
향후 구현 예정인 통신 방식입니다.

#### 예정된 이벤트
1. 자산 관련 이벤트
- AssetCreated
- AssetUpdated
- AssetDeleted
- AssetValueChanged

2. 포트폴리오 관련 이벤트
- PortfolioCreated
- PortfolioUpdated
- PortfolioDeleted
- AssetAllocationChanged

3. 거래 관련 이벤트
- TransactionCreated
- TransactionCompleted
- TransactionFailed
- TransactionCancelled

## 에러 처리

### 1. 동기식 통신 에러
- 4xx 에러: 클라이언트 측 에러
- 400: 잘못된 요청
- 401: 인증 실패
- 403: 권한 없음
- 404: 리소스 없음

- 5xx 에러: 서버 측 에러
- 500: 내부 서버 에러
- 503: 서비스 불가

### 2. 비동기식 통신 에러 (예정)
- 이벤트 발행 실패
- 이벤트 처리 실패
- 재시도 메커니즘
- Dead Letter Queue

## 서킷 브레이커 (예정)

### 구현 예정 기능
1. 장애 감지
- 에러율 모니터링
- 응답 시간 모니터링
- 타임아웃 설정

2. 상태 관리
- Closed: 정상 동작
- Open: 차단 상태
- Half-Open: 부분 허용

3. 폴백 메커니즘
- 캐시된 데이터 사용
- 기본값 반환
- 대체 서비스 사용

## 모니터링

### 1. 통신 메트릭
- 요청 수
- 응답 시간
- 에러율
- 타임아웃 수

### 2. 알림
- 높은 에러율 발생
- 느린 응답 시간
- 서비스 불가 상태
- 비정상적인 트래픽

## API 버저닝
현재 v1을 사용 중이며, 향후 변경 시 다음 규칙을 따릅니다:

### 버전 관리 규칙
1. URI 경로 버저닝
- 예: `/api/v1/assets`

2. 하위 호환성 유지
- 기존 API 유지
- 점진적 마이그레이션
- 충분한 공지 기간