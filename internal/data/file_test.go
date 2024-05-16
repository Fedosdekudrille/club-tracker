package data

import (
	"bufio"
	"strings"
	"testing"
)

func TestMustParseClubManager(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		if MustParseClubManager(bufio.NewScanner(strings.NewReader("3\n09:00 19:00\n10\n"))).GetStartTime().Hour != 9 {
			t.Error("Wrong start time")
		}
	})
	t.Run("Negative", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("Must panic")
			}
		}()
		MustParseClubManager(bufio.NewScanner(strings.NewReader("3\n09:00 19:00 asdfasdf\n11\n")))
	})
	t.Run("WrongTimeOrder", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("Must panic")
			}
		}()
		MustParseClubManager(bufio.NewScanner(strings.NewReader("3\n19:00 09:00\n10\n")))
	})
}

func TestMustQueueEvents(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		manager := MustParseClubManager(bufio.NewScanner(strings.NewReader("3\n09:00 19:00\n10\n")))
		queue := MustQueueEvents(bufio.NewScanner(strings.NewReader("08:00 1 Dan")), manager)
		if queue.Pop() != "08:00 1 Dan" {
			t.Error("Doesn't return event")
		}
		if queue.Pop() != "08:00 13 NotOpenYet" {
			t.Error("Doesn't return mistake")
		}
	})
	t.Run("Negative", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("Must panic")
			}
		}()
		manager := MustParseClubManager(bufio.NewScanner(strings.NewReader("3\n09:00 19:00\n10\n")))
		MustQueueEvents(bufio.NewScanner(strings.NewReader("09:00 19:00\n10\n")), manager)
	})
	t.Run("WrongTimeOrder", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("Must panic")
			}
		}()
		manager := MustParseClubManager(bufio.NewScanner(strings.NewReader("3\n09:00 19:00\n10\n")))
		MustQueueEvents(bufio.NewScanner(strings.NewReader("08:00 1 Dan\n07:00 1 Man")), manager)
	})
	t.Run("ExitTime", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("Must panic")
			}
		}()
		manager := MustParseClubManager(bufio.NewScanner(strings.NewReader("3\n09:00 19:00\n10\n")))
		MustQueueEvents(bufio.NewScanner(strings.NewReader("20:00 1 Dan")), manager)
	})
}
