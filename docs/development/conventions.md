# 개발 컨벤션

## 코드 스타일

### 1. Go 코드 스타일
- [Effective Go](https://golang.org/doc/effective_go) 가이드라인 준수
- `gofmt` 사용 필수
- `golangci-lint` 규칙 준수

### 2. 패키지 구조
```
service/
├── cmd/                    # 실행 파일
│   └── server/
│       └── main.go
├── internal/              # 내부 패키지
│   ├── api/              # API 핸들러
│   ├── domain/           # 도메인 모델
│   └── infrastructure/   # 인프라 구현
└── pkg/                  # 공개 패키지
```

### 3. 네이밍 컨벤션
- 파일명: 스네이크 케이스 (`user_repository.go`)
- 패키지명: 소문자, 단일 단어
- 인터페이스: 동사+er (`Reader`, `Writer`)
- 메서드/함수: 카멜 케이스 (`GetUserByID`)
- 상수: 대문자 스네이크 케이스 (`MAX_RETRY_COUNT`)

### 4. 주석 규칙
- 모든 공개 API에 주석 필수
- 한글 주석 사용
- 패키지 설명 필수
```go
// Package user는 사용자 관련 기능을 제공합니다.
package user

// User는 사용자 정보를 나타냅니다.
type User struct {
// ID는 사용자의 고유 식별자입니다.
ID string
}
```

## 테스트

### 1. 테스트 규칙
- 모든 공개 함수에 대한 테스트 작성
- 테이블 기반 테스트 사용
- 테스트 커버리지 80% 이상 유지

### 2. 테스트 네이밍
```go
func TestUserRepository_GetUser_should_return_user_when_exists(t *testing.T)
func TestUserRepository_GetUser_should_return_error_when_not_found(t *testing.T)
```

### 3. 테스트 구조
```go
// Given
// 테스트 데이터 및 조건 설정

// When
// 테스트할 기능 실행

// Then
// 결과 검증
```

## 에러 처리

### 1. 에러 타입
```go
// 도메인 에러
type DomainError struct {
Code    string
Message string
}

// 인프라 에러
type InfrastructureError struct {
Code    string
Message string
Cause   error
}
```

### 2. 에러 메시지
- 명확하고 구체적인 메시지
- 한글 메시지 사용
- 해결 방법 포함 (가능한 경우)

## 로깅

### 1. 로그 레벨
- DEBUG: 개발 시 상세 정보
- INFO: 일반적인 작업 정보
- WARN: 잠재적 문제
- ERROR: 처리된 에러
- FATAL: 복구 불가능한 에러

### 2. 로그 포맷
```json
{
"level": "INFO",
"timestamp": "2024-02-16T12:00:00Z",
"service": "asset",
"trace_id": "abc123",
"message": "자산 생성 완료",
"data": {
"asset_id": "123",
"user_id": "456"
}
}
```

## 버전 관리

### 1. 브랜치 전략
- main: 프로덕션 코드
- develop: 개발 코드
- feature/*: 기능 개발
- bugfix/*: 버그 수정
- release/*: 릴리스 준비

### 2. 커밋 메시지
```
feat: 새로운 기능 추가
fix: 버그 수정
docs: 문서 수정
style: 코드 포맷팅
refactor: 코드 리팩토링
test: 테스트 코드
chore: 빌드 프로세스 변경
```

### 3. Pull Request
- 제목: `[타입] 작업 내용 요약`
- 본문: 작업 내용 상세 설명
- 리뷰어: 최소 1명 이상
- 테스트 결과 포함 