package cronexpr

import (
	"testing"
	"time"
)

func getTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}

	var layouts = []string{
		"Mon Jan 2 15:04 2006",
		"Mon Jan 2 15:04:05 2006",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, time.UTC); err == nil {
			return t
		}
	}
	if t, err := time.ParseInLocation("2006-01-02T15:04:05-0700", value, time.UTC); err == nil {
		return t
	}
	panic("could not parse time value " + value)
}

func TestSchedule_Next(t *testing.T) {
	runs := []struct {
		time, spec string
		expected   string
	}{
		// Simple
		{"Mon Jul 9 14:45 2012", "0/15 * * * *", "Mon Jul 9 15:00 2012"},
		{"Mon Jul 9 14:59 2012", "0/15 * * * *", "Mon Jul 9 15:00 2012"},

		// Wrap around hours
		{"Mon Jul 9 15:45 2012", "20-35/15 * * * *", "Mon Jul 9 16:20 2012"},

		// Wrap around days
		{"Mon Jul 9 23:46 2012", "0/15 * * * *", "Tue Jul 10 00:00 2012"},
		{"Mon Jul 9 23:45 2012", "20-35/15 * * * *", "Tue Jul 10 00:20 2012"},

		// Wrap around months
		{"Mon Jul 9 23:35 2012", "0 0 9 Apr-Oct *", "Thu Aug 9 00:00 2012"},
		{"Mon Jul 9 23:35 2012", "0 0 * Apr,Aug,Oct *", "Tue Aug 1 00:00 2012"},
	}

	for _, c := range runs {
		sched := MustParse(c.spec)
		actual := sched.Next(getTime(c.time))
		expected := getTime(c.expected)
		if !actual.Equal(expected) {
			t.Errorf("Next(%s) =  %v, wanted %v", c.time, actual, expected)
		}
	}
}

func TestSchedule_Prev(t *testing.T) {
	runs := []struct {
		time, spec string
		expected   string
	}{
		// Simple
		{"Mon Jul 9 14:45 2012", "0/15 * * * *", "Mon Jul 9 14:30 2012"},
		{"Mon Jul 9 14:59 2012", "0/15 * * * *", "Mon Jul 9 14:45 2012"},

		// Wrap around hours
		{"Mon Jul 9 15:45 2012", "20-35/15 * * * *", "Mon Jul 9 15:35 2012"},

		// Wrap around days
		{"Tue Jul 10 00:00 2012", "0/15 * * * *", "Mon Jul 9 23:45 2012"},
		{"Tue Jul 10 00:20 2012", "20-35/15 * * * *", "Mon Jul 9 23:35 2012"},

		// Wrap around months
		{"Thu Aug 9 00:00 2012", "0 0 9 Apr-Oct *", "Mon Jul 9 00:00 2012"},
		{"Tue Aug 1 00:00 2012", "0 0 * Apr,Aug,Oct *", "Mon Apr 30 00:00 2012"},
	}

	for _, c := range runs {
		sched := MustParse(c.spec)
		actual := sched.Prev(getTime(c.time))
		expected := getTime(c.expected)
		if !actual.Equal(expected) {
			t.Errorf("Prev(%s) =  %v, wanted %v", c.time, actual, expected)
		}
	}
}
