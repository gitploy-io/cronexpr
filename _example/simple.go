package main

import (
	"fmt"
	"time"

	"github.com/gitploy-io/cronexpr"
)

func main() {
	t := time.Date(2013, 8, 29, 9, 28, 0, 0, time.UTC)

	nextTime := cronexpr.MustParse("0 0 29 2 *").Next(t)
	fmt.Print(nextTime)
}
