package shift

import (
	"errors"
	"strconv"
	"strings"
)

type Time struct {
	Hour, Minute int
}

func NewTime(hour, minute int) Time {
	return Time{
		Hour:   hour,
		Minute: minute,
	}
}

func Parse(stringTime string) (Time, error) {
	parts := strings.Split(stringTime, ":")
	if len(parts) != 2 || len(parts[0]) != 2 || len(parts[1]) != 2 {
		return Time{}, errors.New("wrong time format")
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return Time{}, errors.New("can't parse hour " + parts[0])
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return Time{}, errors.New("can't parse minute " + parts[1])
	}
	return NewTime(hour, minute), nil
}

// Compare returns 1 if a > b, -1 if a < b, 0 if a == b
func Compare(a, b Time) int {
	hourRes := compareInts(a.Hour, b.Hour)
	if hourRes != 0 {
		return hourRes
	}
	return compareInts(a.Minute, b.Minute)
}

func compareInts(a, b int) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

func (t Time) Sub(start Time) Time {
	if t.Minute < start.Minute {
		return Time{
			Hour:   t.Hour - 1 - start.Hour,
			Minute: t.Minute + 60 - start.Minute,
		}
	}
	return Time{
		Hour:   t.Hour - start.Hour,
		Minute: t.Minute - start.Minute,
	}
}

func (t Time) Add(other Time) Time {
	minutes := t.Minute + other.Minute
	if minutes > 59 {
		t.Hour++
		minutes -= 60
	}
	return Time{
		Hour:   t.Hour + other.Hour,
		Minute: minutes,
	}
}

func (t Time) GetCeilHours() int {
	if t.Minute > 0 {
		return t.Hour + 1
	}
	return t.Hour
}

func (t Time) String() string {
	var strTime string
	if t.Hour < 10 {
		strTime += "0"
	}
	strTime += strconv.Itoa(t.Hour) + ":"
	if t.Minute < 10 {
		strTime += "0"
	}
	strTime += strconv.Itoa(t.Minute)
	return strTime
}
