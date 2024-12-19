package timeutils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/datngo2sgtech/go-packages/timeutils"
)

func TestTimeUtils_GetCurrentTime(t *testing.T) {
	t.Parallel()

	ti := timeutils.NewTimeUtils()
	assert.Equal(t, ti.GetCurrentTime().Round(time.Second), time.Now().UTC().Round(time.Second))
}

func TestTimeUtils_GetCurrentTimeStr(t *testing.T) {
	t.Parallel()

	ti := timeutils.NewTimeUtils()
	format := "2024-02-01"
	assert.Equal(t, ti.GetCurrentTimeStr(format), time.Now().UTC().Format(format))
}

func TestTimeUtils_GetDayOffsetTime(t *testing.T) {
	t.Parallel()

	ti := timeutils.NewTimeUtils()
	offset := -3
	assert.Equal(
		t,
		ti.GetDayOffsetTime(offset).Round(time.Second),
		time.Now().UTC().AddDate(0, 0, offset).Round(time.Second),
	)
}

func TestTimeUtils_GetDayOffsetTimeStr(t *testing.T) {
	t.Parallel()

	ti := timeutils.NewTimeUtils()
	format := "2023-01-02"
	assert.Equal(
		t,
		ti.GetDayOffsetTimeStr(-3, format),
		time.Now().UTC().AddDate(0, 0, -3).Format(format),
	)
}

func TestTimeUtils_GetCurrentTimeUnix(t *testing.T) {
	t.Parallel()

	ti := timeutils.NewTimeUtils()
	timeUnix := ti.GetCurrentTimeUnix()
	assert.Equal(t, time.Now().UTC().Unix(), timeUnix)
}

func TestTimeUtils_GetDayOffsetTimeUnix(t *testing.T) {
	t.Parallel()

	ti := timeutils.NewTimeUtils()
	offset := -3
	timeUnix := ti.GetDayOffsetTimeUnix(offset)
	assert.Equal(t, ti.GetDayOffsetTime(offset).Unix(), timeUnix)
}
