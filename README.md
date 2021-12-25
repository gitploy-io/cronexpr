# Cron expression parser

**Given a cron expression, you can get the previous timestamp or the next timestamp that satisfies the cron expression.**

I have used cron expression syntax to implement a new feature called *deploy freeze window* in Gitploy.

## Install

```shell
go get github.com/gitploy-io/cronexpr
```

## Usage

Import the package, first:

```go
import "time"
import "github.com/gitploy-io/cronexpr"
```

```go
prevTime := cronexpr.MustParse("0 0 29 * *").Prev(time.Now())
nextTime := cronexpr.MustParse("0 0 29 * *").Next(time.Now())
```

You can check the detail in the `_example` directory.


## Implementation


```
Field name     Mandatory?   Allowed values    Allowed special characters
----------     ----------   --------------    --------------------------
Minutes        Yes          0-59              * / , -
Hours          Yes          0-23              * / , -
Day of month   Yes          1-31              * / , - 
Month          Yes          1-12 or JAN-DEC   * / , -
Day of week    Yes          0-6 or SUN-SAT    * / , - 
```

### Asterisk ( * )
The asterisk indicates that the cron expression matches for all values of the field. E.g., using an asterisk in the 4th field (month) indicates every month. 

### Slash ( / )
Slashes describe increments of ranges. For example `3-59/15` in the minute field indicate the third minute of the hour and every 15 minutes thereafter. The form `*/...` is equivalent to the form "first-last/...", that is, an increment over the largest possible range of the field.

### Comma ( , )
Commas are used to separate items of a list. For example, using `MON,WED,FRI` in the 5th field (day of week) means Mondays, Wednesdays and Fridays.

### Hyphen ( - )
Hyphens define ranges. For example, 2000-2010 indicates every year between 2000 and 2010 AD, inclusive.

## Details

* The return value of `Next` and `Prev` is zero if the pattern doesn't match in five years.
