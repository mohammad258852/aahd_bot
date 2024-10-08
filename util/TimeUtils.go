package util

import (
	"log"
	"time"
)

var tehranTime *time.Location

func LoadTehranTime() *time.Location {
	if tehranTime != nil {
		return tehranTime
	}
	var err error
	tehranTime, err = time.LoadLocation("Asia/Tehran")
	if err != nil {
		log.Print(err)
	}
	return tehranTime
}

func GetDurationTilNextDay() time.Duration {
	t := GetCurrentLocalTime()
	n := time.Date(t.Year(), t.Month(), t.Day(), 0, 1, 0, 0, t.Location())
	d := n.Sub(t)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	return d
}

func GetDurationTilNextWeekDay(weekday time.Weekday) time.Duration {
	t := GetCurrentLocalTime()
	n := time.Date(t.Year(), t.Month(), t.Day(), 12, 15, 0, 0, t.Location())
	if n.Sub(t) < 0 {
		n = n.Add(24 * time.Hour)
	}
	for n.Weekday() != weekday {
		n = n.Add(24 * time.Hour)
	}
	return n.Sub(t)
}

func GetCurrentLocalTime() time.Time {
	return time.Now().In(LoadTehranTime())
}
