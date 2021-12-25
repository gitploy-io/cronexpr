package cronexpr

import (
	"time"
)

type (
	Schedule struct {
		Minute, Hour, Dom, Month, Dow bitset

		Location *time.Location
	}
)

func (s *Schedule) Next(t time.Time) time.Time {
	loc := time.UTC
	if s.Location != nil {
		loc = s.Location
	}

	t.In(loc)

	added := false

	// Start at the earliest possible time (the upcoming second).
	t = t.Add(1*time.Minute - time.Duration(t.Nanosecond())*time.Nanosecond)

	yearLimit := t.Year() + 5

L:
	if t.Year() > yearLimit {
		return time.Time{}
	}

	// Find the first applicable month.
	// If it's this month, then do nothing.
	year := t.Year()
	for 1<<uint(t.Month())&s.Month == 0 {
		// If we have to add a month, reset the other parts to 0.
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
		}

		t = t.AddDate(0, 1, 0)

		if t.Year() != year {
			goto L
		}
	}

	// Now get a day in that month.
	month := t.Month()
	for !dayMatches(s, t) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
		}

		t = t.AddDate(0, 0, 1)

		if t.Month() != month {
			goto L
		}
	}

	day := t.Day()
	for 1<<uint(t.Hour())&s.Hour == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, loc)
		}

		t = t.Add(1 * time.Hour)

		if t.Day() != day {
			goto L
		}
	}

	hour := t.Hour()
	for 1<<uint(t.Minute())&s.Minute == 0 {
		if !added {
			added = true
			t = t.Truncate(time.Minute)
		}

		t = t.Add(1 * time.Minute)

		if t.Hour() != hour {
			goto L
		}
	}

	return t
}

func (s *Schedule) Prev(t time.Time) time.Time {
	loc := time.UTC
	if s.Location != nil {
		loc = s.Location
	}

	t.In(loc)

	subtracted := false

	// Start at the earliest possible time (the upcoming second).
	t = t.Add(-1*time.Minute + time.Duration(t.Nanosecond())*time.Nanosecond)

	yearLimit := t.Year() - 5

L:
	if t.Year() < yearLimit {
		return time.Time{}
	}

	year := t.Year()
	for 1<<uint(t.Month())&s.Month == 0 {
		// If we have to add a month, reset with the next month before.
		if !subtracted {
			subtracted = true
			t = time.Date(t.Year(), t.Month()+1, 0, 23, 59, 0, 0, loc)
		}

		// Change the time into the last day of the previous month.
		// Note that AddDate(0, -1, 0) has a bug by the normalization.
		// E.g) time.Date(2021, 6, 0, 23, 59, 59, 0, time.UTC).AddDate(0, -1, 0)
		t = time.Date(t.Year(), t.Month(), 0, 23, 59, 0, 0, loc)

		if t.Year() != year {
			goto L
		}
	}

	// Now get a day in that month.
	month := t.Month()
	for !dayMatches(s, t) {
		if !subtracted {
			subtracted = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 0, 0, loc)
		}

		t = t.AddDate(0, 0, -1)

		if t.Month() != month {
			goto L
		}
	}

	day := t.Day()
	for 1<<uint(t.Hour())&s.Hour == 0 {
		if !subtracted {
			subtracted = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 0, 0, loc)
		}

		t = t.Add(-1 * time.Hour)

		if t.Day() != day {
			goto L
		}
	}

	hour := t.Hour()
	for 1<<uint(t.Minute())&s.Minute == 0 {
		if !subtracted {
			subtracted = true
			t = t.Truncate(-time.Minute)
		}

		t = t.Add(-1 * time.Minute)

		if t.Hour() != hour {
			goto L
		}
	}

	return t
}

// dayMatches returns true if the schedule's day-of-week and day-of-month
// restrictions are satisfied by the given time.
func dayMatches(s *Schedule, t time.Time) bool {
	var (
		domMatch bool = 1<<uint(t.Day())&s.Dom > 0
		dowMatch bool = 1<<uint(t.Weekday())&s.Dow > 0
	)
	if s.Dom&bitsetStar > 0 || s.Dow&bitsetStar > 0 {
		return domMatch && dowMatch
	}
	return domMatch || dowMatch
}
