package cronexpr

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	bitset uint64

	bound struct {
		min, max int
	}

	translater map[string]int
)

const (
	bitsetStar = 1<<64 - 1
)

var (
	boundMinute = bound{0, 59}
	boundHour   = bound{0, 24}
	boundDOM    = bound{1, 31}
	boundMonth  = bound{1, 12}
	boundDOW    = bound{0, 6}
)

var (
	translaterMonth = translater{
		"JAN": 1,
		"FEB": 2,
		"MAR": 3,
		"APR": 4,
		"MAY": 5,
		"JUN": 6,
		"JUL": 7,
		"AUG": 8,
		"SEP": 9,
		"OCT": 10,
		"NOV": 11,
		"DEC": 12,
	}

	translaterDay = translater{
		"SUN": 0,
		"MON": 1,
		"TUE": 2,
		"WED": 3,
		"THU": 4,
		"FRI": 5,
		"SAT": 6,
	}
)

func MustParse(expr string) *Schedule {
	s, err := Parse(expr)
	if err != nil {
		panic(err)
	}

	return s
}

func Parse(expr string) (*Schedule, error) {
	err := verifyExpr(expr)
	if err != nil {
		return nil, err
	}

	var (
		minute, hour, dom, month, dow bitset
	)

	fields := strings.Fields(strings.TrimSpace(expr))

	if minute, err = parseField(fields[0], boundMinute, translater{}); err != nil {
		return nil, err
	}

	if hour, err = parseField(fields[1], boundHour, translater{}); err != nil {
		return nil, err
	}

	if dom, err = parseField(fields[2], boundDOM, translater{}); err != nil {
		return nil, err
	}

	if month, err = parseField(fields[3], boundMonth, translaterMonth); err != nil {
		return nil, err
	}

	if dow, err = parseField(fields[4], boundDOW, translaterDay); err != nil {
		return nil, err
	}

	return &Schedule{
		Minute: minute,
		Hour:   hour,
		Dom:    dom,
		Month:  month,
		Dow:    dow,
	}, nil
}

// parseField returns an int with the bits set representing all of the times that
// the field represents or error parsing field value.
func parseField(field string, b bound, t translater) (bitset, error) {
	var bitsets bitset = 0

	// Split with "," (OR).
	fieldexprs := strings.Split(field, ",")
	for _, fieldexpr := range fieldexprs {
		b, err := parseFieldExpr(fieldexpr, b, t)
		if err != nil {
			return 0, err
		}

		bitsets = bitsets | b
	}

	return bitsets, nil
}

// parseFieldExpr returns the bits indicated by the given expression:
//   number | number "-" number [ "/" number ]
func parseFieldExpr(fieldexpr string, b bound, t translater) (bitset, error) {
	// Replace "*" into "min-max".
	newexpr := strings.Replace(fieldexpr, "*", fmt.Sprintf("%d-%d", b.min, b.max), 1)

	rangeAndStep := strings.Split(newexpr, "/")
	if !(len(rangeAndStep) == 1 || len(rangeAndStep) == 2) {
		return 0, fmt.Errorf("Failed to parse the expr '%s', too many '/'", fieldexpr)
	}

	hasStep := len(rangeAndStep) == 2

	// Parse the range, first.
	var (
		begin, end int
	)
	{
		lowAndHigh := strings.Split(rangeAndStep[0], "-")
		if !(len(lowAndHigh) == 1 || len(lowAndHigh) == 2) {
			return 0, fmt.Errorf("Failed to parse the expr '%s', too many '-'", fieldexpr)
		}

		low, err := parseInt(lowAndHigh[0], t)
		if err != nil {
			return 0, fmt.Errorf("Failed to parse the expr '%s': %w", fieldexpr, err)
		}

		begin = low

		// Special handling: "N/step" means "N-max/step".
		if len(lowAndHigh) == 1 && hasStep {
			end = b.max
		} else if len(lowAndHigh) == 1 && !hasStep {
			end = low
		} else if len(lowAndHigh) == 2 {
			high, err := parseInt(lowAndHigh[1], t)
			if err != nil {
				return 0, fmt.Errorf("Failed to parse the expr '%s': %w", fieldexpr, err)
			}

			end = high
		}
	}

	// Parse the step, second.
	step := 1
	if hasStep {
		var err error
		if step, err = strconv.Atoi(rangeAndStep[1]); err != nil {
			return 0, fmt.Errorf("Failed to parse the expr '%s': %w", fieldexpr, err)
		}
	}

	return buildBitset(begin, end, step), nil
}

func parseInt(s string, t translater) (int, error) {
	if i, err := strconv.Atoi(s); err == nil {
		return i, nil
	}

	i, ok := t[strings.ToUpper(s)]
	if !ok {
		return 0, fmt.Errorf("'%s' is out of reserved words", s)
	}

	return i, nil
}

func buildBitset(min, max, step int) bitset {
	var b bitset

	for i := min; i <= max; i += step {
		b = b | (1 << i)
	}

	return b
}

func verifyExpr(expr string) error {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return fmt.Errorf("The length of fields must be five.")
	}

	return nil
}
