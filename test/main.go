package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Date(2021, 6, 0, 23, 59, 59, 0, time.UTC)
	fmt.Println(t)

	t = time.Date(t.Year(), t.Month(), 0, 23, 59, 59, 0, time.UTC)
	fmt.Println(t)

	t = time.Date(t.Year(), t.Month(), 0, 23, 59, 59, 0, time.UTC)
	fmt.Println(t)
}
