package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Run("NewError should create error with valid data", func(t *testing.T) {
		err := NewError("test", ErrCodeNotFound, "not found")
		assert.NotNil(t, err)
		assert.Equal(t, "test", err.Domain())
		assert.Equal(t, ErrCodeNotFound, err.Code())
		assert.Equal(t, "[test] NOT_FOUND: not found", err.Error())
	})

	t.Run("BaseError should implement Error interface", func(_ *testing.T) {
		var _ Error = &BaseError{}
	})

	t.Run("Error codes should be defined", func(t *testing.T) {
		assert.Equal(t, "NOT_FOUND", ErrCodeNotFound)
		assert.Equal(t, "ALREADY_EXISTS", ErrCodeAlreadyExists)
		assert.Equal(t, "INVALID_ARGUMENT", ErrCodeInvalidArgument)
		assert.Equal(t, "INVALID_OPERATION", ErrCodeInvalidOperation)
		assert.Equal(t, "NOT_IMPLEMENTED", ErrCodeNotImplemented)
		assert.Equal(t, "INTERNAL", ErrCodeInternal)
	})
}
