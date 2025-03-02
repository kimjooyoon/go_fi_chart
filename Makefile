# Go 관련 변수
GO=go
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin

# 도구 버전
GOLANGCI_LINT_VERSION=v1.55.2
GOSEC_VERSION=v2.18.2

# 프로젝트 설정
PROJECT_NAME=go_fi_chart
MAIN_PACKAGE=./cmd/server
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# 서비스 목록
SERVICES=asset portfolio transaction monitoring

.PHONY: all tidy verify tools lint test test-services sec build clean all-with-service-tests

# 기본 all 타겟에서는 test-services를 제외합니다
all: tools tidy verify lint test sec build

# 모든 테스트(전체 + 서비스별)를 실행하는 새로운 타겟
all-with-service-tests: tools tidy verify lint test test-services sec build

tools:
	@echo "개발 도구 설치 중..."
	@if ! which golangci-lint >/dev/null; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION); \
	fi
	@if ! which gosec >/dev/null; then \
		go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION); \
	fi
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest

tidy:
	@echo "의존성 정리 중..."
	$(GO) mod tidy
	@echo "pkg 의존성 정리 중..."
	(cd pkg && $(GO) mod tidy)
	@echo "서비스 의존성 정리 중..."
	@for service in $(SERVICES); do \
		echo "의존성 정리: $$service"; \
		(cd services/$$service && $(GO) mod tidy); \
	done
	@echo "go work sync 실행..."
	$(GO) work sync
	@go mod verify

verify:
	@go mod verify
	@go mod why -m all

lint:
	@echo "전체 린트 검사 중..."
	golangci-lint run ./...
	@goimports -w .
	@staticcheck ./...
	@go vet ./...

# 메인 테스트 타겟을 수정하여 전체 테스트만 실행합니다
test:
	@echo "전체 테스트 실행 중..."
	$(GO) test -v -race -coverprofile=$(COVERAGE_FILE) ./...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

# 서비스별 테스트를 별도의 타겟으로 분리합니다
test-services:
	@echo "서비스별 테스트 실행 중..."
	@for service in $(SERVICES); do \
		echo "테스트 실행: $$service"; \
		(cd services/$$service && \
		$(GO) test -v -race -coverprofile=coverage.out ./... && \
		$(GO) tool cover -html=coverage.out -o coverage.html); \
	done

sec:
	@echo "전체 보안 검사 중..."
	gosec ./...

build:
	$(GO) build -v ./...

clean:
	@echo "임시 파일 정리 중..."
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@for service in $(SERVICES); do \
		rm -f services/$$service/coverage.out services/$$service/coverage.html; \
	done
	go clean
	rm -f $(shell find . -name '*.test')
	rm -f $(shell find . -name '*.out')

# 도움말
.PHONY: help
help:
	@echo "사용 가능한 명령어:"
	@echo "  make init          - 개발 환경 초기화 (도구 설치 및 의존성 정리)"
	@echo "  make test          - 모든 테스트 실행"
	@echo "  make test-services - 서비스별 테스트 실행"
	@echo "  make lint          - 코드 린트 검사"
	@echo "  make security      - 보안 취약점 검사"
	@echo "  make coverage      - 테스트 커버리지 리포트 생성"
	@echo "  make run           - 서버 실행"
	@echo "  make clean         - 임시 파일 정리"
	@echo "  make all           - 전체 검사 실행 (init, lint, test, security, coverage)"
	@echo "  make all-with-service-tests - 전체 검사 + 서비스별 테스트 실행" 