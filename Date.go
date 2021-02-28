package main

import (
	"fmt"
	"math"
	"time"
)

func timeSince(tm time.Time) string {
	t := int(math.Round(time.Now().Sub(tm).Seconds()))
	ago := ""
	secondsInADay := 3600 * 24
	if t < 120 {
		ago = fmt.Sprintf("%ds", t)
	} else if t < 3600 {
		mins := int(t / 60)
		secs := t % 60
		ago = fmt.Sprintf("%dm%ds", mins, secs)
	} else if t < secondsInADay {
		hours := int(t / 3600)
		t -= 3600 * hours
		mins := int(t / 60)
		ago = fmt.Sprintf("%dh%dm", hours, mins)
	} else if t < secondsInADay*5 {
		days := int(t / secondsInADay)
		t -= days * secondsInADay
		hours := int(t / 3600)
		t -= 3600 * hours
		ago = fmt.Sprintf("%dd%dh", days, hours)
	} else {
		days := int(t / secondsInADay)
		t -= days * secondsInADay
		hours := int(t / 3600)
		t -= 3600 * hours
		ago = fmt.Sprintf("%dd", days)
	}
	return ago
}
