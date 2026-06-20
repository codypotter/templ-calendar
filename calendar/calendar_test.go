package calendar

import (
	"testing"
	"time"
)

func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		year  int
		month time.Month
		want  int
	}{
		{2024, time.January, 31},
		{2024, time.February, 29}, // leap year
		{2023, time.February, 28}, // non-leap year
		{2024, time.April, 30},
		{2024, time.December, 31},
	}
	for _, tt := range tests {
		got := daysInMonth(tt.year, tt.month)
		if got != tt.want {
			t.Errorf("daysInMonth(%d, %s) = %d, want %d", tt.year, tt.month, got, tt.want)
		}
	}
}

func TestFirstWeekday(t *testing.T) {
	tests := []struct {
		year  int
		month time.Month
		want  time.Weekday
	}{
		{2024, time.January, time.Monday},   // Jan 1 2024 = Monday
		{2024, time.February, time.Thursday}, // Jan has 31 days (31%7=3), Mon+3=Thu
		{2024, time.June, time.Saturday},    // verified by counting forward from Jan
		{2025, time.January, time.Wednesday}, // 2024 leap year (366 days, 366%7=2), Mon+2=Wed
	}
	for _, tt := range tests {
		got := firstWeekday(tt.year, tt.month)
		if got != tt.want {
			t.Errorf("firstWeekday(%d, %s) = %s, want %s", tt.year, tt.month, got, tt.want)
		}
	}
}

func TestTrailingDays(t *testing.T) {
	tests := []struct {
		year  int
		month time.Month
		want  int
	}{
		// Jan 2024: Mon start (offset 1), 31 days → total 32 → 5 rows → 35 cells → 3 trailing
		{2024, time.January, 3},
		// Feb 2024: Thu start (offset 4), 29 days → total 33 → 5 rows → 35 cells → 2 trailing
		{2024, time.February, 2},
		// Jun 2024: Sat start (offset 6), 30 days → total 36 → 6 rows → 42 cells → 6 trailing
		{2024, time.June, 6},
	}
	for _, tt := range tests {
		got := trailingDays(tt.year, tt.month)
		if got != tt.want {
			t.Errorf("trailingDays(%d, %s) = %d, want %d", tt.year, tt.month, got, tt.want)
		}
	}
}

func TestEventsForDay(t *testing.T) {
	loc := time.UTC

	single := Event{
		ID:    "single",
		Title: "Standup",
		Start: time.Date(2026, time.June, 10, 9, 0, 0, 0, loc),
		End:   time.Date(2026, time.June, 10, 9, 30, 0, 0, loc),
	}
	multi := Event{
		ID:    "multi",
		Title: "Conference",
		Start: time.Date(2026, time.June, 15, 8, 0, 0, 0, loc),
		End:   time.Date(2026, time.June, 17, 18, 0, 0, 0, loc),
	}
	// all-day style: midnight to midnight (exclusive end)
	allDay := Event{
		ID:    "allday",
		Title: "Holiday",
		Start: time.Date(2026, time.June, 20, 0, 0, 0, 0, loc),
		End:   time.Date(2026, time.June, 21, 0, 0, 0, 0, loc),
	}

	events := []Event{single, multi, allDay}

	tests := []struct {
		day      int
		wantIDs  []string
	}{
		{9, nil},                        // day before single
		{10, []string{"single"}},        // single day event
		{11, nil},                       // day after single
		{14, nil},                       // day before multi
		{15, []string{"multi"}},         // first day of multi
		{16, []string{"multi"}},         // middle of multi
		{17, []string{"multi"}},         // last day (end is 18:00, still within day)
		{18, nil},                       // day after multi ends
		{20, []string{"allday"}},        // all-day event
		{21, nil},                       // exclusive end: midnight end does not bleed into next day
	}

	for _, tt := range tests {
		got := eventsForDay(events, 2026, time.June, tt.day, loc)
		if len(got) != len(tt.wantIDs) {
			t.Errorf("day %d: got %d events, want %d", tt.day, len(got), len(tt.wantIDs))
			continue
		}
		for i, id := range tt.wantIDs {
			if got[i].ID != id {
				t.Errorf("day %d: event[%d].ID = %q, want %q", tt.day, i, got[i].ID, id)
			}
		}
	}
}

func TestEventsForDayTimezone(t *testing.T) {
	eastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("timezone data unavailable")
	}

	// Event at 11pm Eastern = 4am UTC next day.
	// With Eastern loc it should appear on Jun 10; with UTC it should not.
	e := Event{
		ID:    "tz",
		Title: "Late night",
		Start: time.Date(2026, time.June, 10, 23, 0, 0, 0, eastern),
		End:   time.Date(2026, time.June, 10, 23, 59, 0, 0, eastern),
	}

	got := eventsForDay([]Event{e}, 2026, time.June, 10, eastern)
	if len(got) != 1 {
		t.Errorf("eastern loc, day 10: got %d events, want 1", len(got))
	}

	got = eventsForDay([]Event{e}, 2026, time.June, 10, time.UTC)
	if len(got) != 0 {
		t.Errorf("UTC loc, day 10: got %d events, want 0 (event is on Jun 11 in UTC)", len(got))
	}

	got = eventsForDay([]Event{e}, 2026, time.June, 11, time.UTC)
	if len(got) != 1 {
		t.Errorf("UTC loc, day 11: got %d events, want 1", len(got))
	}
}
