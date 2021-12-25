package main

import (
	"log"
	"time"

	"github.com/gitploy-io/cronexpr"
)

// verifyFreezeWindow verifies the time satisfy being between two cronexpr.
func verifyFreezeWindow(begin, end string, t time.Time) (bool, error) {
	b, err := cronexpr.Parse(begin)
	if err != nil {
		return false, err
	}

	e, err := cronexpr.Parse(end)
	if err != nil {
		return false, err
	}

	if !(t.Before(b.Prev(t)) && t.After(e.Next(t))) {
		return false, nil
	}

	return true, nil
}

func main() {
	t := time.Date(2021, 12, 25, 0, 0, 0, 0, time.UTC)

	ok, _ := verifyFreezeWindow("* 23 24 DEC *", "* 1 25 DEC *", t)
	if !ok {
		log.Printf("Blocked to deploy at the midnight of X-mas.")
	}
}
