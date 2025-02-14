package asset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_should_generate_unique_id_when_called_multiple_times(t *testing.T) {
	// Given
	iterations := 1000
	idMap := make(map[string]bool)

	// When
	for i := 0; i < iterations; i++ {
		id := generateID()

		// Then
		// 이미 생성된 ID가 없어야 함
		assert.False(t, idMap[id], "ID should be unique")
		idMap[id] = true

		// ID 형식이 올바른지 검증
		assert.Len(t, id, 36, "ID should be a valid UUID string")
		assert.Regexp(t, "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$", id,
			"ID should be a valid UUID v4 string")
	}
}

func Test_should_generate_valid_uuid_v4_when_called(t *testing.T) {
	// When
	id := generateID()

	// Then
	assert.Len(t, id, 36, "ID should be a valid UUID string")
	assert.Regexp(t, "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$", id,
		"ID should be a valid UUID v4 string")
}

func Test_should_not_return_empty_string_when_called(t *testing.T) {
	// When
	id := generateID()

	// Then
	assert.NotEmpty(t, id, "ID should not be empty")
}

func Test_should_be_thread_safe_when_called_concurrently(t *testing.T) {
	// Given
	iterations := 1000
	idMap := make(map[string]bool)
	idChan := make(chan string, iterations)
	done := make(chan bool)

	// When
	for i := 0; i < iterations; i++ {
		go func() {
			id := generateID()
			idChan <- id
		}()
	}

	// 모든 ID를 수집
	go func() {
		for i := 0; i < iterations; i++ {
			id := <-idChan
			// Then
			assert.False(t, idMap[id], "ID should be unique even in concurrent execution")
			idMap[id] = true
		}
		done <- true
	}()

	<-done
}
