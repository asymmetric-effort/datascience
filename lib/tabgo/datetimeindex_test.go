//go:build unit

package tabgo

import (
	"testing"
	"time"
)

func TestDatetimeIndex_Basic(t *testing.T) {
	dates := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(dates, "dates")
	if di.Len() != 3 {
		t.Errorf("Len() = %d, want 3", di.Len())
	}
	if di.Name() != "dates" {
		t.Errorf("Name() = %q, want dates", di.Name())
	}
}

func TestDateRange_FiveDays(t *testing.T) {
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	di := DateRange(start, end, "D")
	if di.Len() != 5 {
		t.Errorf("DateRange periods=5 Len() = %d, want 5", di.Len())
	}
	if di.Get(0) != start {
		t.Errorf("DateRange[0] = %v, want %v", di.Get(0), start)
	}
	expected := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	if di.Get(4) != expected {
		t.Errorf("DateRange[4] = %v, want %v", di.Get(4), expected)
	}
}

func TestDateRange_StartEnd(t *testing.T) {
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)
	di := DateRange(start, end, "D")
	if di.Len() != 3 {
		t.Errorf("DateRange start-end Len() = %d, want 3", di.Len())
	}
}

func TestDateRange_Hourly(t *testing.T) {
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 3, 0, 0, 0, time.UTC)
	di := DateRange(start, end, "H")
	if di.Len() != 4 {
		t.Errorf("DateRange hourly Len() = %d, want 4", di.Len())
	}
	if di.Get(1).Hour() != 1 {
		t.Errorf("DateRange hourly[1].Hour() = %d, want 1", di.Get(1).Hour())
	}
}

func TestDateRange_Monthly(t *testing.T) {
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)
	di := DateRange(start, end, "M")
	if di.Len() != 3 {
		t.Errorf("Len() = %d, want 3", di.Len())
	}
	if int(di.Get(1).Month()) != 2 {
		t.Errorf("Monthly[1].Month() = %d, want 2", di.Get(1).Month())
	}
}

func TestDatetimeIndex_Year(t *testing.T) {
	dates := []time.Time{
		time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(dates, "d")
	years := di.Year()
	vals := years.Values()
	if vals[0] != 2023 || vals[1] != 2024 {
		t.Errorf("Year() = %v, want [2023 2024]", vals)
	}
}

func TestDatetimeIndex_Month(t *testing.T) {
	dates := []time.Time{
		time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(dates, "d")
	months := di.Month()
	vals := months.Values()
	if vals[0] != 6 {
		t.Errorf("Month() = %v, want [6]", vals)
	}
}

func TestDatetimeIndex_Day(t *testing.T) {
	dates := []time.Time{
		time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(dates, "d")
	days := di.Day()
	vals := days.Values()
	if vals[0] != 15 {
		t.Errorf("Day() = %v, want [15]", vals)
	}
}

func TestDatetimeIndex_Hour(t *testing.T) {
	dates := []time.Time{
		time.Date(2023, 6, 15, 14, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(dates, "d")
	hours := di.Hour()
	vals := hours.Values()
	if vals[0] != 14 {
		t.Errorf("Hour() = %v, want [14]", vals)
	}
}

func TestDatetimeIndex_Weekday(t *testing.T) {
	// 2023-01-01 is a Sunday
	dates := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), // Monday
	}
	di := NewDatetimeIndex(dates, "d")
	wd := di.Weekday()
	vals := wd.Values()
	if vals[0] != 0 || vals[1] != 1 {
		t.Errorf("Weekday() = %v, want [0 1]", vals)
	}
}

func TestDatetimeIndex_Slice(t *testing.T) {
	dates := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(dates, "d")
	sliced := di.Slice(1, 3)
	if sliced.Len() != 2 {
		t.Errorf("Slice(1,3).Len() = %d, want 2", sliced.Len())
	}
}

func TestDatetimeIndex_ToSeries(t *testing.T) {
	dates := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(dates, "d")
	s := di.ToSeries()
	if s.Len() != 1 {
		t.Errorf("ToSeries().Len() = %d, want 1", s.Len())
	}
	if s.Name() != "d" {
		t.Errorf("ToSeries().Name() = %q, want d", s.Name())
	}
}

func TestDatetimeIndex_String(t *testing.T) {
	dates := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(dates, "test")
	s := di.String()
	if len(s) == 0 {
		t.Error("String() should not be empty")
	}
}

func TestDatetimeIndex_Values(t *testing.T) {
	dates := []time.Time{
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(dates, "d")
	vals := di.Values()
	if len(vals) != 1 {
		t.Errorf("Values() length = %d, want 1", len(vals))
	}
	// Mutating should not affect original
	vals[0] = time.Now()
	if di.Get(0) != dates[0] {
		t.Error("Values() did not return a copy")
	}
}

func TestDatetimeIndex_BetweenTime_Legacy(t *testing.T) {
	dates := []time.Time{
		time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 1, 18, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 1, 22, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(dates, "d")
	start := time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)
	end := time.Date(0, 1, 1, 19, 0, 0, 0, time.UTC)
	indices := di.BetweenTime(start, end)
	if len(indices) != 2 {
		t.Fatalf("BetweenTime length = %d, want 2", len(indices))
	}
	if indices[0] != 1 || indices[1] != 2 {
		t.Errorf("BetweenTime = %v, want [1 2]", indices)
	}
}
