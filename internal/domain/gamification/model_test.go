package gamification

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewProfile_should_create_profile_with_valid_data(t *testing.T) {
	// Given
	userID := "test-user"

	// When
	profile := NewProfile(userID)

	// Then
	assert.NotEmpty(t, profile.ID)
	assert.Equal(t, userID, profile.UserID)
	assert.Equal(t, 1, profile.Level.Value)
	assert.Equal(t, 0, profile.Experience)
	assert.Equal(t, "초보 투자자", profile.Level.Title)
	assert.Empty(t, profile.Achievements)
	assert.Empty(t, profile.Badges)
	assert.Empty(t, profile.Streaks)
	assert.NotZero(t, profile.CreatedAt)
	assert.NotZero(t, profile.UpdatedAt)
}

func Test_Profile_AddExperience_should_level_up_when_enough_exp(t *testing.T) {
	// Given
	profile := NewProfile("test-user")
	initialLevel := profile.Level.Value
	expToAdd := 100 // 레벨 2로 레벨업하기 위한 경험치

	// When
	leveledUp := profile.AddExperience(expToAdd)

	// Then
	assert.True(t, leveledUp)
	assert.Equal(t, initialLevel+1, profile.Level.Value)
	assert.Equal(t, expToAdd, profile.Experience)
	assert.Equal(t, "초보 투자자", profile.Level.Title)
}

func Test_Profile_AddExperience_should_not_level_up_when_not_enough_exp(t *testing.T) {
	// Given
	profile := NewProfile("test-user")
	initialLevel := profile.Level.Value
	expToAdd := 50 // 다음 레벨에 필요한 경험치보다 적음

	// When
	leveledUp := profile.AddExperience(expToAdd)

	// Then
	assert.False(t, leveledUp)
	assert.Equal(t, initialLevel, profile.Level.Value)
	assert.Equal(t, expToAdd, profile.Experience)
	assert.Equal(t, "초보 투자자", profile.Level.Title)
}

func Test_Profile_AddBadge_should_add_badge_and_update_stats(t *testing.T) {
	// Given
	profile := NewProfile("test-user")
	initialBadgeCount := len(profile.Badges)

	// When
	profile.AddBadge(BadgeTypeSaving, BadgeTierBronze, "절약왕", "첫 번째 저축 목표 달성")

	// Then
	assert.Len(t, profile.Badges, initialBadgeCount+1)
	assert.Equal(t, 1, profile.Stats.BadgesEarned)

	badge := profile.Badges[0]
	assert.Equal(t, BadgeTypeSaving, badge.Type)
	assert.Equal(t, BadgeTierBronze, badge.Tier)
	assert.Equal(t, "절약왕", badge.Title)
	assert.Equal(t, "첫 번째 저축 목표 달성", badge.Description)
	assert.NotZero(t, badge.UnlockedAt)
}

func Test_Profile_UpdateStreak_should_increment_existing_streak(t *testing.T) {
	// Given
	profile := NewProfile("test-user")
	profile.UpdateStreak(StreakTypeDaily)
	initialCount := profile.Streaks[0].Count

	// When
	profile.UpdateStreak(StreakTypeDaily)

	// Then
	assert.Len(t, profile.Streaks, 1)
	assert.Equal(t, initialCount+1, profile.Streaks[0].Count)
	assert.Equal(t, initialCount+1, profile.Streaks[0].MaxCount)
}

func Test_Profile_UpdateStreak_should_reset_when_expired(t *testing.T) {
	// Given
	profile := NewProfile("test-user")
	profile.UpdateStreak(StreakTypeDaily)
	streak := &profile.Streaks[0]
	streak.LastUpdated = time.Now().Add(-25 * time.Hour) // 하루가 지난 상태

	// When
	profile.UpdateStreak(StreakTypeDaily)

	// Then
	assert.Len(t, profile.Streaks, 1)
	assert.Equal(t, 1, profile.Streaks[0].Count)
	assert.Equal(t, 1, profile.Streaks[0].MaxCount)
}

func Test_Profile_UpdateStreak_should_create_new_streak_type(t *testing.T) {
	// Given
	profile := NewProfile("test-user")
	profile.UpdateStreak(StreakTypeDaily)
	initialStreakCount := len(profile.Streaks)

	// When
	profile.UpdateStreak(StreakTypeWeekly)

	// Then
	assert.Len(t, profile.Streaks, initialStreakCount+1)
	assert.Equal(t, StreakTypeWeekly, profile.Streaks[1].Type)
	assert.Equal(t, 1, profile.Streaks[1].Count)
}

func Test_Profile_should_update_level_title_based_on_level(t *testing.T) {
	// Given
	profile := NewProfile("test-user")
	testCases := []struct {
		exp           int
		expectedTitle string
	}{
		{40000, "중급 투자자"},  // Level 20
		{90000, "숙련된 투자자"}, // Level 30
		{160000, "투자 마스터"}, // Level 40
		{250000, "투자의 신"},  // Level 50
	}

	for _, tc := range testCases {
		// When
		profile.Experience = 0 // 각 테스트 케이스마다 경험치 초기화
		profile.AddExperience(tc.exp)

		// Then
		assert.Equal(t, tc.expectedTitle, profile.Level.Title)
	}
}
