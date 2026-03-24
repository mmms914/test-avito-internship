package converter

import (
	"fmt"
	"strconv"
	"strings"
)

func StringToIntArray(s string) ([]int, error) {
	if s == "" {
		return []int{}, nil
	}

	s = strings.Trim(s, "{}")
	if s == "" {
		return []int{}, nil
	}

	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("failed to parse '%s' as int: %w", part, err)
		}
		result = append(result, num)
	}

	return result, nil
}
