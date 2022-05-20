package common

import (
	"time"
)

func GetFormattedTime(timeNow time.Time, timeFormat string) string {
	return timeNow.Format(timeFormat)
}
