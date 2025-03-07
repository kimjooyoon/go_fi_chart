# 바운디드 컨텍스트 매핑

## 컨텍스트 간 관계 정의

### Asset ↔ Portfolio
- **관계 유형**: Partnership (협력)
- **통신 방식**:
- 동기: GraphQL API (자산 조회)
- 비동기: 이벤트 스트림 (자산 상태 변경)
- 서비스 간: gRPC (내부 통신)
- **데이터 공유**:
- Portfolio → Asset: GraphQL 쿼리
- Asset → Portfolio: 도메인 이벤트
- **정책**:
- 결과적 일관성 모델 적용
- 이벤트 기반 상태 동기화
- 읽기 전용 복제본 유지

### Portfolio ↔ Transaction
- **관계 유형**: Customer-Supplier
- **통신 방식**:
- 동기: GraphQL API (거래 실행)
- 비동기: 이벤트 스트림 (거래 상태)
- 서비스 간: gRPC (내부 통신)
- **데이터 공유**:
- Portfolio → Transaction: 거래 커맨드
- Transaction → Portfolio: 도메인 이벤트
- **정책**:
- 이벤트 소싱 기반 상태 관리
- 보상 트랜잭션 통한 일관성 유지
- CQRS 패턴 적용

### Transaction ↔ Asset
- **관계 유형**: Conformist
- **통신 방식**:
- 동기: GraphQL API (자산 조회)
- 비동기: 이벤트 스트림 (거래 처리)
- 서비스 간: gRPC (내부 검증)
- **데이터 공유**:
- Transaction → Asset: 거래 이벤트
- Asset → Transaction: 상태 변경 이벤트
- **정책**:
- 이벤트 기반 워크플로우
- 멱등성 보장
- 결과적 일관성 허용

### Monitoring과의 관계
- **관계 유형**: Published Language
- **통신 방식**:
- 비동기: 메트릭 스트림
- 동기: 상태 확인 API
- **데이터 공유**:
- 각 서비스 → Monitoring: 이벤트, 메트릭
- Monitoring → 각 서비스: 알림 이벤트
- **정책**:
- 이벤트 기반 메트릭 수집
- 분산 추적 통합
- 실시간 알림 처리

## 통신 패턴

### 동기 통신
1. **외부 API**
- GraphQL: 클라이언트 통신
- 용도: 데이터 조회, 커맨드 전송
- 타임아웃: 3초
- 재시도: 최대 3회

2. **서비스 간 통신**
- gRPC: 내부 서비스 통신
- 용도: 검증, 조회
- 고려사항:
- 타입 안전성
- 성능 최적화
- 서비스 디스커버리

### 비동기 통신
1. **이벤트 스트림**
- 이벤트 저장소: MongoDB
- 이벤트 발행/구독
- 특성:
- 순서 보장
- 멱등성
- 영구 저장

2. **메시지 큐**
- 용도: 작업 큐
- 특성:
- 신뢰성 있는 전달
- 백프레셔 처리
- 데드레터 큐

## 데이터 일관성

### 커맨드 처리
1. **이벤트 소싱**
- MongoDB 이벤트 저장소
- 버전 기반 동시성 제어
- 이벤트 로그 기반 감사

2. **CQRS**
- Command: MongoDB
- Query: PostgreSQL
- 뷰 갱신: 이벤트 구독

### 데이터 동기화
1. **읽기 모델**
- PostgreSQL 기반 뷰
- 이벤트 기반 갱신
- 캐시 무효화 이벤트

2. **쓰기 모델**
- 이벤트 소싱
- 스냅샷 관리
- 버전 관리

## 장애 처리

### 회복성
1. **서킷 브레이커**
- resilience4j 사용
- 이벤트 기반 상태 전파
- 부분적 장애 허용

2. **폴백 전략**
- 캐시된 데이터 사용
- 읽기 전용 모드
- 점진적 성능 저하

### 모니터링
1. **상태 확인**
- 서비스 헬스 체크
- 이벤트 처리 상태
- 일관성 메트릭

2. **메트릭 수집**
- 이벤트 처리 지연
- 일관성 지연
- 시스템 성능