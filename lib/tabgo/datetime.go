package tabgo

import (
	"fmt"
	"time"
)

// DatetimeIndex represents an index of datetime values with an optional frequency.
type DatetimeIndex struct {
	times []time.Time
	name  string
	freq  string
}

// NewDatetimeIndex creates a new DatetimeIndex from a slice of times and a name.
func NewDatetimeIndex(times []time.Time, name string) *DatetimeIndex {
	cp := make([]time.Time, len(times))
	copy(cp, times)
	return &DatetimeIndex{times: cp, name: name}
}

// Len returns the number of elements in the index.
func (di *DatetimeIndex) Len() int {
	return len(di.times)
}

// Name returns the name of the index.
func (di *DatetimeIndex) Name() string {
	return di.name
}

// Times returns a copy of the underlying time values.
func (di *DatetimeIndex) Times() []time.Time {
	cp := make([]time.Time, len(di.times))
	copy(cp, di.times)
	return cp
}

// Freq returns the frequency string of the index.
func (di *DatetimeIndex) Freq() string {
	return di.freq
}

// Get returns the time at position i.
func (di *DatetimeIndex) Get(i int) time.Time {
	return di.times[i]
}

// Values returns a copy of the underlying times (alias for Times).
func (di *DatetimeIndex) Values() []time.Time {
	return di.Times()
}

// Slice returns a new DatetimeIndex containing elements [start, end).
func (di *DatetimeIndex) Slice(start, end int) *DatetimeIndex {
	if start < 0 {
		start = 0
	}
	if end > len(di.times) {
		end = len(di.times)
	}
	cp := make([]time.Time, end-start)
	copy(cp, di.times[start:end])
	return &DatetimeIndex{times: cp, name: di.name, freq: di.freq}
}

// ToSeries converts the DatetimeIndex to a Series of time.Time values.
func (di *DatetimeIndex) ToSeries() *Series {
	vals := make([]any, len(di.times))
	for i, t := range di.times {
		vals[i] = t
	}
	return NewSeries(di.name, vals)
}

// Year returns a Series with the year of each element.
func (di *DatetimeIndex) Year() *Series {
	vals := make([]any, len(di.times))
	for i, t := range di.times {
		vals[i] = t.Year()
	}
	return &Series{name: di.name, values: vals}
}

// Month returns a Series with the month (1-12) of each element.
func (di *DatetimeIndex) Month() *Series {
	vals := make([]any, len(di.times))
	for i, t := range di.times {
		vals[i] = int(t.Month())
	}
	return &Series{name: di.name, values: vals}
}

// Day returns a Series with the day of month of each element.
func (di *DatetimeIndex) Day() *Series {
	vals := make([]any, len(di.times))
	for i, t := range di.times {
		vals[i] = t.Day()
	}
	return &Series{name: di.name, values: vals}
}

// Hour returns a Series with the hour of each element.
func (di *DatetimeIndex) Hour() *Series {
	vals := make([]any, len(di.times))
	for i, t := range di.times {
		vals[i] = t.Hour()
	}
	return &Series{name: di.name, values: vals}
}

// Weekday returns a Series with the weekday (0=Sunday) of each element.
func (di *DatetimeIndex) Weekday() *Series {
	vals := make([]any, len(di.times))
	for i, t := range di.times {
		vals[i] = int(t.Weekday())
	}
	return &Series{name: di.name, values: vals}
}

// String returns a string representation of the DatetimeIndex.
func (di *DatetimeIndex) String() string {
	if len(di.times) == 0 {
		return fmt.Sprintf("DatetimeIndex([], name=%q)", di.name)
	}
	n := len(di.times)
	limit := n
	if limit > 5 {
		limit = 5
	}
	s := "DatetimeIndex(["
	for i := 0; i < limit; i++ {
		if i > 0 {
			s += ", "
		}
		s += di.times[i].Format(time.RFC3339)
	}
	if n > 5 {
		s += fmt.Sprintf(", ... (%d more)", n-5)
	}
	s += fmt.Sprintf("], name=%q, freq=%q)", di.name, di.freq)
	return s
}

