# Transaction 서비스

## 개요
Transaction 서비스는 자산 거래 내역을 관리하는 마이크로서비스입니다.

## 주요 기능
- 거래 내역 생성, 조회, 수정, 삭제
- 자산별 거래 내역 관리
- 포트폴리오별 거래 내역 관리

## 기술 스택
- Go 1.24.0
- Chi 라우터
- In-memory 저장소 (현재)

## API 엔드포인트
### 거래 내역 관리
- `GET /api/v1/transactions` - 거래 내역 목록 조회
- `POST /api/v1/transactions` - 새 거래 내역 생성
- `GET /api/v1/transactions/{id}` - 특정 거래 내역 조회
- `PUT /api/v1/transactions/{id}` - 거래 내역 수정
- `DELETE /api/v1/transactions/{id}` - 거래 내역 삭제
- `GET /api/v1/transactions/user/{userID}` - 사용자별 거래 내역 조회
- `GET /api/v1/transactions/portfolio/{portfolioID}` - 포트폴리오별 거래 내역 조회
- `GET /api/v1/transactions/asset/{assetID}` - 자산별 거래 내역 조회

## 의존성
- Asset 서비스: 거래 대상 자산 정보 조회
- Portfolio 서비스: 포트폴리오 정보 연동
- Monitoring 서비스: 거래 활동 모니터링

## 설정
### 환경 변수
- `PORT`: 서비스 포트 (기본값: 8082)

## 로컬 개발 환경 설정
```bash
# 서비스 실행
go run cmd/server/main.go

# 테스트 실행
go test ./...
```

## 도메인 모델
- Transaction: 거래 내역 정보
- Money: 거래 금액 값 객체
- TransactionType: 거래 유형 (매수/매도) 