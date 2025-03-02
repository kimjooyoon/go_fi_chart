package valueobjects

import (
	"fmt"
	"time"
)

// TimeRange 시간 범위를 나타냅니다.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// NewTimeRange TimeRange 값 객체를 생성합니다.
func NewTimeRange(start, end time.Time) (TimeRange, error) {
	if end.Before(start) {
		return TimeRange{}, fmt.Errorf("종료 시간은 시작 시간보다 이후여야 합니다")
	}
	return TimeRange{
		Start: start,
		End:   end,
	}, nil
}

// Duration 시간 범위의 기간을 반환합니다.
func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// Contains 주어진 시간이 범위 내에 있는지 확인합니다.
func (tr TimeRange) Contains(t time.Time) bool {
	return (t.Equal(tr.Start) || t.After(tr.Start)) && (t.Equal(tr.End) || t.Before(tr.End))
}

// Overlaps 다른 시간 범위와 겹치는지 확인합니다.
func (tr TimeRange) Overlaps(other TimeRange) bool {
	return (tr.Start.Before(other.End) && other.Start.Before(tr.End))
}

// Equals 두 시간 범위가 같은지 확인합니다.
func (tr TimeRange) Equals(other TimeRange) bool {
	return tr.Start.Equal(other.Start) && tr.End.Equal(other.End)
}

// String 시간 범위를 문자열로 변환합니다.
func (tr TimeRange) String() string {
	return fmt.Sprintf("%s - %s", tr.Start.Format(time.RFC3339), tr.End.Format(time.RFC3339))
}

// IsZero 시간 범위가 비어있는지 확인합니다.
func (tr TimeRange) IsZero() bool {
	return tr.Start.IsZero() && tr.End.IsZero()
}

// Extend 시간 범위를 주어진 기간만큼 연장합니다.
func (tr TimeRange) Extend(d time.Duration) TimeRange {
	return TimeRange{
		Start: tr.Start,
		End:   tr.End.Add(d),
	}
}

// Shift 시간 범위를 주어진 기간만큼 이동합니다.
func (tr TimeRange) Shift(d time.Duration) TimeRange {
	return TimeRange{
		Start: tr.Start.Add(d),
		End:   tr.End.Add(d),
	}
}

// Split은 시간 범위를 주어진 간격으로 분할합니다.
func (tr TimeRange) Split(interval time.Duration) []TimeRange {
	if interval <= 0 {
		return []TimeRange{}
	}

	if interval >= tr.Duration() {
		return []TimeRange{tr}
	}

	var ranges []TimeRange
	start := tr.Start
	for start.Before(tr.End) {
		end := start.Add(interval)
		if end.After(tr.End) {
			end = tr.End
		}
		ranges = append(ranges, TimeRange{Start: start, End: end})
		start = end
	}
	return ranges
}
