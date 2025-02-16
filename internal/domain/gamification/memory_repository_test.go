package gamification

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewMemoryRepository_should_create_empty_repository(t *testing.T) {
	// When
	repo := NewMemoryRepository()

	// Then
	assert.NotNil(t, repo)
	assert.Empty(t, repo.data)
}

func Test_MemoryRepository_should_save_and_find_profile_by_id(t *testing.T) {
	// Given
	repo := NewMemoryRepository()
	profile := NewProfile("test-user")

	// When
	err := repo.Save(context.Background(), profile)
	assert.NoError(t, err)

	found, err := repo.FindByID(context.Background(), profile.ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, profile.ID, found.ID)
	assert.Equal(t, profile.UserID, found.UserID)
}

func Test_MemoryRepository_should_update_profile(t *testing.T) {
	// Given
	repo := NewMemoryRepository()
	profile := NewProfile("test-user")
	err := repo.Save(context.Background(), profile)
	assert.NoError(t, err)

	// When
	profile.AddExperience(100)
	err = repo.Update(context.Background(), profile)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), profile.ID)
	assert.NoError(t, err)
	assert.Equal(t, 100, found.Experience)
}

func Test_MemoryRepository_should_delete_profile(t *testing.T) {
	// Given
	repo := NewMemoryRepository()
	profile := NewProfile("test-user")
	err := repo.Save(context.Background(), profile)
	assert.NoError(t, err)

	// When
	err = repo.Delete(context.Background(), profile.ID)

	// Then
	assert.NoError(t, err)
	_, err = repo.FindByID(context.Background(), profile.ID)
	assert.Error(t, err)
}

func Test_MemoryRepository_should_find_profile_by_user_id(t *testing.T) {
	// Given
	repo := NewMemoryRepository()
	userID := "test-user"
	profile := NewProfile(userID)
	err := repo.Save(context.Background(), profile)
	assert.NoError(t, err)

	// When
	found, err := repo.FindByUserID(context.Background(), userID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, profile.ID, found.ID)
	assert.Equal(t, userID, found.UserID)
}

func Test_MemoryRepository_should_update_experience(t *testing.T) {
	// Given
	repo := NewMemoryRepository()
	profile := NewProfile("test-user")
	err := repo.Save(context.Background(), profile)
	assert.NoError(t, err)

	// When
	err = repo.UpdateExperience(context.Background(), profile.ID, 100)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), profile.ID)
	assert.NoError(t, err)
	assert.Equal(t, 100, found.Experience)
}

func Test_MemoryRepository_should_update_stats(t *testing.T) {
	// Given
	repo := NewMemoryRepository()
	profile := NewProfile("test-user")
	err := repo.Save(context.Background(), profile)
	assert.NoError(t, err)

	stats := Statistics{
		TotalSavings:      1000000,
		TotalInvestments:  5000000,
		GoalsCompleted:    5,
		BadgesEarned:      3,
		LongestStreak:     7,
		CommunityRanking:  10,
		ContributionScore: 85.5,
	}

	// When
	err = repo.UpdateStats(context.Background(), profile.ID, stats)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), profile.ID)
	assert.NoError(t, err)
	assert.Equal(t, stats, found.Stats)
}

func Test_MemoryRepository_should_add_badge(t *testing.T) {
	// Given
	repo := NewMemoryRepository()
	profile := NewProfile("test-user")
	err := repo.Save(context.Background(), profile)
	assert.NoError(t, err)

	badge := Badge{
		ID:          "test-badge",
		Type:        BadgeTypeSaving,
		Title:       "절약왕",
		Description: "첫 번째 저축 목표 달성",
		Tier:        BadgeTierBronze,
		UnlockedAt:  time.Now(),
	}

	// When
	err = repo.AddBadge(context.Background(), profile.ID, badge)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), profile.ID)
	assert.NoError(t, err)
	assert.Len(t, found.Badges, 1)
	assert.Equal(t, badge.ID, found.Badges[0].ID)
	assert.Equal(t, 1, found.Stats.BadgesEarned)
}

func Test_MemoryRepository_should_update_streak(t *testing.T) {
	// Given
	repo := NewMemoryRepository()
	profile := NewProfile("test-user")
	err := repo.Save(context.Background(), profile)
	assert.NoError(t, err)

	// When
	err = repo.UpdateStreak(context.Background(), profile.ID, StreakTypeDaily)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), profile.ID)
	assert.NoError(t, err)
	assert.Len(t, found.Streaks, 1)
	assert.Equal(t, StreakTypeDaily, found.Streaks[0].Type)
	assert.Equal(t, 1, found.Streaks[0].Count)
}

func Test_MemoryRepository_should_be_thread_safe(t *testing.T) {
	// Given
	repo := NewMemoryRepository()
	profile := NewProfile("test-user")
	err := repo.Save(context.Background(), profile)
	assert.NoError(t, err)

	iterations := 1000
	done := make(chan bool)

	// When
	go func() {
		for i := 0; i < iterations; i++ {
			_ = repo.UpdateExperience(context.Background(), profile.ID, 1)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			_ = repo.UpdateStreak(context.Background(), profile.ID, StreakTypeDaily)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			_, _ = repo.FindByID(context.Background(), profile.ID)
		}
		done <- true
	}()

	// Then
	for i := 0; i < 3; i++ {
		<-done
	}

	found, err := repo.FindByID(context.Background(), profile.ID)
	assert.NoError(t, err)
	assert.Equal(t, iterations, found.Experience)
}
