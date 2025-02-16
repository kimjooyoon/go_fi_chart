package gamification

import (
	"context"

	"github.com/aske/go_fi_chart/internal/domain"
)

// Repository 게임화 프로필의 저장소 인터페이스입니다.
type Repository interface {
	domain.Repository[*Profile, string]

	// FindByUserID 사용자 ID로 게임화 프로필을 조회합니다.
	FindByUserID(ctx context.Context, userID string) (*Profile, error)

	// UpdateExperience 사용자의 경험치를 업데이트합니다.
	UpdateExperience(ctx context.Context, id string, exp int) error

	// UpdateStats 사용자의 통계를 업데이트합니다.
	UpdateStats(ctx context.Context, id string, stats Statistics) error

	// AddBadge 사용자에게 뱃지를 추가합니다.
	AddBadge(ctx context.Context, id string, badge Badge) error

	// UpdateStreak 사용자의 연속 달성을 업데이트합니다.
	UpdateStreak(ctx context.Context, id string, streakType StreakType) error
}
