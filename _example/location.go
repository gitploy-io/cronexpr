package main

import (
	"fmt"
	"time"

	"github.com/gitploy-io/cronexpr"
)

func main() {
	nextTime := cronexpr.MustParseInLocation("0 * * * *", "Asia/Seoul").Next(time.Now())
	fmt.Printf("Parse the cron expression in the KR timezone: %s", nextTime)
}
