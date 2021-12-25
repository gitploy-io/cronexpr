package cronexpr

import (
	"reflect"
	"strconv"
	"testing"
)

func Test_Parse(t *testing.T) {
	t.Run("Parse the cron expr", func(t *testing.T) {
		// TODO:
	})
}

func Test_parseField(t *testing.T) {
	t.Run("Parse the field.", func(t *testing.T) {
		cases := []struct {
			value string
			b     bound
			t     translater
			want  bitset
		}{
			{
				value: "*",
				b:     boundMinute,
				want:  buildBitset(0, 59, 1),
			},
			{
				value: "5",
				b:     boundMinute,
				want:  1 << 5,
			},
			{
				value: "5,10",
				b:     boundMinute,
				want:  1<<5 | 1<<10,
			},
			{
				value: "5-20",
				b:     boundMinute,
				want:  buildBitset(5, 20, 1),
			},
			{
				value: "JAN-MAR",
				b:     boundMonth,
				t:     translaterMonth,
				want:  buildBitset(1, 3, 1),
			},
			{
				value: "5-20/5",
				b:     boundMinute,
				want:  (1 << 5) | (1 << 10) | (1 << 15) | (1 << 20),
			},
			{
				value: "*/20",
				b:     boundMinute,
				want:  (1 << 0) | (1 << 20) | (1 << 40),
			},
		}

		for _, c := range cases {
			b, err := parseField(c.value, c.b, c.t)
			if err != nil {
				t.Fatalf("parseField returns an error: %s", err)
			}

			if !reflect.DeepEqual(b, c.want) {
				t.Fatalf("parseField(%s) = %v, wanted %v", c.value, strconv.FormatUint(uint64(b), 2), strconv.FormatUint(uint64(c.want), 2))
			}
		}
	})

	t.Run("Return an error of parsing the field.", func(t *testing.T) {
		cases := []struct {
			value string
			b     bound
			t     translater
		}{
			{
				value: "1-3-5/3",
				b:     boundMinute,
			},
		}

		for _, c := range cases {
			_, err := parseField(c.value, c.b, c.t)
			if err == nil {
				t.Fatal("parseField must returns an error")
			}
		}
	})
}
