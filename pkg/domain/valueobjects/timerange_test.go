package valueobjects

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTimeRange(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)
	earlier := now.Add(-24 * time.Hour)

	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		wantErr bool
	}{
		{
			name:    "유효한 시간 범위",
			start:   now,
			end:     later,
			wantErr: false,
		},
		{
			name:    "시작과 종료가 같은 경우",
			start:   now,
			end:     now,
			wantErr: false,
		},
		{
			name:    "종료가 시작보다 이전인 경우",
			start:   now,
			end:     earlier,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTimeRange(tt.start, tt.end)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.start, got.Start)
				assert.Equal(t, tt.end, got.End)
			}
		})
	}
}

func TestTimeRange_Duration(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)

	tr, _ := NewTimeRange(now, later)
	duration := tr.Duration()

	assert.Equal(t, 24*time.Hour, duration)
}

func TestTimeRange_Contains(t *testing.T) {
	now := time.Now()
	middle := now.Add(12 * time.Hour)
	later := now.Add(24 * time.Hour)
	outside := now.Add(48 * time.Hour)

	tr, _ := NewTimeRange(now, later)

	assert.True(t, tr.Contains(now))      // 시작 시간
	assert.True(t, tr.Contains(middle))   // 중간 시간
	assert.True(t, tr.Contains(later))    // 종료 시간
	assert.False(t, tr.Contains(outside)) // 범위 밖 시간
}

func TestTimeRange_Overlaps(t *testing.T) {
	now := time.Now()
	middle := now.Add(12 * time.Hour)
	later := now.Add(24 * time.Hour)
	afterLater := now.Add(36 * time.Hour)

	tr1, _ := NewTimeRange(now, later)
	tr2, _ := NewTimeRange(middle, afterLater) // 겹치는 범위
	tr3, _ := NewTimeRange(later, afterLater)  // 연속된 범위
	tr4, _ := NewTimeRange(now, middle)        // 포함된 범위

	assert.True(t, tr1.Overlaps(tr2))  // 부분 겹침
	assert.True(t, tr1.Overlaps(tr4))  // 포함 관계
	assert.False(t, tr1.Overlaps(tr3)) // 연속된 범위
}

func TestTimeRange_Equals(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)

	tr1, _ := NewTimeRange(now, later)
	tr2, _ := NewTimeRange(now, later)
	tr3, _ := NewTimeRange(now, now.Add(12*time.Hour))

	assert.True(t, tr1.Equals(tr2))  // 동일한 범위
	assert.False(t, tr1.Equals(tr3)) // 다른 범위
}

func TestTimeRange_String(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)

	tr, _ := NewTimeRange(now, later)
	expected := now.Format(time.RFC3339) + " - " + later.Format(time.RFC3339)

	assert.Equal(t, expected, tr.String())
}

func TestTimeRange_IsZero(t *testing.T) {
	zeroTR := TimeRange{}
	assert.True(t, zeroTR.IsZero())

	now := time.Now()
	later := now.Add(24 * time.Hour)
	tr, _ := NewTimeRange(now, later)
	assert.False(t, tr.IsZero())
}

func TestTimeRange_Extend(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)
	extension := 12 * time.Hour

	tr, _ := NewTimeRange(now, later)
	extended := tr.Extend(extension)

	assert.Equal(t, now, extended.Start)
	assert.Equal(t, later.Add(extension), extended.End)
}

func TestTimeRange_Shift(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)
	shift := 12 * time.Hour

	tr, _ := NewTimeRange(now, later)
	shifted := tr.Shift(shift)

	assert.Equal(t, now.Add(shift), shifted.Start)
	assert.Equal(t, later.Add(shift), shifted.End)
}
