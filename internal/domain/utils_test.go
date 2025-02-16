package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUtils(_ *testing.T) {
	// TODO: Add tests
}

func TestGenerateID(t *testing.T) {
	t.Run("유효한 UUID v4를 생성해야 함", func(t *testing.T) {
		id := GenerateID()

		// UUID 형식 검증
		parsedUUID, err := uuid.Parse(id)
		assert.NoError(t, err, "생성된 ID가 유효한 UUID 형식이어야 함")
		assert.Equal(t, uuid.Version(4), parsedUUID.Version(), "UUID 버전이 4여야 함")
	})

	t.Run("고유한 ID를 생성해야 함", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 1000; i++ {
			id := GenerateID()
			assert.False(t, ids[id], "중복된 ID가 생성되면 안 됨")
			ids[id] = true
		}
	})

	t.Run("빈 문자열을 반환하면 안 됨", func(t *testing.T) {
		id := GenerateID()
		assert.NotEmpty(t, id, "생성된 ID가 빈 문자열이면 안 됨")
	})

	t.Run("동시에 호출해도 안전해야 함", func(t *testing.T) {
		ids := make(chan string, 1000)
		done := make(chan bool)

		for i := 0; i < 10; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					ids <- GenerateID()
				}
				done <- true
			}()
		}

		// 모든 고루틴이 완료될 때까지 대기
		for i := 0; i < 10; i++ {
			<-done
		}
		close(ids)

		// 중복 검사
		uniqueIDs := make(map[string]bool)
		for id := range ids {
			assert.False(t, uniqueIDs[id], "동시 실행 시에도 중복된 ID가 생성되면 안 됨")
			uniqueIDs[id] = true
		}
	})
}
