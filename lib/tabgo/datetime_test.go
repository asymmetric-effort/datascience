//go:build unit

package tabgo

import (
	"testing"
	"time"
)

func TestNewDatetimeIndex(t *testing.T) {
	times := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(times, "date")

	if di.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", di.Len())
	}
	if di.Name() != "date" {
		t.Fatalf("Name() = %q, want %q", di.Name(), "date")
	}

	// Mutating original should not affect index.
	times[0] = time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)
	if di.Times()[0].Year() == 2099 {
		t.Fatal("NewDatetimeIndex did not copy input")
	}
}

func TestDatetimeIndexTimes(t *testing.T) {
	times := []time.Time{time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)}
	di := NewDatetimeIndex(times, "ts")
	got := di.Times()
	// Mutating returned slice should not affect index.
	got[0] = time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)
	if di.Times()[0].Year() == 2099 {
		t.Fatal("Times() did not return a copy")
	}
}

func TestDateRangeDaily(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	di := DateRange(start, end, "D")

	if di.Len() != 5 {
		t.Fatalf("DateRange D len = %d, want 5", di.Len())
	}
	if di.Freq() != "D" {
		t.Fatalf("Freq() = %q, want D", di.Freq())
	}
	if !di.Times()[0].Equal(start) {
		t.Fatalf("first time = %v, want %v", di.Times()[0], start)
	}
	if !di.Times()[4].Equal(end) {
		t.Fatalf("last time = %v, want %v", di.Times()[4], end)
	}
}

func TestDateRangeHourly(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 3, 0, 0, 0, time.UTC)
	di := DateRange(start, end, "H")

	if di.Len() != 4 {
		t.Fatalf("DateRange H len = %d, want 4", di.Len())
	}
}

func TestDateRangeMonthly(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	di := DateRange(start, end, "M")

	if di.Len() != 6 {
		t.Fatalf("DateRange M len = %d, want 6", di.Len())
	}
	// Check months are correct.
	for i, ts := range di.Times() {
		wantMonth := time.Month(i + 1)
		if ts.Month() != wantMonth {
			t.Fatalf("month[%d] = %v, want %v", i, ts.Month(), wantMonth)
		}
	}
}

func TestDateRangeWeekly(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC)
	di := DateRange(start, end, "W")

	if di.Len() != 5 {
		t.Fatalf("DateRange W len = %d, want 5", di.Len())
	}
}

func TestDateRangeMinute(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 0, 5, 0, 0, time.UTC)
	di := DateRange(start, end, "T")

	if di.Len() != 6 {
		t.Fatalf("DateRange T len = %d, want 6", di.Len())
	}
}

func TestDateRangeEmpty(t *testing.T) {
	start := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	di := DateRange(start, end, "D")

	if di.Len() != 0 {
		t.Fatalf("DateRange with end < start should be empty, got len = %d", di.Len())
	}
}

func TestResample(t *testing.T) {
	// Hourly data resampled to daily.
	times := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 6, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 6, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(times, "ts")
	resampled := di.Resample("D")

	if resampled.Len() != 2 {
		t.Fatalf("Resample D len = %d, want 2", resampled.Len())
	}
	if resampled.Freq() != "D" {
		t.Fatalf("Resample Freq() = %q, want D", resampled.Freq())
	}
}

func TestResampleMonthly(t *testing.T) {
	times := []time.Time{
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(times, "ts")
	resampled := di.Resample("M")

	if resampled.Len() != 3 {
		t.Fatalf("Resample M len = %d, want 3", resampled.Len())
	}
}

func TestShift(t *testing.T) {
	times := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(times, "ts")

	// Shift forward 3 days.
	shifted := di.Shift(3, "D")
	if shifted.Len() != 2 {
		t.Fatalf("Shift len = %d, want 2", shifted.Len())
	}
	expected := time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)
	if !shifted.Times()[0].Equal(expected) {
		t.Fatalf("Shift(3, D)[0] = %v, want %v", shifted.Times()[0], expected)
	}

	// Shift backward.
	shiftedBack := di.Shift(-1, "D")
	expectedBack := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	if !shiftedBack.Times()[0].Equal(expectedBack) {
		t.Fatalf("Shift(-1, D)[0] = %v, want %v", shiftedBack.Times()[0], expectedBack)
	}
}

func TestShiftHourly(t *testing.T) {
	times := []time.Time{time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)}
	di := NewDatetimeIndex(times, "ts")
	shifted := di.Shift(5, "H")
	expected := time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC)
	if !shifted.Times()[0].Equal(expected) {
		t.Fatalf("Shift(5, H) = %v, want %v", shifted.Times()[0], expected)
	}
}

