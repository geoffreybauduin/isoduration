package isoduration_test

import (
	"testing"

	"github.com/geoffreybauduin/isoduration"
	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	t.Parallel()

	// test with bad format
	_, err := isoduration.FromString("asdf")
	assert.Equal(t, err, isoduration.ErrBadFormat)

	// test with good full string
	dur, err := isoduration.FromString("P1Y2DT3H4M5S")
	assert.Nil(t, err)
	assert.Equal(t, 1, dur.Years)
	assert.Equal(t, 2, dur.Days)
	assert.Equal(t, 3, dur.Hours)
	assert.Equal(t, 4, dur.Minutes)
	assert.Equal(t, 5, dur.Seconds)

	// test with good week string
	dur, err = isoduration.FromString("P1W")
	assert.Nil(t, err)
	assert.Equal(t, 1, dur.Weeks)

	// test with month
	dur, err = isoduration.FromString("P1M")
	assert.NoError(t, err)
	assert.Equal(t, 1, dur.Months)
}

func TestString(t *testing.T) {
	t.Parallel()

	// test empty
	d := isoduration.Duration{}
	assert.Equal(t, d.String(), "P")

	// test only larger-than-day
	d = isoduration.Duration{Years: 1, Days: 2}
	assert.Equal(t, d.String(), "P1Y2D")

	// test only smaller-than-day
	d = isoduration.Duration{Hours: 1, Minutes: 2, Seconds: 3}
	assert.Equal(t, d.String(), "PT1H2M3S")

	// test full format
	d = isoduration.Duration{Years: 1, Days: 2, Hours: 3, Minutes: 4, Seconds: 5}
	assert.Equal(t, d.String(), "P1Y2DT3H4M5S")

	// test week format
	d = isoduration.Duration{Weeks: 1}
	assert.Equal(t, d.String(), "P1W")
}
