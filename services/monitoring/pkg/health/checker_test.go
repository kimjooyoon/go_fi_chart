package health

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockCheck는 테스트용 헬스 체크입니다.
type MockCheck struct {
	name   string
	status Status
	err    error
}

func (m *MockCheck) Check(_ context.Context) (Result, error) {
	if m.err != nil {
		return Result{
			Name:      m.name,
			Status:    m.status,
			Error:     m.err.Error(),
			Timestamp: time.Now(),
		}, m.err
	}
	return Result{
		Name:      m.name,
		Status:    m.status,
		Timestamp: time.Now(),
	}, nil
}

func (m *MockCheck) Name() string {
	return m.name
}

func Test_NewChecker_should_create_checker(t *testing.T) {
	// When
	checker := NewChecker(5 * time.Second)

	// Then
	assert.NotNil(t, checker)
	assert.Empty(t, checker.checks)
	assert.Equal(t, 5*time.Second, checker.interval)
}

func Test_Checker_should_add_and_check(t *testing.T) {
	// Given
	checker := NewChecker(5 * time.Second)
	mockCheck := &MockCheck{
		name:   "test-check",
		status: StatusUp,
	}

	// When
	checker.AddCheck(mockCheck)
	results := checker.CheckAll(context.Background())

	// Then
	assert.Len(t, results, 1)
	assert.Equal(t, StatusUp, results[0].Status)
	assert.Equal(t, "test-check", results[0].Name)
}

func Test_Checker_should_be_thread_safe(_ *testing.T) {
	// Given
	checker := NewChecker(5 * time.Second)
	iterations := 1000
	done := make(chan bool)

	// When
	go func() {
		for i := 0; i < iterations; i++ {
			checker.AddCheck(&MockCheck{
				name:   "test-check",
				status: StatusUp,
			})
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			checker.CheckAll(context.Background())
		}
		done <- true
	}()

	// Then
	<-done
	<-done
}

func TestChecker_CheckAll(t *testing.T) {
	// Given
	checker := NewChecker(5 * time.Second)
	mockCheck1 := &MockCheck{
		name:   "check1",
		status: StatusUp,
	}
	mockCheck2 := &MockCheck{
		name:   "check2",
		status: StatusDown,
		err:    errors.New("check failed"),
	}

	// When
	checker.AddCheck(mockCheck1)
	checker.AddCheck(mockCheck2)
	results := checker.CheckAll(context.Background())

	// Then
	assert.Len(t, results, 2)
	assert.False(t, results.IsHealthy())
}

func TestChecker_RemoveCheck(t *testing.T) {
	// Given
	checker := NewChecker(5 * time.Second)
	mockCheck := &MockCheck{
		name:   "test-check",
		status: StatusUp,
	}

	// When
	checker.AddCheck(mockCheck)
	checker.RemoveCheck("test-check")
	results := checker.CheckAll(context.Background())

	// Then
	assert.Empty(t, results)
}

func TestResults_IsHealthy(t *testing.T) {
	t.Run("모든 체크가 정상인 경우", func(t *testing.T) {
		results := Results{
			{
				Name:      "check1",
				Status:    StatusUp,
				Timestamp: time.Now(),
			},
			{
				Name:      "check2",
				Status:    StatusUp,
				Timestamp: time.Now(),
			},
		}
		assert.True(t, results.IsHealthy())
	})

	t.Run("하나라도 비정상인 경우", func(t *testing.T) {
		results := Results{
			{
				Name:      "check1",
				Status:    StatusUp,
				Timestamp: time.Now(),
			},
			{
				Name:      "check2",
				Status:    StatusDown,
				Timestamp: time.Now(),
			},
		}
		assert.False(t, results.IsHealthy())
	})
}

func Test_SimpleChecker_should_create_checker_with_up_status(t *testing.T) {
	// When
	checker := NewSimpleChecker()

	// Then
	result, err := checker.Check(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, StatusUp, result.Status)
	assert.Empty(t, result.Error)
	assert.NotZero(t, result.Timestamp)
}

