package timeutils

import "time"

//go:generate mockery --name TimeUtils --filename timeutils.go
type TimeUtils interface {
	GetCurrentTime() time.Time
	GetCurrentTimeStr(format string) string
	GetDayOffsetTime(dayOffset int) time.Time
	GetDayOffsetTimeStr(dayOffset int, format string) string
	GetCurrentTimeUnix() int64
	GetDayOffsetTimeUnix(dayOffset int) int64
}

func NewTimeUtils() TimeUtils {
	return &timeUtils{}
}

type timeUtils struct{}

func (t *timeUtils) GetCurrentTime() time.Time {
	return time.Now().UTC()
}

func (t *timeUtils) GetCurrentTimeStr(format string) string {
	ti := t.GetCurrentTime()
	return ti.Format(format)
}

func (t *timeUtils) GetDayOffsetTime(dayOffset int) time.Time {
	return time.Now().UTC().AddDate(0, 0, dayOffset)
}

func (t *timeUtils) GetDayOffsetTimeStr(dayOffset int, format string) string {
	ti := t.GetDayOffsetTime(dayOffset)
	return ti.Format(format)
}

func (t *timeUtils) GetCurrentTimeUnix() int64 {
	ti := t.GetCurrentTime()
	return ti.Unix()
}

func (t *timeUtils) GetDayOffsetTimeUnix(dayOffset int) int64 {
	ti := t.GetDayOffsetTime(dayOffset)
	return ti.Unix()
}

func MustParse(layout, value string) time.Time {
	v, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return v
}

func MustParseRFC3339(value string) time.Time {
	v, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}

	return v
}
