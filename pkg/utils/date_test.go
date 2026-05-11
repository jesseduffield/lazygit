package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatSecondsAgo(t *testing.T) {
	tests := []struct {
		name string
		args int64
		want string
	}{
		{
			name: "zero",
			args: 0,
			want: "0s",
		},
		{
			name: "one second",
			args: 1,
			want: "1s",
		},
		{
			name: "almost a minute",
			args: 59,
			want: "59s",
		},
		{
			name: "one minute",
			args: 60,
			want: "1m",
		},
		{
			name: "one minute and one second",
			args: 61,
			want: "1m",
		},
		{
			name: "almost one hour",
			args: 3599,
			want: "59m",
		},
		{
			name: "one hour",
			args: 3600,
			want: "1h",
		},
		{
			name: "almost one day",
			args: 86399,
			want: "23h",
		},
		{
			name: "one day",
			args: 86400,
			want: "1d",
		},
		{
			name: "almost a week",
			args: 604799,
			want: "6d",
		},
		{
			name: "one week",
			args: 604800,
			want: "1w",
		},
		{
			name: "six months",
			args: SECONDS_IN_YEAR / 2,
			want: "6M",
		},
		{
			name: "almost one year",
			args: 31535999,
			want: "11M",
		},
		{
			name: "one year",
			args: SECONDS_IN_YEAR,
			want: "1y",
		},
		{
			name: "50 years",
			args: SECONDS_IN_YEAR * 50,
			want: "50y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatSecondsAgo(tt.args); got != tt.want {
				t.Errorf("formatSecondsAgo(%d) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestUnixToDateSmart_UsesNowLocationForFormatting(t *testing.T) {
	timestamp := int64(1577844184) // 2020-01-01 02:03:04 UTC
	now := time.Date(2020, 1, 1, 5, 3, 4, 0, time.UTC)

	assert.Equal(t, "2:03AM", UnixToDateSmart(now, timestamp, "2006-01-02", "3:04PM"))
}

func TestUnixToDateSmart_SameTimestampDifferentNowLocation(t *testing.T) {
	timestamp := int64(1577844184) // 2020-01-01 02:03:04 UTC
	loc := time.FixedZone("UTC+1", 3600)
	now := time.Date(2020, 1, 1, 6, 3, 4, 0, loc)

	assert.Equal(t, "3:03AM", UnixToDateSmart(now, timestamp, "2006-01-02", "3:04PM"))
}
