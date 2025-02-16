package events

import (
	"context"
	"fmt"
	"time"
)

// RetryPolicy는 이벤트 처리 재시도 정책을 정의합니다.
type RetryPolicy interface {
	// ShouldRetry는 주어진 에러에 대해 재시도를 해야 하는지 결정합니다.
	ShouldRetry(err error) bool

	// NextBackoff는 다음 재시도까지의 대기 시간을 반환합니다.
	NextBackoff(attempt int) time.Duration

	// MaxAttempts는 최대 재시도 횟수를 반환합니다.
	MaxAttempts() int
}

// ExponentialBackoff는 지수 백오프 재시도 정책을 구현합니다.
type ExponentialBackoff struct {
	initialInterval time.Duration
	maxInterval     time.Duration
	maxAttempts     int
	multiplier      float64
}

// NewExponentialBackoff는 새로운 ExponentialBackoff를 생성합니다.
func NewExponentialBackoff(
	initialInterval time.Duration,
	maxInterval time.Duration,
	maxAttempts int,
	multiplier float64,
) *ExponentialBackoff {
	return &ExponentialBackoff{
		initialInterval: initialInterval,
		maxInterval:     maxInterval,
		maxAttempts:     maxAttempts,
		multiplier:      multiplier,
	}
}

// ShouldRetry는 주어진 에러에 대해 재시도를 해야 하는지 결정합니다.
func (e *ExponentialBackoff) ShouldRetry(err error) bool {
	return true
}

// NextBackoff는 다음 재시도까지의 대기 시간을 반환합니다.
func (e *ExponentialBackoff) NextBackoff(attempt int) time.Duration {
	if attempt <= 0 {
		return e.initialInterval
	}

	backoff := float64(e.initialInterval) * pow(e.multiplier, float64(attempt-1))
	if backoff > float64(e.maxInterval) {
		return e.maxInterval
	}
	return time.Duration(backoff)
}

// MaxAttempts는 최대 재시도 횟수를 반환합니다.
func (e *ExponentialBackoff) MaxAttempts() int {
	return e.maxAttempts
}

// pow는 거듭제곱을 계산합니다.
func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}

// RetryableEventHandler는 재시도 정책을 적용한 이벤트 핸들러입니다.
type RetryableEventHandler struct {
	handler EventHandler
	policy  RetryPolicy
}

// NewRetryableEventHandler는 새로운 RetryableEventHandler를 생성합니다.
func NewRetryableEventHandler(handler EventHandler, policy RetryPolicy) *RetryableEventHandler {
	return &RetryableEventHandler{
		handler: handler,
		policy:  policy,
	}
}

// HandleEvent는 재시도 정책을 적용하여 이벤트를 처리합니다.
func (h *RetryableEventHandler) HandleEvent(ctx context.Context, event Event) error {
	var lastErr error
	attempt := 0

	for attempt < h.policy.MaxAttempts() {
		err := h.handler.HandleEvent(ctx, event)
		if err == nil {
			return nil
		}

		lastErr = err
		if !h.policy.ShouldRetry(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		attempt++
		if attempt >= h.policy.MaxAttempts() {
			break
		}

		backoff := h.policy.NextBackoff(attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			continue
		}
	}

	return fmt.Errorf("max retry attempts (%d) reached: %w", h.policy.MaxAttempts(), lastErr)
}

// HandlerType는 원본 핸들러의 타입을 반환합니다.
func (h *RetryableEventHandler) HandlerType() string {
	return h.handler.HandlerType()
}

// DefaultRetryPolicy는 기본 재시도 정책을 제공합니다.
func DefaultRetryPolicy() RetryPolicy {
	return NewExponentialBackoff(
		100*time.Millisecond, // 초기 대기 시간
		10*time.Second,       // 최대 대기 시간
		5,                    // 최대 재시도 횟수
		2.0,                  // 배수
	)
}
