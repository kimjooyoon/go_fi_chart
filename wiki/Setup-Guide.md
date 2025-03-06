# 환경 설정 가이드

이 문서는 Go Finance Chart 프로젝트의 개발 환경 설정 방법을 안내합니다.

## 필수 요구사항

### 소프트웨어 요구사항
- Go 1.18 이상 (제네릭 지원을 위해)
- Git
- 코드 에디터 (VS Code 권장)
- 웹 브라우저 (Chrome, Firefox 권장)

### 선택적 도구
- Docker (컨테이너화된 개발 환경을 위한)
- Make (빌드 자동화)
- Air (실시간 리로드)

## 개발 환경 구축

### 1. Go 설치

#### macOS (Homebrew 사용)
```bash
brew install go
```

#### Linux
```bash
# Debian/Ubuntu
sudo apt update
sudo apt install golang-go

# RHEL/CentOS/Fedora
sudo dnf install golang
```

#### Windows
- [Go 공식 웹사이트](https://golang.org/dl/)에서 인스톨러 다운로드
- 인스톨러 실행 및 화면의 지시 따르기

설치 후 확인:
```bash
go version
```

### 2. 프로젝트 클론

```bash
git clone https://github.com/kimjooyoon/go_fi_chart.git
cd go_fi_chart
```

### 3. 의존성 설치

```bash
go mod download
```

또는 처음부터 모듈 초기화:

```bash
go mod init github.com/kimjooyoon/go_fi_chart
go mod tidy
```

### 4. 환경 변수 설정

개발 환경에 따라 다음 환경 변수를 설정해야 할 수 있습니다:

```bash
# 예시
export API_KEY_YAHOO_FINANCE="your_api_key_here"
export API_KEY_NEWS_SERVICE="your_api_key_here"
export DEV_MODE=true
```

Windows PowerShell에서는:
```powershell
$env:API_KEY_YAHOO_FINANCE="your_api_key_here"
$env:API_KEY_NEWS_SERVICE="your_api_key_here"
$env:DEV_MODE="true"
```

프로젝트 루트에 `.env` 파일을 생성하여 환경 변수를 관리할 수도 있습니다:
```
API_KEY_YAHOO_FINANCE=your_api_key_here
API_KEY_NEWS_SERVICE=your_api_key_here
DEV_MODE=true
```

### 5. 빌드 및 실행

#### 빌드
```bash
go build -o go_fi_chart ./cmd/main.go
```

#### 실행
```bash
./go_fi_chart
```

Windows에서는:
```bash
go_fi_chart.exe
```

### 6. 개발 모드로 실행 (실시간 리로드)

Air를 사용하는 경우:
```bash
# Air 설치 (처음 한 번만)
go install github.com/cosmtrek/air@latest

# Air 실행
air
```

## 테스트 실행

### 전체 테스트 실행
```bash
go test ./...
```

### 특정 패키지 테스트
```bash
go test ./internal/repository
```

### 커버리지 리포트 생성
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 코드 스타일 및 린트

### Go 코드 스타일 체크
```bash
# golangci-lint 설치 (처음 한 번만)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 린트 실행
golangci-lint run
```

### 코드 형식 자동 정리
```bash
go fmt ./...
```

## Docker로 개발 환경 구성 (선택적)

프로젝트 루트에 `Dockerfile`과 `docker-compose.yml`을 생성할 수 있습니다.

### Dockerfile 예시
```dockerfile
FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go_fi_chart ./cmd/main.go

EXPOSE 8080

CMD ["/go_fi_chart"]
```

### docker-compose.yml 예시
```yaml
version: '3'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - .:/app
    environment:
      - API_KEY_YAHOO_FINANCE=${API_KEY_YAHOO_FINANCE}
      - API_KEY_NEWS_SERVICE=${API_KEY_NEWS_SERVICE}
      - DEV_MODE=true
    command: go run cmd/main.go
```

### Docker Compose 실행
```bash
docker-compose up
```

## API 키 취득 방법

### Yahoo Finance API 키
1. [Yahoo Developer Network](https://developer.yahoo.com/) 접속
2. 계정 생성 및 로그인
3. 새로운 애플리케이션 등록
4. Finance API 사용 권한 요청
5. 발급받은 API 키 사용

### 뉴스 서비스 API 키 (선택적)
여러 뉴스 서비스 중 하나 또는 여러 개를 선택하여 사용할 수 있습니다:

1. [NewsAPI](https://newsapi.org/)
2. [Alpha Vantage](https://www.alphavantage.co/)
3. [Financial Times](https://developer.ft.com/)

각 서비스의 개발자 포털에서 API 키를 발급받을 수 있습니다.

## 트러블슈팅

### 일반적인 문제

#### "package XXX is not in GOROOT" 오류
```bash
go mod tidy
```

#### 환경 변수 관련 오류
```bash
# 환경 변수가 제대로 설정되었는지 확인
echo $API_KEY_YAHOO_FINANCE
```

#### 실행 시 접속 권한 오류
```bash
# 실행 파일에 실행 권한 부여
chmod +x go_fi_chart
```

### 도움 받기

더 많은 도움이 필요하시면:
1. 프로젝트 GitHub 이슈 섹션 확인
2. 새로운 이슈 생성

## 다음 단계

환경 설정을 완료하셨다면, 다음 단계로 [개발 로드맵](Development-Roadmap)을 확인하여 프로젝트의 현재 상태와 개발 계획을 파악하세요.