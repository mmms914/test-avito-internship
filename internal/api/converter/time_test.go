package converter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mmms914/test-avito-internship/internal/api/converter"
	"github.com/mmms914/test-avito-internship/pkg/ptr"
)

func TestStringHourMinuteToTime(t *testing.T) {
	tests := map[string]struct {
		timeStr       *string
		expected      time.Duration
		expectedError bool
	}{
		"success - 09:00": {
			timeStr:       ptr.To("09:00"),
			expected:      9 * time.Hour,
			expectedError: false,
		},
		"success - 18:30": {
			timeStr:       ptr.To("18:30"),
			expected:      18*time.Hour + 30*time.Minute,
			expectedError: false,
		},
		"success - 00:00": {
			timeStr:       ptr.To("00:00"),
			expected:      0,
			expectedError: false,
		},
		"success - 23:59": {
			timeStr:       ptr.To("23:59"),
			expected:      23*time.Hour + 59*time.Minute,
			expectedError: false,
		},
		"failure - invalid format": {
			timeStr:       ptr.To("119:00"),
			expectedError: true,
		},
		"failure - wrong format": {
			timeStr:       ptr.To("09-00"),
			expectedError: true,
		},
		"failure - empty string": {
			timeStr:       ptr.To(""),
			expectedError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := converter.StringHourMinuteToTime(tt.timeStr)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestTimeHourMinuteToString(t *testing.T) {
	tests := map[string]struct {
		duration time.Duration
		expected string
	}{
		"09:00": {
			duration: 9 * time.Hour,
			expected: "09:00",
		},
		"18:30": {
			duration: 18*time.Hour + 30*time.Minute,
			expected: "18:30",
		},
		"00:00": {
			duration: 0,
			expected: "00:00",
		},
		"23:59": {
			duration: 23*time.Hour + 59*time.Minute,
			expected: "23:59",
		},
		"09:05": {
			duration: 9*time.Hour + 5*time.Minute,
			expected: "09:05",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := converter.TimeHourMinuteToString(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIntArrayToWeekdays(t *testing.T) {
	tests := map[string]struct {
		daysInt       *[]int
		expected      []time.Weekday
		expectedError bool
	}{
		"success - weekdays 1-5": {
			daysInt:       &[]int{1, 2, 3, 4, 5},
			expected:      []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			expectedError: false,
		},
		"success - weekend": {
			daysInt:       &[]int{6, 7},
			expected:      []time.Weekday{time.Saturday, time.Sunday},
			expectedError: false,
		},
		"success - single day": {
			daysInt:       &[]int{1},
			expected:      []time.Weekday{time.Monday},
			expectedError: false,
		},
		"success - all days": {
			daysInt: &[]int{1, 2, 3, 4, 5, 6, 7},
			expected: []time.Weekday{time.Monday, time.Tuesday,
				time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
			expectedError: false,
		},
		"success - empty array": {
			daysInt:       &[]int{},
			expected:      []time.Weekday{},
			expectedError: false,
		},
		"failure - invalid day 0": {
			daysInt:       &[]int{0},
			expectedError: true,
		},
		"failure - invalid day 8": {
			daysInt:       &[]int{8},
			expectedError: true,
		},
		"failure - mixed valid and invalid": {
			daysInt:       &[]int{1, 2, 9},
			expectedError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := converter.IntArrayToWeekdays(tt.daysInt)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestWeekdaysToIntArray(t *testing.T) {
	tests := map[string]struct {
		days     []time.Weekday
		expected []int
	}{
		"weekdays 1-5": {
			days:     []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			expected: []int{1, 2, 3, 4, 5},
		},
		"weekend": {
			days:     []time.Weekday{time.Saturday, time.Sunday},
			expected: []int{6, 7},
		},
		"single day": {
			days:     []time.Weekday{time.Monday},
			expected: []int{1},
		},
		"all days": {
			days: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday,
				time.Thursday, time.Friday, time.Saturday, time.Sunday},
			expected: []int{1, 2, 3, 4, 5, 6, 7},
		},
		"empty array": {
			days:     []time.Weekday{},
			expected: []int{},
		},
		"nil array": {
			days:     nil,
			expected: []int{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := converter.WeekdaysToIntArray(tt.days)
			assert.Equal(t, tt.expected, result)
		})
	}
}
