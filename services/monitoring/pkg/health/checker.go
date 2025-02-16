package health

import (
	"context"
	"sync"
	"time"
)

// Status 서비스의 상태를 나타냅니다.
type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

// Check 헬스 체크 결과를 나타냅니다.
type Check struct {
	Status    Status    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

// Checker 서비스 상태를 체크하는 인터페이스입니다.
type Checker interface {
	// Check 현재 서비스의 상태를 확인합니다.
	Check(ctx context.Context) Check
}

// SimpleChecker 기본적인 상태 체크 구현체입니다.
type SimpleChecker struct {
	mu     sync.RWMutex
	status Status
	err    error
}

// NewSimpleChecker 새로운 SimpleChecker를 생성합니다.
func NewSimpleChecker() *SimpleChecker {
	return &SimpleChecker{
		status: StatusUp,
	}
}

// Check 현재 서비스의 상태를 반환합니다.
func (c *SimpleChecker) Check(_ context.Context) Check {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var errStr string
	if c.err != nil {
		errStr = c.err.Error()
	}

	return Check{
		Status:    c.status,
		Timestamp: time.Now(),
		Error:     errStr,
	}
}

// SetStatus 서비스의 상태를 설정합니다.
func (c *SimpleChecker) SetStatus(status Status, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.status = status
	c.err = err
}
