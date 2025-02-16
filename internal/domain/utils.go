package domain

import (
	"github.com/google/uuid"
)

// GenerateID UUID v4를 사용하여 고유한 ID를 생성합니다.
func GenerateID() string {
	return uuid.New().String()
}
