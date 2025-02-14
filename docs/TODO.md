# 할 일 목록

## 우선순위 높음 (P0)
1. [x] 도메인 이벤트 스토밍 진행
   - [x] 핵심 도메인 이벤트 식별
   - [x] 커맨드와 액터 정의
   - [x] 애그리게잇 경계 설정

2. [x] 초기 프로젝트 설정
   - [x] Go 모듈 초기화
   - [x] 기본 디렉토리 구조 설정
   - [x] 의존성 관리 설정
   - [x] 코드 품질 도구 설정 (golangci-lint)

3. [ ] 레포지토리 패턴 재구현
   - [x] 제네릭 기반 기본 Repository 인터페이스 정의
   - [x] Asset 도메인 레포지토리 리팩토링
   - [x] ID 생성기 구현
     - [x] UUID 기반 ID 생성 유틸리티 작성
     - [x] ID 생성 테스트 코드 작성
   - [x] 테스트 픽스처 구현
     - [x] Asset 테스트 데이터 정의
     - [x] Transaction 테스트 데이터 정의
     - [x] Portfolio 테스트 데이터 정의
   - [ ] 인메모리 레포지토리 구현 (테스트용)
     - [ ] Asset 인메모리 레포지토리
     - [ ] Transaction 인메모리 레포지토리
     - [ ] Portfolio 인메모리 레포지토리

4. [ ] Value Object 구현
   - [ ] Money Value Object
     - [ ] 기본 구조 정의
     - [ ] 연산 메서드 구현
     - [ ] 통화 변환 지원
   - [ ] Percentage Value Object
     - [ ] 기본 구조 정의
     - [ ] 유효성 검증
     - [ ] 연산 메서드 구현
   - [ ] ID Value Object
     - [ ] 기본 구조 정의
     - [ ] 유효성 검증
     - [ ] 생성 메서드 구현

5. [ ] Entity 보완
   - [ ] Asset Entity
     - [ ] Money Value Object 통합
     - [ ] 유효성 검증 추가
     - [ ] 도메인 이벤트 발행 구현
   - [ ] Transaction Entity
     - [ ] Money Value Object 통합
     - [ ] 유효성 검증 추가
     - [ ] 도메인 이벤트 발행 구현
   - [ ] Portfolio Entity
     - [ ] Percentage Value Object 통합
     - [ ] 유효성 검증 추가
     - [ ] 도메인 이벤트 발행 구현

6. [ ] 도메인 서비스 구현
   - [ ] PortfolioBalanceService
     - [ ] 포트폴리오 밸런싱 로직 구현
     - [ ] 리밸런싱 필요성 계산
     - [ ] 리밸런싱 제안 생성
   - [ ] AssetValuationService
     - [ ] 자산 가치 평가 로직 구현
     - [ ] 시장 데이터 연동 구조 설계
     - [ ] 가치 변동 추적 구현
   - [ ] TransactionValidationService
     - [ ] 거래 유효성 검증 규칙 구현
     - [ ] 예산 제한 검증
     - [ ] 이상 거래 감지

## 우선순위 중간 (P1)
1. [ ] 인프라스트럭처 레이어 구현
   - [ ] MongoDB 연결 설정
   - [ ] MongoDB 레포지토리 구현
   - [ ] 트랜잭션 관리자 구현
   - [ ] 에러 핸들링 구현

2. [ ] 애플리케이션 레이어 구현
   - [ ] UseCase 인터페이스 정의
   - [ ] Asset 관련 UseCase 구현
     - [ ] CreateAssetUseCase
     - [ ] UpdateAssetUseCase
     - [ ] DeleteAssetUseCase
   - [ ] Transaction 관련 UseCase 구현
     - [ ] RecordTransactionUseCase
     - [ ] UpdateTransactionUseCase
   - [ ] Portfolio 관련 UseCase 구현
     - [ ] CreatePortfolioUseCase
     - [ ] UpdatePortfolioUseCase
     - [ ] RebalancePortfolioUseCase

3. [ ] API 레이어 구현
   - [ ] Echo 프레임워크 설정
   - [ ] 미들웨어 구현
     - [ ] 인증 미들웨어
     - [ ] 로깅 미들웨어
     - [ ] 에러 핸들링 미들웨어
   - [ ] API 엔드포인트 구현
     - [ ] Asset API
     - [ ] Transaction API
     - [ ] Portfolio API

## 우선순위 낮음 (P2)
1. [ ] 테스트 자동화
   - [ ] 테스트 헬퍼 구현
   - [ ] 통합 테스트 환경 구성
   - [ ] 성능 테스트 시나리오 작성

2. [ ] 문서화
   - [ ] API 문서 자동화 (Swagger)
   - [ ] 아키텍처 결정 기록 (ADR)
   - [ ] 개발 가이드 작성

3. [ ] 모니터링 시스템
   - [ ] 메트릭 수집 설정
   - [ ] 로깅 시스템 구축
   - [ ] 알림 시스템 구현

## 참고사항
- 각 작업은 LLM_ROLE.md의 가이드라인을 따라 진행
- 완료된 작업은 DONE.md로 이동
- 각 작업은 가능한 한 작은 단위로 분할하여 관리
- 모든 코드 변경은 lint 검사 필수 