// freqDuration returns the duration for a single period of the given frequency
// relative to a reference time. For "M" (monthly) it returns the number of days
// in the reference month; for all other frequencies the reference is unused.
func freqDuration(freq string, ref time.Time) time.Duration {
	switch freq {
	case "T": // minute
		return time.Minute
	case "H": // hourly
		return time.Hour
	case "D": // daily
		return 24 * time.Hour
	case "W": // weekly
		return 7 * 24 * time.Hour
	case "M": // monthly — variable length
		next := time.Date(ref.Year(), ref.Month()+1, ref.Day(), ref.Hour(), ref.Minute(), ref.Second(), ref.Nanosecond(), ref.Location())
		return next.Sub(ref)
	default:
		return 24 * time.Hour // default to daily
	}
}

// addFreq advances t by n periods of the given frequency.
func addFreq(t time.Time, n int, freq string) time.Time {
	switch freq {
	case "M":
		return t.AddDate(0, n, 0)
	case "W":
		return t.AddDate(0, 0, 7*n)
	case "D":
		return t.AddDate(0, 0, n)
	case "H":
		return t.Add(time.Duration(n) * time.Hour)
	case "T":
		return t.Add(time.Duration(n) * time.Minute)
	default:
		return t.AddDate(0, 0, n) // default daily
	}
}

// truncateToFreq truncates a time to the start of the period defined by freq.
func truncateToFreq(t time.Time, freq string) time.Time {
	switch freq {
	case "T":
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	case "H":
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	case "D":
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case "W":
		// Truncate to the most recent Monday.
		wd := int(t.Weekday())
		if wd == 0 {
			wd = 7
		}
		d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		return d.AddDate(0, 0, -(wd - 1))
	case "M":
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	default:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}
}

// DateRange generates a DatetimeIndex from start to end (inclusive) at the given frequency.
// Supported frequencies: "D" (daily), "H" (hourly), "M" (monthly), "W" (weekly), "T" (minute).
func DateRange(start, end time.Time, freq string) *DatetimeIndex {
	var times []time.Time
	current := start
	for !current.After(end) {
		times = append(times, current)
		current = addFreq(current, 1, freq)
	}
	return &DatetimeIndex{times: times, name: "", freq: freq}
}

// Resample returns a new DatetimeIndex with timestamps resampled to the given frequency.
// Each original timestamp is truncated to its period boundary.
// Duplicate boundaries are kept (one per original timestamp).
func (di *DatetimeIndex) Resample(freq string) *DatetimeIndex {
	seen := make(map[time.Time]bool)
	var times []time.Time
	for _, t := range di.times {
		truncated := truncateToFreq(t, freq)
		if !seen[truncated] {
			seen[truncated] = true
			times = append(times, truncated)
		}
	}
	return &DatetimeIndex{times: times, name: di.name, freq: freq}
}

// Shift returns a new DatetimeIndex with each timestamp shifted by the given
// number of periods at the specified frequency.
func (di *DatetimeIndex) Shift(periods int, freq string) *DatetimeIndex {
	times := make([]time.Time, len(di.times))
	for i, t := range di.times {
		times[i] = addFreq(t, periods, freq)
	}
	return &DatetimeIndex{times: times, name: di.name, freq: di.freq}
}

// Asfreq returns a new DatetimeIndex converted to the given frequency.
// It generates new timestamps from the first to last element at the new frequency.
func (di *DatetimeIndex) Asfreq(freq string) *DatetimeIndex {
	if len(di.times) == 0 {
		return &DatetimeIndex{name: di.name, freq: freq}
	}
	start := di.times[0]
	end := di.times[len(di.times)-1]
	result := DateRange(start, end, freq)
	result.name = di.name
	return result
}

