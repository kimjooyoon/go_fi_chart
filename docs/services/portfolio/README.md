# Portfolio 서비스

## 개요
Portfolio 서비스는 사용자의 투자 포트폴리오를 관리하는 마이크로서비스입니다.

## 주요 기능
- 포트폴리오 생성, 조회, 수정, 삭제
- 자산 배분 관리
- 포트폴리오 성과 분석

## 기술 스택
- Go 1.24.0
- Chi 라우터
- In-memory 저장소 (현재)

## API 엔드포인트
### 포트폴리오 관리
- `GET /api/v1/portfolios` - 포트폴리오 목록 조회
- `POST /api/v1/portfolios` - 새 포트폴리오 생성
- `GET /api/v1/portfolios/{id}` - 특정 포트폴리오 조회
- `PUT /api/v1/portfolios/{id}` - 포트폴리오 정보 수정
- `DELETE /api/v1/portfolios/{id}` - 포트폴리오 삭제
- `POST /api/v1/portfolios/{id}/assets` - 자산 추가
- `PUT /api/v1/portfolios/{id}/assets/{assetId}` - 자산 비중 수정
- `DELETE /api/v1/portfolios/{id}/assets/{assetId}` - 자산 제거

## 의존성
- Asset 서비스: 포트폴리오 구성 자산 정보 조회
- Transaction 서비스: 포트폴리오 거래 내역 연동
- Monitoring 서비스: 포트폴리오 상태 모니터링

## 설정
### 환경 변수
- `PORT`: 서비스 포트 (기본값: 8081)

## 로컬 개발 환경 설정
```bash
# 서비스 실행
go run cmd/server/main.go

# 테스트 실행
go test ./...
```

## 도메인 모델
- Portfolio: 포트폴리오 정보
- PortfolioAsset: 포트폴리오 구성 자산
- Percentage: 자산 비중 값 객체 