func Test_SimpleChecker_should_update_status_and_error(t *testing.T) {
	// Given
	checker := NewSimpleChecker()
	expectedErr := errors.New("service unavailable")

	// When
	checker.SetStatus(StatusDown, expectedErr)

	// Then
	result, err := checker.Check(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, StatusDown, result.Status)
	assert.Equal(t, expectedErr.Error(), result.Error)
	assert.NotZero(t, result.Timestamp)
}

func Test_SimpleChecker_should_be_thread_safe(_ *testing.T) {
	// Given
	checker := NewSimpleChecker()
	iterations := 1000
	done := make(chan bool)

	// When
	go func() {
		for i := 0; i < iterations; i++ {
			checker.SetStatus(StatusDown, errors.New("error"))
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			checker.SetStatus(StatusUp, nil)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			_, _ = checker.Check(context.Background())
		}
		done <- true
	}()

	// Then
	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestCheck_Check(t *testing.T) {
	t.Run("성공적인 헬스 체크", func(t *testing.T) {
		// Given
		check := &MockCheck{
			name:   "test-check",
			status: StatusUp,
		}

		// When
		result, err := check.Check(context.Background())

		// Then
		assert.NoError(t, err)
		assert.Equal(t, StatusUp, result.Status)
		assert.NotZero(t, result.Timestamp)
	})

	t.Run("실패한 헬스 체크", func(t *testing.T) {
		// Given
		check := &MockCheck{
			name:   "test-check",
			status: StatusDown,
			err:    errors.New("check failed"),
		}

		// When
		result, err := check.Check(context.Background())

		// Then
		assert.Error(t, err)
		assert.Equal(t, StatusDown, result.Status)
		assert.NotZero(t, result.Timestamp)
	})
}

func TestChecker_AddCheck(t *testing.T) {
	// Given
	checker := NewChecker(time.Second)
	check := &MockCheck{
		name:   "test-check",
		status: StatusUp,
	}

	// When
	checker.AddCheck(check)

	// Then
	assert.Len(t, checker.checks, 1)
}

func TestChecker_Start(t *testing.T) {
	// Given
	checker := NewChecker(100 * time.Millisecond)
	check := &MockCheck{
		name:   "test-check",
		status: StatusUp,
	}
	checker.AddCheck(check)

	// When
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go checker.Start(ctx)

	// Then
	<-ctx.Done()
	results := checker.CheckAll(context.Background())
	assert.NotEmpty(t, results)
}

func TestChecker_Stop(t *testing.T) {
	// Given
	checker := NewChecker(100 * time.Millisecond)
	check := &MockCheck{
		name:   "test-check",
		status: StatusUp,
	}
	checker.AddCheck(check)

	// When
	ctx, cancel := context.WithCancel(context.Background())
	go checker.Start(ctx)
	time.Sleep(200 * time.Millisecond)
	cancel()

	// Then
	results := checker.CheckAll(context.Background())
	assert.NotEmpty(t, results)
}

func TestChecker_CheckWithError(t *testing.T) {
	// Given
	checker := NewChecker(100 * time.Millisecond)
	check := &MockCheck{
		name:   "test-check",
		status: StatusDown,
		err:    errors.New("check failed"),
	}
	checker.AddCheck(check)

	// When
	results := checker.CheckAll(context.Background())

	// Then
	assert.Len(t, results, 1)
	assert.Equal(t, StatusDown, results[0].Status)
}

func TestResult(t *testing.T) {
	t.Run("결과 생성", func(t *testing.T) {
		result := NewResult(true, nil)
		assert.Equal(t, StatusUp, result.Status)
		assert.Empty(t, result.Error)
	})

	t.Run("에러 추가", func(t *testing.T) {
		errors := map[string]error{
			"check1": errors.New("health check failed"),
			"check2": errors.New("health check timeout"),
		}
		result := NewResult(false, errors)
		assert.Equal(t, StatusDown, result.Status)
		assert.Contains(t, result.Error, "check1")
		assert.Contains(t, result.Error, "check2")
	})
}