func TestShiftMonthly(t *testing.T) {
	times := []time.Time{time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)}
	di := NewDatetimeIndex(times, "ts")
	shifted := di.Shift(2, "M")
	expected := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if !shifted.Times()[0].Equal(expected) {
		t.Fatalf("Shift(2, M) = %v, want %v", shifted.Times()[0], expected)
	}
}

func TestAsfreq(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
	di := DateRange(start, end, "D")

	// Convert daily to hourly.
	hourly := di.Asfreq("H")
	// 2 full days = 48 hours + hour 0 of day 3 = 49 entries.
	if hourly.Len() != 49 {
		t.Fatalf("Asfreq H len = %d, want 49", hourly.Len())
	}
}

func TestAsfreqEmpty(t *testing.T) {
	di := NewDatetimeIndex(nil, "empty")
	result := di.Asfreq("D")
	if result.Len() != 0 {
		t.Fatalf("Asfreq on empty should return empty, got len = %d", result.Len())
	}
}

func TestBetweenTime(t *testing.T) {
	times := []time.Time{
		time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 22, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(times, "ts")

	start := time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)
	end := time.Date(0, 1, 1, 18, 0, 0, 0, time.UTC)
	indices := di.BetweenTime(start, end)

	// Should match: index 1 (12:00), index 2 (18:00), index 3 (9:00).
	if len(indices) != 3 {
		t.Fatalf("BetweenTime len = %d, want 3, got indices: %v", len(indices), indices)
	}
	expected := []int{1, 2, 3}
	for i, idx := range indices {
		if idx != expected[i] {
			t.Fatalf("BetweenTime[%d] = %d, want %d", i, idx, expected[i])
		}
	}
}

func TestBetweenTimeWrapAround(t *testing.T) {
	// Test wrap-around (e.g., 22:00 to 02:00).
	times := []time.Time{
		time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 23, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(times, "ts")

	start := time.Date(0, 1, 1, 22, 0, 0, 0, time.UTC)
	end := time.Date(0, 1, 1, 2, 0, 0, 0, time.UTC)
	indices := di.BetweenTime(start, end)

	// Should match: index 0 (01:00) and index 2 (23:00).
	if len(indices) != 2 {
		t.Fatalf("BetweenTime wrap len = %d, want 2, indices: %v", len(indices), indices)
	}
}

func TestAtTime(t *testing.T) {
	times := []time.Time{
		time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 9, 30, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 15, 0, 0, 0, time.UTC),
	}
	di := NewDatetimeIndex(times, "ts")

	target := time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC)
	indices := di.AtTime(target)

	if len(indices) != 2 {
		t.Fatalf("AtTime len = %d, want 2", len(indices))
	}
	if indices[0] != 0 || indices[1] != 2 {
		t.Fatalf("AtTime indices = %v, want [0 2]", indices)
	}
}

func TestAtTimeNoMatch(t *testing.T) {
	times := []time.Time{time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)}
	di := NewDatetimeIndex(times, "ts")

	target := time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC)
	indices := di.AtTime(target)
	if len(indices) != 0 {
		t.Fatalf("AtTime should return empty for no match, got %v", indices)
	}
}

func TestTZLocalize(t *testing.T) {
	times := []time.Time{time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)}
	di := NewDatetimeIndex(times, "ts")

	localized, err := di.TZLocalize("America/New_York")
	if err != nil {
		t.Fatalf("TZLocalize error: %v", err)
	}
	loc := localized.Times()[0].Location()
	if loc.String() != "America/New_York" {
		t.Fatalf("TZLocalize location = %q, want America/New_York", loc.String())
	}
	// The wall clock time should be preserved.
	if localized.Times()[0].Hour() != 12 {
		t.Fatalf("TZLocalize hour = %d, want 12", localized.Times()[0].Hour())
	}
}

func TestTZLocalizeInvalidTZ(t *testing.T) {
	di := NewDatetimeIndex(nil, "ts")
	_, err := di.TZLocalize("Invalid/Zone")
	if err == nil {
		t.Fatal("expected error for invalid timezone")
	}
}

func TestTZConvert(t *testing.T) {
	// Create a time in UTC.
	utcTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	di := NewDatetimeIndex([]time.Time{utcTime}, "ts")

	converted, err := di.TZConvert("America/New_York")
	if err != nil {
		t.Fatalf("TZConvert error: %v", err)
	}
	ct := converted.Times()[0]
	if ct.Location().String() != "America/New_York" {
		t.Fatalf("TZConvert location = %q, want America/New_York", ct.Location().String())
	}
	// UTC 12:00 should be ET 08:00 in June (EDT).
	if ct.Hour() != 8 {
		t.Fatalf("TZConvert hour = %d, want 8 (EDT)", ct.Hour())
	}
}