// BetweenTime returns the indices of timestamps whose time-of-day (hour, minute, second)
// falls between start and end (inclusive). Only the time portion is compared;
// the date portion of start and end is ignored.
func (di *DatetimeIndex) BetweenTime(start, end time.Time) []int {
	startSec := start.Hour()*3600 + start.Minute()*60 + start.Second()
	endSec := end.Hour()*3600 + end.Minute()*60 + end.Second()
	var indices []int
	for i, t := range di.times {
		tSec := t.Hour()*3600 + t.Minute()*60 + t.Second()
		if startSec <= endSec {
			if tSec >= startSec && tSec <= endSec {
				indices = append(indices, i)
			}
		} else {
			// Wraps around midnight.
			if tSec >= startSec || tSec <= endSec {
				indices = append(indices, i)
			}
		}
	}
	return indices
}

// AtTime returns the indices of timestamps that match the given time-of-day
// (hour, minute, second). The date portion of t is ignored.
func (di *DatetimeIndex) AtTime(t time.Time) []int {
	targetH, targetM, targetS := t.Hour(), t.Minute(), t.Second()
	var indices []int
	for i, ts := range di.times {
		if ts.Hour() == targetH && ts.Minute() == targetM && ts.Second() == targetS {
			indices = append(indices, i)
		}
	}
	return indices
}

// TZLocalize sets the timezone of each timestamp (assumes they are currently
// in an unspecified/UTC zone). Returns an error if the timezone name is invalid.
func (di *DatetimeIndex) TZLocalize(tz string) (*DatetimeIndex, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("tabgo: DatetimeIndex.TZLocalize: %w", err)
	}
	times := make([]time.Time, len(di.times))
	for i, t := range di.times {
		times[i] = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	}
	return &DatetimeIndex{times: times, name: di.name, freq: di.freq}, nil
}

// TZConvert converts each timestamp to the given timezone.
// Returns an error if the timezone name is invalid.
func (di *DatetimeIndex) TZConvert(tz string) (*DatetimeIndex, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("tabgo: DatetimeIndex.TZConvert: %w", err)
	}
	times := make([]time.Time, len(di.times))
	for i, t := range di.times {
		times[i] = t.In(loc)
	}
	return &DatetimeIndex{times: times, name: di.name, freq: di.freq}, nil
}

// -------------------------------------------------------------------
// Timedelta
// -------------------------------------------------------------------

// Timedelta wraps a time.Duration for pandas-like timedelta operations.
type Timedelta struct {
	duration time.Duration
}

// NewTimedelta creates a new Timedelta from a time.Duration.
func NewTimedelta(d time.Duration) Timedelta {
	return Timedelta{duration: d}
}

// Duration returns the underlying time.Duration.
func (td Timedelta) Duration() time.Duration {
	return td.duration
}

// Add returns a new Timedelta that is the sum of td and other.
func (td Timedelta) Add(other Timedelta) Timedelta {
	return Timedelta{duration: td.duration + other.duration}
}

// Sub returns a new Timedelta that is the difference of td and other.
func (td Timedelta) Sub(other Timedelta) Timedelta {
	return Timedelta{duration: td.duration - other.duration}
}

// Mul returns a new Timedelta scaled by n.
func (td Timedelta) Mul(n int) Timedelta {
	return Timedelta{duration: td.duration * time.Duration(n)}
}

// String returns a human-readable representation of the Timedelta.
func (td Timedelta) String() string {
	return td.duration.String()
}

// -------------------------------------------------------------------
// Period
// -------------------------------------------------------------------

// Period represents a time span at a given frequency, anchored to a start time.
type Period struct {
	start time.Time
	freq  string
}

// NewPeriod creates a new Period starting at start with the given frequency.
func NewPeriod(start time.Time, freq string) Period {
	return Period{start: start, freq: freq}
}

// Start returns the start time of the period.
func (p Period) Start() time.Time {
	return p.start
}

// End returns the end time of the period (exclusive).
// For monthly frequency this is the first day of the next month, etc.
func (p Period) End() time.Time {
	return addFreq(p.start, 1, p.freq)
}

// String returns a human-readable representation of the Period.
func (p Period) String() string {
	return fmt.Sprintf("Period(%s, %s)", p.start.Format(time.RFC3339), p.freq)
}
