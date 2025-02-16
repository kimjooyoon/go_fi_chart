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

.PHONY: all
all: init lint test security coverage

# 초기 설정
.PHONY: init
init: install-tools tidy
	@echo "초기화 완료"

.PHONY: install-tools
install-tools:
	@echo "개발 도구 설치 중..."
	@if ! which golangci-lint >/dev/null; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION); \
	fi
	@if ! which gosec >/dev/null; then \
		go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION); \
	fi

.PHONY: tidy
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

# 테스트
.PHONY: test
test: test-all test-services

.PHONY: test-all
test-all:
	@echo "전체 테스트 실행 중..."
	$(GO) test -v -race ./...

.PHONY: test-services
test-services:
	@echo "서비스별 테스트 실행 중..."
	@for service in $(SERVICES); do \
		echo "테스트: $$service"; \
		$(GO) test -v -race ./services/$$service/...; \
	done

# 린트
.PHONY: lint
lint: lint-all lint-services

.PHONY: lint-all
lint-all:
	@echo "전체 린트 검사 중..."
	golangci-lint run ./...

.PHONY: lint-services
lint-services:
	@echo "서비스별 린트 검사 중..."
	@for service in $(SERVICES); do \
		echo "린트: $$service"; \
		golangci-lint run ./services/$$service/...; \
	done

# 보안 검사
.PHONY: security
security: security-all security-services

.PHONY: security-all
security-all:
	@echo "전체 보안 검사 중..."
	gosec ./...

.PHONY: security-services
security-services:
	@echo "서비스별 보안 검사 중..."
	@for service in $(SERVICES); do \
		echo "보안 검사: $$service"; \
		gosec ./services/$$service/...; \
	done

# 커버리지
.PHONY: coverage
coverage: coverage-html coverage-services

.PHONY: coverage-html
coverage-html:
	@echo "전체 커버리지 리포트 생성 중..."
	$(GO) test -coverprofile=$(COVERAGE_FILE) ./...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

.PHONY: coverage-services
coverage-services:
	@echo "서비스별 커버리지 리포트 생성 중..."
	@for service in $(SERVICES); do \
		echo "커버리지: $$service"; \
		$(GO) test -coverprofile=services/$$service/coverage.out ./services/$$service/...; \
		$(GO) tool cover -html=services/$$service/coverage.out -o services/$$service/coverage.html; \
	done

# 실행
.PHONY: run
run:
	@echo "서버 실행 중..."
	$(GO) run $(MAIN_PACKAGE)

# 정리
.PHONY: clean
clean:
	@echo "임시 파일 정리 중..."
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@for service in $(SERVICES); do \
		rm -f services/$$service/coverage.out services/$$service/coverage.html; \
	done

# 도움말
.PHONY: help
help:
	@echo "사용 가능한 명령어:"
	@echo "  make init          - 개발 환경 초기화 (도구 설치 및 의존성 정리)"
	@echo "  make test          - 모든 테스트 실행"
	@echo "  make lint          - 코드 린트 검사"
	@echo "  make security      - 보안 취약점 검사"
	@echo "  make coverage      - 테스트 커버리지 리포트 생성"
	@echo "  make run           - 서버 실행"
	@echo "  make clean         - 임시 파일 정리"
	@echo "  make all           - 전체 검사 실행 (init, lint, test, security, coverage)" 