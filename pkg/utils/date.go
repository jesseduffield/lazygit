package utils

import (
	"fmt"
	"time"
)

func UnixToTimeAgo(timestamp int64) string {
	now := time.Now().Unix()
	delta := float64(now - timestamp)
	// we go seconds, minutes, hours, days, weeks, months, years
	conversions := []float64{60, 60, 24, 7, 4.34524, 12}
	labels := []string{"s", "m", "h", "d", "w", "m", "y"}
	for i, conversion := range conversions {
		if delta < conversion {
			return fmt.Sprintf("%d%s", int(delta), labels[i])
		}
		delta /= conversion
	}
	return fmt.Sprintf("%dy", int(delta))
}

func UnixToDate(timestamp int64, timeFormat string) string {
	return time.Unix(timestamp, 0).Format(timeFormat)
}
