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

// formats the date in a smart way, if the date is today, it will show the time, otherwise it will show the date
func UnixToDateSmart(now time.Time, timestamp int64, longTimeFormat string, shortTimeFormat string) string {
	date := time.Unix(timestamp, 0)

	if date.Day() == now.Day() && date.Month() == now.Month() && date.Year() == now.Year() {
		return date.Format(shortTimeFormat)
	}

	return date.Format(longTimeFormat)
}
