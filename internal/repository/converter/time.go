package converter

import (
	"fmt"
	"time"
)

func IntArrayToWeekdays(daysInt *[]int) ([]time.Weekday, error) {
	weekdaysFromInt := map[int]time.Weekday{
		1: time.Monday,
		2: time.Tuesday,
		3: time.Wednesday,
		4: time.Thursday,
		5: time.Friday,
		6: time.Saturday,
		7: time.Sunday,
	}

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
	//nolint:mnd // useless
	intFromWeekdays := map[time.Weekday]int{
		time.Monday:    1,
		time.Tuesday:   2,
		time.Wednesday: 3,
		time.Thursday:  4,
		time.Friday:    5,
		time.Saturday:  6,
		time.Sunday:    7,
	}

	convertedDays := make([]int, 0, len(days))
	for _, day := range days {
		convertedDays = append(convertedDays, intFromWeekdays[day])
	}

	return convertedDays
}

func DurationToStringTime(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60 //nolint:mnd // obviously
	seconds := int(d.Seconds()) % 60 //nolint:mnd // obviously
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func StringToDuration(timeStr string) (time.Duration, error) {
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		t, err = time.Parse("15:04", timeStr)
		if err != nil {
			return 0, fmt.Errorf("invalid time format: %w", err)
		}
	}

	duration := time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second

	return duration, nil
}
