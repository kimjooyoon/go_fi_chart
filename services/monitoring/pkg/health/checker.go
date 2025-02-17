package health

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Status는 헬스 체크 상태를 나타냅니다.
type Status string

const (
	// StatusUp은 정상 상태를 나타냅니다.
	StatusUp Status = "UP"
	// StatusDown은 비정상 상태를 나타냅니다.
	StatusDown Status = "DOWN"
)

// Result는 헬스 체크 결과를 나타냅니다.
type Result struct {
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// NewResult는 새로운 Result를 생성합니다.
func NewResult(isHealthy bool, errors map[string]error) Result {
	status := StatusDown
	if isHealthy {
		status = StatusUp
	}

	var errStr string
	if len(errors) > 0 {
		for name, err := range errors {
			if errStr != "" {
				errStr += ", "
			}
			errStr += name + ": " + err.Error()
		}
	}

	return Result{
		Name:      "health",
		Status:    status,
		Error:     errStr,
		Timestamp: time.Now(),
	}
}

// Results는 여러 헬스 체크 결과를 나타냅니다.
type Results []Result

// IsHealthy는 모든 헬스 체크가 정상인지 확인합니다.
func (r Results) IsHealthy() bool {
	for _, result := range r {
		if result.Status != StatusUp {
			return false
		}
	}
	return true
}

// SimpleChecker는 간단한 헬스 체크 구현체입니다.
type SimpleChecker struct {
	mu     sync.RWMutex
	status Status
	err    error
}

// NewSimpleChecker는 새로운 SimpleChecker를 생성합니다.
func NewSimpleChecker() *SimpleChecker {
	return &SimpleChecker{
		status: StatusUp,
	}
}

// Check는 현재 상태를 반환합니다.
func (c *SimpleChecker) Check(_ context.Context) (Result, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var errStr string
	if c.err != nil {
		errStr = c.err.Error()
	}

	return Result{
		Name:      "simple",
		Status:    c.status,
		Error:     errStr,
		Timestamp: time.Now(),
	}, nil
}

// SetStatus는 체커의 상태를 설정합니다.
func (c *SimpleChecker) SetStatus(status Status, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status = status
	c.err = err
}

// Name은 체커의 이름을 반환합니다.
func (c *SimpleChecker) Name() string {
	return "simple"
}

// Check는 헬스 체크 인터페이스입니다.
type Check interface {
	// Check는 헬스 체크를 수행합니다.
	Check(ctx context.Context) (Result, error)
	// Name은 헬스 체크의 이름을 반환합니다.
	Name() string
}

// Checker는 여러 헬스 체크를 관리하는 구조체입니다.
type Checker struct {
	checks   map[string]Check
	interval time.Duration
	mu       sync.RWMutex
}

// NewChecker는 새로운 Checker를 생성합니다.
func NewChecker(interval time.Duration) *Checker {
	return &Checker{
		checks:   make(map[string]Check),
		interval: interval,
	}
}

// AddCheck는 헬스 체크를 추가합니다.
func (c *Checker) AddCheck(check Check) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checks[check.Name()] = check
}

// RemoveCheck는 헬스 체크를 제거합니다.
func (c *Checker) RemoveCheck(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.checks, name)
}

// CheckAll은 모든 헬스 체크를 수행합니다.
func (c *Checker) CheckAll(ctx context.Context) Results {
	c.mu.RLock()
	defer c.mu.RUnlock()

	results := make(Results, 0, len(c.checks))
	for _, check := range c.checks {
		result, err := check.Check(ctx)
		if err != nil {
			result = Result{
				Name:      check.Name(),
				Status:    StatusDown,
				Error:     err.Error(),
				Timestamp: time.Now(),
			}
		}
		results = append(results, result)
	}
	return results
}

// Start는 주기적인 헬스 체크를 시작합니다.
func (c *Checker) Start(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.CheckAll(ctx)
		}
	}
}

var (
	// ErrHealthCheckFailed 헬스 체크 실패 에러
	ErrHealthCheckFailed = errors.New("health check failed")
	// ErrHealthCheckTimeout 헬스 체크 타임아웃 에러
	ErrHealthCheckTimeout = errors.New("health check timeout")
)
