package isoduration

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"text/template"
)

var (
	// ErrBadFormat is returned when parsing fails
	ErrBadFormat = errors.New("bad format")

	tmpl = template.Must(template.New("duration").Parse(`P{{if .Years}}{{.Years}}Y{{end}}{{if .Weeks}}{{.Weeks}}W{{end}}{{if .Days}}{{.Days}}D{{end}}{{if .HasTimePart}}T{{end }}{{if .Hours}}{{.Hours}}H{{end}}{{if .Minutes}}{{.Minutes}}M{{end}}{{if .Seconds}}{{.Seconds}}S{{end}}`))

	full  = regexp.MustCompile(`P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+)S)?)?`)
	week  = regexp.MustCompile(`P((?P<week>\d+)W)`)
	month = regexp.MustCompile(`P((?P<month>\d+)M)`)
)

type Duration struct {
	Years   int
	Months  int
	Weeks   int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// FromString parses the ISO duration string into a Duration object
func FromString(dur string) (*Duration, error) {
	var (
		match []string
		re    *regexp.Regexp
	)

	if week.MatchString(dur) {
		match = week.FindStringSubmatch(dur)
		re = week
	} else if month.MatchString(dur) {
		match = month.FindStringSubmatch(dur)
		re = month
	} else if full.MatchString(dur) {
		match = full.FindStringSubmatch(dur)
		re = full
	} else {
		return nil, ErrBadFormat
	}

	d := &Duration{}

	for i, name := range re.SubexpNames() {
		part := match[i]
		if i == 0 || name == "" || part == "" {
			continue
		}

		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		switch name {
		case "year":
			d.Years = val
		case "month":
			d.Months = val
		case "week":
			d.Weeks = val
		case "day":
			d.Days = val
		case "hour":
			d.Hours = val
		case "minute":
			d.Minutes = val
		case "second":
			d.Seconds = val
		default:
			return nil, errors.New(fmt.Sprintf("unknown field %s", name))
		}
	}

	return d, nil
}

// String prints out the value passed in. It's not strictly according to the
// ISO spec, but it's pretty close. In particular, to completely conform it
// would need to round up to the next largest unit. 61 seconds to 1 minute 1
// second, for example. It would also need to disallow weeks mingling with
// other units.
func (d *Duration) String() string {
	var s bytes.Buffer

	err := tmpl.Execute(&s, d)
	if err != nil {
		panic(err)
	}

	return s.String()
}

// HasTimePart returns true if the current duration has a time part. Used for templating
func (d *Duration) HasTimePart() bool {
	return d.Hours != 0 || d.Minutes != 0 || d.Seconds != 0
}
