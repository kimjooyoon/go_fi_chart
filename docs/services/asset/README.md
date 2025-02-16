# Asset 서비스

## 개요
Asset 서비스는 사용자의 자산을 관리하는 마이크로서비스입니다.

## 주요 기능
- 자산 생성, 조회, 수정, 삭제
- 자산 가치 평가
- 자산 성과 추적

## 기술 스택
- Go 1.24.0
- Chi 라우터
- In-memory 저장소 (현재)

## API 엔드포인트
### 자산 관리
- `GET /api/v1/assets` - 자산 목록 조회
- `POST /api/v1/assets` - 새 자산 생성
- `GET /api/v1/assets/{id}` - 특정 자산 조회
- `PUT /api/v1/assets/{id}` - 자산 정보 수정
- `DELETE /api/v1/assets/{id}` - 자산 삭제

## 의존성
- Portfolio 서비스: 포트폴리오 구성을 위한 자산 정보 제공
- Transaction 서비스: 자산 거래 정보 연동
- Monitoring 서비스: 자산 상태 모니터링

## 설정
### 환경 변수
- `PORT`: 서비스 포트 (기본값: 8080)

## 로컬 개발 환경 설정
```bash
# 서비스 실행
go run cmd/server/main.go

# 테스트 실행
go test ./...
```

## 도메인 모델
- Asset: 자산 정보
- Money: 화폐 값 객체
- Performance: 자산 성과 정보 