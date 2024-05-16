package main

import (
	"bufio"
	"club-tracker/internal/data"
	"fmt"
	"os"
)

// standard program version, takes data from provided file, then processes it
func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if len(os.Args) < 2 {
		fmt.Println("Please provide path to the file")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

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
}
