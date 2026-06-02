//go:build unit

package tabgo

import (
	"testing"
	"time"
)

func TestDtYear(t *testing.T) {
	dates := []any{
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC),
		nil,
	}
	s := NewSeries("dates", dates)
	years := s.Dt().Year()
	vals := years.Values()
	if vals[0] != 2023 || vals[1] != 2024 || vals[2] != nil {
		t.Errorf("Year() = %v", vals)
	}
}

func TestDtMonth(t *testing.T) {
	dates := []any{
		time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 20, 0, 0, 0, 0, time.UTC),
	}
	s := NewSeries("dates", dates)
	months := s.Dt().Month()
	vals := months.Values()
	if vals[0] != 3 || vals[1] != 12 {
		t.Errorf("Month() = %v", vals)
	}
}

func TestDtDay(t *testing.T) {
	dates := []any{
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
	}
	s := NewSeries("dates", dates)
	days := s.Dt().Day()
	vals := days.Values()
	if vals[0] != 15 || vals[1] != 1 {
		t.Errorf("Day() = %v", vals)
	}
}

func TestDtHour(t *testing.T) {
	dates := []any{
		time.Date(2023, 1, 15, 14, 30, 0, 0, time.UTC),
		time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC),
	}
	s := NewSeries("dates", dates)
	hours := s.Dt().Hour()
	vals := hours.Values()
	if vals[0] != 14 || vals[1] != 0 {
		t.Errorf("Hour() = %v", vals)
	}
}

func TestDtMinute(t *testing.T) {
	dates := []any{
		time.Date(2023, 1, 15, 14, 30, 0, 0, time.UTC),
	}
	s := NewSeries("dates", dates)
	mins := s.Dt().Minute()
	vals := mins.Values()
	if vals[0] != 30 {
		t.Errorf("Minute() = %v", vals)
	}
}

func TestDtSecond(t *testing.T) {
	dates := []any{
		time.Date(2023, 1, 15, 14, 30, 45, 0, time.UTC),
	}
	s := NewSeries("dates", dates)
	secs := s.Dt().Second()
	vals := secs.Values()
	if vals[0] != 45 {
		t.Errorf("Second() = %v", vals)
	}
}

func TestDtWeekday(t *testing.T) {
	// 2023-01-15 is a Sunday (weekday 0)
	dates := []any{
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 16, 0, 0, 0, 0, time.UTC), // Monday (1)
	}
	s := NewSeries("dates", dates)
	wd := s.Dt().Weekday()
	vals := wd.Values()
	if vals[0] != 0 || vals[1] != 1 {
		t.Errorf("Weekday() = %v", vals)
	}
}

func TestDtDate(t *testing.T) {
	dates := []any{
		time.Date(2023, 1, 15, 14, 30, 45, 0, time.UTC),
	}
	s := NewSeries("dates", dates)
	d := s.Dt().Date()
	vals := d.Values()
	expected := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
	if v, ok := vals[0].(time.Time); !ok || !v.Equal(expected) {
		t.Errorf("Date() = %v, want %v", vals[0], expected)
	}
}

func TestDtQuarter(t *testing.T) {
	dates := []any{
		time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),  // Q1
		time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC),  // Q2
		time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC),  // Q3
		time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC), // Q4
	}
	s := NewSeries("dates", dates)
	q := s.Dt().Quarter()
	vals := q.Values()
	if vals[0] != 1 || vals[1] != 2 || vals[2] != 3 || vals[3] != 4 {
		t.Errorf("Quarter() = %v, want [1 2 3 4]", vals)
	}
}

func TestDtDayOfYear(t *testing.T) {
	dates := []any{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), // day 1
		time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC), // day 32
	}
	s := NewSeries("dates", dates)
	doy := s.Dt().DayOfYear()
	vals := doy.Values()
	if vals[0] != 1 || vals[1] != 32 {
		t.Errorf("DayOfYear() = %v, want [1 32]", vals)
	}
}

func TestDtStringParsing(t *testing.T) {
	dates := []any{
		"2023-01-15",
		"2024-06-20T14:30:00",
	}
	s := NewSeries("dates", dates)
	years := s.Dt().Year()
	vals := years.Values()
	if vals[0] != 2023 || vals[1] != 2024 {
		t.Errorf("Year() from strings = %v, want [2023 2024]", vals)
	}
}
