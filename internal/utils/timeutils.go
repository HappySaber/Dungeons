package timeutil

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func ParseTime(s string) (time.Time, error) {
	return time.Parse("15:04:05", s)
}

func ParseDuration(s string) (time.Duration, error) {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return 0, err
	}

	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second, nil
}
