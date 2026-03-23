package converter

import (
	"fmt"
	"time"
)

var weekdays = map[int]time.Weekday{
	1: time.Monday,
	2: time.Tuesday,
	3: time.Wednesday,
	4: time.Thursday,
	5: time.Friday,
	6: time.Saturday,
	7: time.Sunday,
}

func StringHourMinuteToTime(timeStr *string) (time.Duration, error) {
	t, err := time.Parse("15:04", *timeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid time format: %w", err)
	}

	duration := time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
	return duration, nil
}

func IntArrayToWeekdays(daysInt *[]int) ([]time.Weekday, error) {
	convertedDays := make([]time.Weekday, 0, len(*daysInt))
	for _, dayI := range *daysInt {
		if dayI < 1 || dayI > 7 {
			return nil, fmt.Errorf("invalid day %d", dayI)
		}

		convertedDays = append(convertedDays, weekdays[dayI])
	}

	return convertedDays, nil
}