func TestTZConvertInvalidTZ(t *testing.T) {
	di := NewDatetimeIndex(nil, "ts")
	_, err := di.TZConvert("Invalid/Zone")
	if err == nil {
		t.Fatal("expected error for invalid timezone")
	}
}

// --- Timedelta tests ---

func TestNewTimedelta(t *testing.T) {
	td := NewTimedelta(5 * time.Hour)
	if td.Duration() != 5*time.Hour {
		t.Fatalf("Duration() = %v, want 5h", td.Duration())
	}
}

func TestTimedeltaAdd(t *testing.T) {
	a := NewTimedelta(2 * time.Hour)
	b := NewTimedelta(30 * time.Minute)
	result := a.Add(b)
	expected := 2*time.Hour + 30*time.Minute
	if result.Duration() != expected {
		t.Fatalf("Add = %v, want %v", result.Duration(), expected)
	}
}

func TestTimedeltaSub(t *testing.T) {
	a := NewTimedelta(5 * time.Hour)
	b := NewTimedelta(2 * time.Hour)
	result := a.Sub(b)
	if result.Duration() != 3*time.Hour {
		t.Fatalf("Sub = %v, want 3h", result.Duration())
	}
}

func TestTimedeltaMul(t *testing.T) {
	td := NewTimedelta(10 * time.Minute)
	result := td.Mul(3)
	if result.Duration() != 30*time.Minute {
		t.Fatalf("Mul = %v, want 30m", result.Duration())
	}
}

func TestTimedeltaString(t *testing.T) {
	td := NewTimedelta(90 * time.Minute)
	s := td.String()
	if s != "1h30m0s" {
		t.Fatalf("String() = %q, want 1h30m0s", s)
	}
}

func TestTimedeltaZero(t *testing.T) {
	td := NewTimedelta(0)
	if td.Duration() != 0 {
		t.Fatalf("zero timedelta Duration = %v", td.Duration())
	}
	if td.String() != "0s" {
		t.Fatalf("zero timedelta String = %q", td.String())
	}
}

func TestTimedeltaNegative(t *testing.T) {
	a := NewTimedelta(time.Hour)
	b := NewTimedelta(2 * time.Hour)
	result := a.Sub(b)
	if result.Duration() != -time.Hour {
		t.Fatalf("negative Sub = %v, want -1h", result.Duration())
	}
}

// --- Period tests ---

func TestNewPeriod(t *testing.T) {
	start := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, "M")

	if !p.Start().Equal(start) {
		t.Fatalf("Start() = %v, want %v", p.Start(), start)
	}
}

func TestPeriodEndMonthly(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, "M")

	expected := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	if !p.End().Equal(expected) {
		t.Fatalf("End() = %v, want %v", p.End(), expected)
	}
}

func TestPeriodEndDaily(t *testing.T) {
	start := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, "D")

	expected := time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC)
	if !p.End().Equal(expected) {
		t.Fatalf("End() = %v, want %v", p.End(), expected)
	}
}

func TestPeriodEndHourly(t *testing.T) {
	start := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	p := NewPeriod(start, "H")

	expected := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	if !p.End().Equal(expected) {
		t.Fatalf("End() = %v, want %v", p.End(), expected)
	}
}

func TestPeriodEndWeekly(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Monday
	p := NewPeriod(start, "W")

	expected := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	if !p.End().Equal(expected) {
		t.Fatalf("End() = %v, want %v", p.End(), expected)
	}
}

func TestPeriodString(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	p := NewPeriod(start, "D")
	s := p.String()
	if s == "" {
		t.Fatal("String() should not be empty")
	}
	if len(s) < 10 {
		t.Fatalf("String() too short: %q", s)
	}
}

func TestPeriodEndMinute(t *testing.T) {
	start := time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC)
	p := NewPeriod(start, "T")

	expected := time.Date(2024, 1, 1, 10, 31, 0, 0, time.UTC)
	if !p.End().Equal(expected) {
		t.Fatalf("End() = %v, want %v", p.End(), expected)
	}
}

// --- Integration-style tests ---

func TestDateRangeResampleRoundTrip(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 23, 0, 0, 0, time.UTC)

	hourly := DateRange(start, end, "H")
	if hourly.Len() != 24 {
		t.Fatalf("hourly len = %d, want 24", hourly.Len())
	}

	daily := hourly.Resample("D")
	if daily.Len() != 1 {
		t.Fatalf("daily resample len = %d, want 1", daily.Len())
	}
}

func TestShiftPreservesName(t *testing.T) {
	di := NewDatetimeIndex(
		[]time.Time{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		"myindex",
	)
	shifted := di.Shift(1, "D")
	if shifted.Name() != "myindex" {
		t.Fatalf("Shift changed name to %q", shifted.Name())
	}
}
