package tabgo

import (
	"time"
)

// DtAccessor provides vectorized datetime operations on a Series.
// Values are expected to be time.Time; non-time values yield nil.
type DtAccessor struct {
	series *Series
}

// Dt returns a DtAccessor for vectorized datetime operations.
func (s *Series) Dt() *DtAccessor {
	return &DtAccessor{series: s}
}

// toTime attempts to convert a value to time.Time.
// Returns the zero time and false if conversion fails.
func toTime(v any) (time.Time, bool) {
	switch t := v.(type) {
	case time.Time:
		return t, true
	case *time.Time:
		if t == nil {
			return time.Time{}, false
		}
		return *t, true
	case string:
		// Try common formats.
		for _, layout := range []string{
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
			"01/02/2006",
			"Jan 2, 2006",
		} {
			if parsed, err := time.Parse(layout, t); err == nil {
				return parsed, true
			}
		}
		return time.Time{}, false
	default:
		return time.Time{}, false
	}
}

// Year returns a new Series with the year of each datetime element.
func (da *DtAccessor) Year() *Series {
	out := make([]any, len(da.series.values))
	for i, v := range da.series.values {
		if t, ok := toTime(v); ok {
			out[i] = t.Year()
		} else {
			out[i] = nil
		}
	}
	return &Series{name: da.series.name, values: out}
}

// Month returns a new Series with the month (1-12) of each datetime element.
func (da *DtAccessor) Month() *Series {
	out := make([]any, len(da.series.values))
	for i, v := range da.series.values {
		if t, ok := toTime(v); ok {
			out[i] = int(t.Month())
		} else {
			out[i] = nil
		}
	}
	return &Series{name: da.series.name, values: out}
}

// Day returns a new Series with the day of month of each datetime element.
func (da *DtAccessor) Day() *Series {
	out := make([]any, len(da.series.values))
	for i, v := range da.series.values {
		if t, ok := toTime(v); ok {
			out[i] = t.Day()
		} else {
			out[i] = nil
		}
	}
	return &Series{name: da.series.name, values: out}
}

// Hour returns a new Series with the hour of each datetime element.
func (da *DtAccessor) Hour() *Series {
	out := make([]any, len(da.series.values))
	for i, v := range da.series.values {
		if t, ok := toTime(v); ok {
			out[i] = t.Hour()
		} else {
			out[i] = nil
		}
	}
	return &Series{name: da.series.name, values: out}
}

// Minute returns a new Series with the minute of each datetime element.
func (da *DtAccessor) Minute() *Series {
	out := make([]any, len(da.series.values))
	for i, v := range da.series.values {
		if t, ok := toTime(v); ok {
			out[i] = t.Minute()
		} else {
			out[i] = nil
		}
	}
	return &Series{name: da.series.name, values: out}
}

// Second returns a new Series with the second of each datetime element.
func (da *DtAccessor) Second() *Series {
	out := make([]any, len(da.series.values))
	for i, v := range da.series.values {
		if t, ok := toTime(v); ok {
			out[i] = t.Second()
		} else {
			out[i] = nil
		}
	}
	return &Series{name: da.series.name, values: out}
}

// Weekday returns a new Series with the day of the week (0=Sunday, 6=Saturday).
func (da *DtAccessor) Weekday() *Series {
	out := make([]any, len(da.series.values))
	for i, v := range da.series.values {
		if t, ok := toTime(v); ok {
			out[i] = int(t.Weekday())
		} else {
			out[i] = nil
		}
	}
	return &Series{name: da.series.name, values: out}
}

// Date returns a new Series with each datetime truncated to date (midnight).
func (da *DtAccessor) Date() *Series {
	out := make([]any, len(da.series.values))
	for i, v := range da.series.values {
		if t, ok := toTime(v); ok {
			out[i] = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		} else {
			out[i] = nil
		}
	}
	return &Series{name: da.series.name, values: out}
}

// DayOfYear returns a new Series with the day of the year (1-366).
func (da *DtAccessor) DayOfYear() *Series {
	out := make([]any, len(da.series.values))
	for i, v := range da.series.values {
		if t, ok := toTime(v); ok {
			out[i] = t.YearDay()
		} else {
			out[i] = nil
		}
	}
	return &Series{name: da.series.name, values: out}
}

// Quarter returns a new Series with the quarter (1-4) of each datetime element.
func (da *DtAccessor) Quarter() *Series {
	out := make([]any, len(da.series.values))
	for i, v := range da.series.values {
		if t, ok := toTime(v); ok {
			out[i] = (int(t.Month())-1)/3 + 1
		} else {
			out[i] = nil
		}
	}
	return &Series{name: da.series.name, values: out}
}
