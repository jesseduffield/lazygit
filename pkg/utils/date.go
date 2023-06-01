package utils

import (
	"fmt"
	"time"
)

func UnixToTimeAgo(timestamp int64) string {
	now := time.Now().Unix()
	return formatSecondsAgo(now - timestamp)
}

const (
	SECONDS_IN_SECOND = 1
	SECONDS_IN_MINUTE = 60
	SECONDS_IN_HOUR   = 3600
	SECONDS_IN_DAY    = 86400
	SECONDS_IN_WEEK   = 604800
	SECONDS_IN_YEAR   = 31536000
	SECONDS_IN_MONTH  = SECONDS_IN_YEAR / 12
)

type period struct {
	label           string
	secondsInPeriod int64
}

var periods = []period{
	{"s", SECONDS_IN_SECOND},
	{"m", SECONDS_IN_MINUTE},
	{"h", SECONDS_IN_HOUR},
	{"d", SECONDS_IN_DAY},
	{"w", SECONDS_IN_WEEK},
	// we're using 'm' for both minutes and months which is ambiguous but
	// disambiguating with another character feels like overkill.
	{"m", SECONDS_IN_MONTH},
	{"y", SECONDS_IN_YEAR},
}

func formatSecondsAgo(secondsAgo int64) string {
	for i, period := range periods {
		if i == 0 {
			continue
		}

		if secondsAgo < period.secondsInPeriod {
			return fmt.Sprintf("%d%s",
				secondsAgo/periods[i-1].secondsInPeriod,
				periods[i-1].label,
			)
		}
	}

	return fmt.Sprintf("%d%s",
		secondsAgo/periods[len(periods)-1].secondsInPeriod,
		periods[len(periods)-1].label,
	)
}

// formats the date in a smart way, if the date is today, it will show the time, otherwise it will show the date
func UnixToDateSmart(now time.Time, timestamp int64, longTimeFormat string, shortTimeFormat string) string {
	date := time.Unix(timestamp, 0)

	if date.Day() == now.Day() && date.Month() == now.Month() && date.Year() == now.Year() {
		return date.Format(shortTimeFormat)
	}

	return date.Format(longTimeFormat)
}
