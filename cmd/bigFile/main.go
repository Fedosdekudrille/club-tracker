package main

import (
	"bufio"
	"club-tracker/internal/data"
	"club-tracker/internal/shift"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// tests program with lots of data, does not take from file, generates by itself
func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	bigData := generateBigData()
	scanner := bufio.NewScanner(strings.NewReader(bigData))

	start := time.Now()

	clubManager := data.MustParseClubManager(scanner)
	outputEvents := data.MustQueueEvents(scanner, clubManager)

	fmt.Println(clubManager.GetStartTime())
	for !outputEvents.IsEmpty() {
		fmt.Println(outputEvents.Pop())
	}
	fmt.Println(clubManager.GetEndTime())

	tables := clubManager.GetAllTables()
	for id, table := range tables {
		fmt.Println(id+1, table)
	}

	fmt.Println(time.Since(start))
}

// generates 24 * 60 = 1440 lines of various data
func generateBigData() string {
	time := shift.NewTime(0, 0)
	builder := strings.Builder{}
	const tableNum = 5
	builder.WriteString(strconv.Itoa(tableNum) + "\n")
	builder.WriteString("00:00 24:00\n")
	builder.WriteString("15\n")

	for time.Hour < 24 {
		builder.WriteString(time.String() + " ")
		code := rand.Intn(4) + 1
		builder.WriteString(strconv.Itoa(code) + " ")
		for range 1 {
			builder.WriteByte('a' + byte(rand.Intn(26)))
		}
		if data.ClientSatAtTable == code {
			builder.WriteString(" " + strconv.Itoa(rand.Intn(tableNum)+1))
		}
		builder.WriteString("\n")
		time = time.Add(shift.Time{Minute: 1})
	}
	return builder.String()
}
