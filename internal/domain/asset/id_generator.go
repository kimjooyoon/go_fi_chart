package asset

import (
	"github.com/google/uuid"
)

// generateID는 UUID v4를 사용하여 고유한 ID를 생성합니다.
// UUID v4는 랜덤하게 생성되며, 충돌 가능성이 극히 낮습니다.
// 또한 내부적으로 동시성을 지원하므로 별도의 동기화가 필요하지 않습니다.
func generateID() string {
	return uuid.New().String()
}
