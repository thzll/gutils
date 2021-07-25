package myutils

import (
	"fmt"
)

func FormatTimeDuration(sec int64) string {
	second := sec % 60
	Minute := sec / 60 % 60
	Hour := sec / 3600 % 24
	Day := sec / (3600 * 24)
	if Day > 0 {
		return fmt.Sprintf("%d %d:%d:%d", Day, Hour, Minute, second)
	} else if Hour > 0 {
		return fmt.Sprintf("%d:%d:%d", Hour, Minute, second)
	} else if Minute > 0 {
		return fmt.Sprintf("%d:%d", Minute, second)
	} else {
		return fmt.Sprintf("%d", second)
	}
}
