package utils

import (
	"time"
)

type GenFixedTimePoint struct {
	fixedHour int
}

func NewGenFixTimePoint(fixedHour int) *GenFixedTimePoint {
	return &GenFixedTimePoint{
		fixedHour: fixedHour,
	}
}

func (gf *GenFixedTimePoint) GetFixedTimeStr() string {

	now := time.Now()
	newfixed := now
	if now.Hour() >= gf.fixedHour {
		newfixed = newfixed.Add(24 * time.Hour)
	}
	return newfixed.Format("2006-01-02")
}

func (gf *GenFixedTimePoint) GetPrevFixedTimeStr() string {

	now := time.Now().Add(-24 * time.Hour)
	newfixed := now
	if now.Hour() >= gf.fixedHour {
		newfixed = newfixed.Add(24 * time.Hour)
	}
	return newfixed.Format("2006-01-02")
}

func (gf *GenFixedTimePoint) GetFixedTime() time.Time {

	now := time.Now()
	newfixed := now
	if now.Hour() >= gf.fixedHour {
		newfixed = newfixed.Add(24 * time.Hour)
	}
	return newfixed
}

func GetMiddleNight() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
func GetNextMiddleNight() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
}
