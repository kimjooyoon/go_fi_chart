package health

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSimpleChecker_should_create_checker_with_up_status(t *testing.T) {
	// When
	checker := NewSimpleChecker()

	// Then
	check := checker.Check(context.Background())
	assert.Equal(t, StatusUp, check.Status)
	assert.Empty(t, check.Error)
	assert.NotZero(t, check.Timestamp)
}

func Test_SimpleChecker_should_update_status_and_error(t *testing.T) {
	// Given
	checker := NewSimpleChecker()
	expectedErr := errors.New("service unavailable")

	// When
	checker.SetStatus(StatusDown, expectedErr)

	// Then
	check := checker.Check(context.Background())
	assert.Equal(t, StatusDown, check.Status)
	assert.Equal(t, expectedErr.Error(), check.Error)
	assert.NotZero(t, check.Timestamp)
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
			_ = checker.Check(context.Background())
		}
		done <- true
	}()

	// Then
	for i := 0; i < 3; i++ {
		<-done
	}
}
