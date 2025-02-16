package gamification

import (
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
)

// Level 사용자의 레벨을 나타냅니다.
type Level struct {
	Value       int
	Experience  int
	NextLevel   int
	Title       string
	UnlockedAt  time.Time
	Description string
}

// Profile 사용자의 게임화 프로필을 나타냅니다.
type Profile struct {
	ID           string
	UserID       string
	Level        Level
	Experience   int
	Achievements []string // Achievement IDs
	Badges       []Badge
	Streaks      []Streak
	Stats        Statistics
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (p *Profile) GetID() string {
	return p.ID
}

func (p *Profile) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Profile) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

// Badge 사용자가 획득한 뱃지를 나타냅니다.
type Badge struct {
	ID          string
	Type        BadgeType
	Title       string
	Description string
	Tier        BadgeTier
	UnlockedAt  time.Time
}

// BadgeType 뱃지의 유형을 나타냅니다.
type BadgeType string

const (
	BadgeTypeSaving    BadgeType = "SAVING"
	BadgeTypeInvesting BadgeType = "INVESTING"
	BadgeTypeCommunity BadgeType = "COMMUNITY"
	BadgeTypeChallenge BadgeType = "CHALLENGE"
)

// BadgeTier 뱃지의 등급을 나타냅니다.
type BadgeTier string

const (
	BadgeTierBronze  BadgeTier = "BRONZE"
	BadgeTierSilver  BadgeTier = "SILVER"
	BadgeTierGold    BadgeTier = "GOLD"
	BadgeTierDiamond BadgeTier = "DIAMOND"
)

// Streak 연속 달성을 나타냅니다.
type Streak struct {
	Type        StreakType
	Count       int
	LastUpdated time.Time
	MaxCount    int
}

// StreakType 연속 달성의 유형을 나타냅니다.
type StreakType string

const (
	StreakTypeDaily   StreakType = "DAILY"
	StreakTypeWeekly  StreakType = "WEEKLY"
	StreakTypeMonthly StreakType = "MONTHLY"
)

// Statistics 사용자의 게임화 통계를 나타냅니다.
type Statistics struct {
	TotalSavings      float64
	TotalInvestments  float64
	GoalsCompleted    int
	BadgesEarned      int
	LongestStreak     int
	CommunityRanking  int
	ContributionScore float64
}

// NewProfile 새로운 게임화 프로필을 생성합니다.
func NewProfile(userID string) *Profile {
	now := time.Now()
	return &Profile{
		ID:     domain.GenerateID(),
		UserID: userID,
		Level: Level{
			Value:       1,
			Experience:  0,
			NextLevel:   100,
			Title:       "초보 투자자",
			UnlockedAt:  now,
			Description: "자산 관리의 첫 걸음을 내딛었습니다.",
		},
		Experience:   0,
		Achievements: make([]string, 0),
		Badges:       make([]Badge, 0),
		Streaks:      make([]Streak, 0),
		Stats:        Statistics{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddExperience 경험치를 추가하고 레벨업 여부를 반환합니다.
func (p *Profile) AddExperience(exp int) bool {
	p.Experience += exp
	p.UpdatedAt = time.Now()

	leveledUp := false
	for p.Experience >= p.Level.NextLevel {
		p.levelUp()
		leveledUp = true
	}
	return leveledUp
}

// levelUp 레벨을 올립니다.
func (p *Profile) levelUp() {
	p.Level.Value++
	p.Level.Experience = p.Experience
	p.Level.NextLevel = calculateNextLevelExp(p.Level.Value)
	p.Level.UnlockedAt = time.Now()
	p.updateLevelTitle()
}

// calculateNextLevelExp 다음 레벨에 필요한 경험치를 계산합니다.
func calculateNextLevelExp(level int) int {
	return level * level * 100
}

// updateLevelTitle 레벨에 따른 타이틀을 업데이트합니다.
func (p *Profile) updateLevelTitle() {
	switch {
	case p.Level.Value >= 50:
		p.Level.Title = "투자의 신"
		p.Level.Description = "최고 수준의 투자 전문가입니다."
	case p.Level.Value >= 40:
		p.Level.Title = "투자 마스터"
		p.Level.Description = "뛰어난 투자 실력을 보유하고 있습니다."
	case p.Level.Value >= 30:
		p.Level.Title = "숙련된 투자자"
		p.Level.Description = "안정적인 투자 능력을 보유하고 있습니다."
	case p.Level.Value >= 20:
		p.Level.Title = "중급 투자자"
		p.Level.Description = "투자의 기본을 완벽히 이해했습니다."
	default:
		p.Level.Title = "초보 투자자"
		p.Level.Description = "투자의 기초를 배우고 있습니다."
	}
}

// AddBadge 뱃지를 추가합니다.
func (p *Profile) AddBadge(badgeType BadgeType, tier BadgeTier, title, description string) {
	badge := Badge{
		ID:          domain.GenerateID(),
		Type:        badgeType,
		Title:       title,
		Description: description,
		Tier:        tier,
		UnlockedAt:  time.Now(),
	}
	p.Badges = append(p.Badges, badge)
	p.Stats.BadgesEarned++
	p.UpdatedAt = time.Now()
}

// UpdateStreak 연속 달성을 업데이트합니다.
func (p *Profile) UpdateStreak(streakType StreakType) {
	now := time.Now()
	for i, streak := range p.Streaks {
		if streak.Type == streakType {
			if now.Sub(streak.LastUpdated) <= getStreakTimeout(streakType) {
				p.Streaks[i].Count++
				if p.Streaks[i].Count > p.Streaks[i].MaxCount {
					p.Streaks[i].MaxCount = p.Streaks[i].Count
				}
			} else {
				p.Streaks[i].Count = 1
			}
			p.Streaks[i].LastUpdated = now
			return
		}
	}

	// 새로운 스트릭 추가
	p.Streaks = append(p.Streaks, Streak{
		Type:        streakType,
		Count:       1,
		LastUpdated: now,
		MaxCount:    1,
	})
}

// getStreakTimeout 스트릭 타임아웃을 반환합니다.
func getStreakTimeout(streakType StreakType) time.Duration {
	switch streakType {
	case StreakTypeDaily:
		return 24 * time.Hour
	case StreakTypeWeekly:
		return 7 * 24 * time.Hour
	case StreakTypeMonthly:
		return 30 * 24 * time.Hour
	default:
		return 24 * time.Hour
	}
}
