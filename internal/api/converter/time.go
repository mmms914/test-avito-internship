package converter

import (
	"fmt"
	"time"
)

var weekdaysFromInt = map[int]time.Weekday{
	1: time.Monday,
	2: time.Tuesday,
	3: time.Wednesday,
	4: time.Thursday,
	5: time.Friday,
	6: time.Saturday,
	7: time.Sunday,
}
var intFromWeekdays = map[time.Weekday]int{
	time.Monday:    1,
	time.Tuesday:   2,
	time.Wednesday: 3,
	time.Thursday:  4,
	time.Friday:    5,
	time.Saturday:  6,
	time.Sunday:    7,
}

func StringHourMinuteToTime(timeStr *string) (time.Duration, error) {
	t, err := time.Parse("15:04", *timeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid time format: %w", err)
	}

	duration := time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
	return duration, nil
}

func TimeHourMinuteToString(timeDur time.Duration) string {
	hours := int(timeDur.Hours())
	minutes := int(timeDur.Minutes())

	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func IntArrayToWeekdays(daysInt *[]int) ([]time.Weekday, error) {
	convertedDays := make([]time.Weekday, 0, len(*daysInt))
	for _, dayI := range *daysInt {
		if dayI < 1 || dayI > 7 {
			return nil, fmt.Errorf("invalid day %d", dayI)
		}

		convertedDays = append(convertedDays, weekdaysFromInt[dayI])
	}

	return convertedDays, nil
}

func WeekdaysToIntArray(days []time.Weekday) []int {
	convertedDays := make([]int, 0, len(days))
	for _, day := range days {
		convertedDays = append(convertedDays, intFromWeekdays[day])
	}

	return convertedDays
}
