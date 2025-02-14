# FIN-RPG: 게임화 요소를 활용한 개인 자산 관리 시스템

[![CI](https://github.com/aske/go_fi_chart/actions/workflows/ci.yml/badge.svg)](https://github.com/aske/go_fi_chart/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/aske/go_fi_chart/branch/main/graph/badge.svg)](https://codecov.io/gh/aske/go_fi_chart)
[![Go Report Card](https://goreportcard.com/badge/github.com/aske/go_fi_chart)](https://goreportcard.com/report/github.com/aske/go_fi_chart)

FIN-RPG는 개인 자산 관리의 진입 장벽을 낮추고 지속적인 관리 동기를 부여하기 위해 게임화(Gamification) 요소를 도입한 시스템입니다. 복잡한 재무 관리를 재미있고 이해하기 쉽게 만들어, 사용자가 자연스럽게 올바른 자산 관리 습관을 기를 수 있도록 돕습니다.

## 핵심 기능

### 자산 관리 (Core)
- 수입/지출 추적 및 분석
  - 상세한 거래 내역 관리
  - 카테고리별 지출 분석
  - 예산 대비 실제 지출 추적
- 자산 포트폴리오 관리
  - 자산 유형별 현황 파악
  - 포트폴리오 다각화 지표
  - 수익률 분석 및 리밸런싱 알림
- 재무 목표 설정 및 추적
  - 단기/중기/장기 목표 관리
  - 목표 달성률 모니터링
  - 맞춤형 조언 제공

### 의사결정 지원 (Support)
- 재무 상태 분석
  - 핵심 재무 지표 계산
  - 연령대별 평균 비교
  - 개선 포인트 제안
- 투자 전략 지원
  - 위험 성향 분석
  - 맞춤형 포트폴리오 제안
  - 시장 상황 알림
- 지출 패턴 분석
  - AI 기반 소비 패턴 분석
  - 최적화 제안
  - 이상 지출 감지

### 동기 부여 시스템 (Gamification)
- 성장 트래킹
  - 재무 건전성 점수
  - 저축률 달성 보상
  - 투자 다각화 점수
- 목표 달성 보상
  - 단계별 목표 설정
  - 달성 시 뱃지 획득
  - 특별 콘텐츠 해금
- 커뮤니티 참여
  - 노하우 공유 및 평가
  - 성공 사례 분석
  - 멘토링 시스템

## 기술 스택

- **백엔드**: Go
  - Clean Architecture
  - Domain-Driven Design
  - Event-Driven Architecture
- **프론트엔드**: Flutter
  - 반응형 디자인
  - 실시간 데이터 시각화
  - 크로스 플랫폼 지원
- **데이터베이스**: MongoDB
  - 유연한 데이터 모델
  - 실시간 집계 기능
  - 확장성 있는 구조

## 설치 방법

1. 이 저장소를 클론합니다.
   ```bash
   git clone https://github.com/username/repo-name.git
   ```
2. 필요한 패키지를 설치합니다.
   ```bash
   cd repo-name
   # Flutter 설치
   flutter pub get
   ```
3. 서버를 시작합니다.
   ```bash
   # Go 서버 실행
   go run main.go
   ```

## 게임 시작하기

1. 계정을 생성하고 초기 스킬트리를 선택합니다.
2. 일일 퀘스트를 확인하고 수행합니다.
3. 자산을 관리하며 경험치를 획득합니다.
4. 수익 코인을 모아 새로운 전략을 해금합니다.
5. 커뮤니티와 상호작용하며 랭킹에 도전합니다.

## 프로젝트 구조

```
repo-name/
├── frontend/       # Flutter 프론트엔드 코드
├── backend/        # Go 백엔드 코드
├── database/       # MongoDB 데이터베이스 스키마
└── docs/          # 프로젝트 문서
    ├── LLM_ROLE.md  # 개발 프로세스 및 규칙
    ├── TODO.md      # 진행할 작업 목록
    ├── DONE.md      # 완료된 작업 목록
    └── PROBLEMS.md  # 프로젝트 진행 중 발생하는 문제점을 관리합니다.
```

## 문서 관리 시스템

이 프로젝트는 자동화된 문서 관리 시스템을 통해 개발 프로세스를 추적하고 관리합니다:

1. **LLM_ROLE.md**: 개발 프로세스와 규칙을 정의합니다.
   - DDD, TDD 등의 개발 방법론 가이드
   - 코드 품질 기준
   - 테스트 전략

2. **TODO.md**: 진행할 작업 목록을 관리합니다.
   - 우선순위별 작업 분류 (P0, P1, P2)
   - 세부 작업 항목
   - 작업 의존성 관리

3. **DONE.md**: 완료된 작업을 기록합니다.
   - 날짜별 작업 완료 기록
   - 주요 결정사항 문서화
   - 변경사항 추적

4. **PROBLEMS.md**: 프로젝트 진행 중 발생하는 문제점을 관리합니다.
   - 문제점 발견 및 추적
   - 심각도 및 영향도 평가
   - 해결 방안 제시
   - 해결 과정 문서화

## 기여 방법

기여를 원하시는 분은 다음 단계를 따라 주세요:

1. 이 저장소를 포크합니다.
2. 새로운 브랜치를 생성합니다.
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. 변경 사항을 커밋합니다.
   ```bash
   git commit -m "설명 메시지"
   ```
4. 브랜치를 푸시합니다.
   ```bash
   git push origin feature/your-feature-name
   ```
5. 풀 리퀘스트를 제출합니다.

## 연락처

프로젝트에 대한 질문이나 피드백이 있으시면 아래의 연락처로 문의해 주세요:

- 이메일: your-email@example.com
- GitHub: [your-github-username](https://github.com/your-github-username)

## 라이센스

이 프로젝트는 MIT 라이센스 하에 배포됩니다. 