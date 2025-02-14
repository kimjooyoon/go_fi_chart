.PHONY: init test lint build run clean mock docs ci

# 기본 포맷팅
fmt:
	go fmt ./...

# 초기 개발 환경 설정
init:
	go mod tidy
	go install github.com/golangci-lint/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang/mock/mockgen@latest
	go install golang.org/x/tools/cmd/godoc@latest

# 테스트 실행 (커버리지 리포트 생성)
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 코드 포맷팅 및 린트 검사
lint:
	go fmt ./...
	golangci-lint run ./...

# 프로젝트 빌드
build:
	go build -v -o bin/server ./cmd/server

# 서버 실행
run: build
	./bin/server

# 생성된 파일들 정리
clean:
	rm -rf bin coverage.out coverage.html

# 모든 인터페이스에 대한 mock 생성
mock:
	mockgen -source=internal/domain/repository.go -destination=internal/domain/mock/repository_mock.go
	mockgen -source=internal/domain/asset/repository.go -destination=internal/domain/asset/mock/repository_mock.go

# API 문서 생성 및 실행
docs:
	godoc -http=:6060

# 전체 CI 파이프라인
ci: lint test build 