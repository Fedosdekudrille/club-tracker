package data

import (
	"bufio"
	"club-tracker/internal/club"
	"club-tracker/internal/shift"
	"club-tracker/pkg/queue"
	"regexp"
	"strconv"
	"strings"
)

const (
	ClientCame = iota + 1
	ClientSatAtTable
	ClientWaits
	ClientLeft
)

const namePattern = "^[A-Za-z0-9_]+$"

// MustParseClubManager parses all managment data from scanner
func MustParseClubManager(scanner *bufio.Scanner) *club.Manager {
	scanner.Scan()
	tableNumTxt := scanner.Text()
	tableNum, err := strconv.Atoi(tableNumTxt)
	if err != nil {
		panic(tableNumTxt)
	}

	scanner.Scan()
	strTimes := scanner.Text()
	openCloseTimes := strings.Split(strTimes, " ")
	if len(openCloseTimes) != 2 {
		panic(strTimes)
	}
	shiftTime := make([]shift.Time, 2)
	for id, time := range openCloseTimes {
		shiftTime[id], err = shift.Parse(time)
		if err != nil {
			panic(strTimes)
		}
	}
	if shiftTime[0].Compare(shiftTime[1]) > 0 {
		panic(strTimes)
	}

	scanner.Scan()
	costPerHourTxt := scanner.Text()
	costPerHour, err := strconv.Atoi(costPerHourTxt)
	if err != nil {
		panic(costPerHourTxt)
	}

	return club.NewManager(shiftTime[0], shiftTime[1], costPerHour, tableNum)
}

// MustQueueEvents parses all events from scanner and queues them adding internal events
func MustQueueEvents(scanner *bufio.Scanner, manager *club.Manager) queue.Queue[string] {
	outputQueue := queue.NewQueue[string]()
	prevTime := shift.NewTime(0, 0)
	for scanner.Scan() {
		event := scanner.Text()
		eventParts := strings.Split(event, " ")
		if len(eventParts) < 3 {
			panic(event)
		}

		time, err := shift.Parse(eventParts[0])
		if err != nil || time.Compare(prevTime) < 0 {
			panic(event)
		}
		prevTime = time
		if time.Compare(manager.GetEndTime()) > 0 {
			panic(event)
		}

		code, err := strconv.Atoi(eventParts[1])
		if err != nil {
			panic(event)
		}

		nameRegex := regexp.MustCompile(namePattern)
		if !nameRegex.MatchString(eventParts[2]) {
			panic(event)
		}

		switch code {
		case ClientCame:
			outputQueue.Push(event)
			answer := manager.AddClient(eventParts[2], time)
			if answer.Code != 0 {
				outputQueue.Push(answer.String())
			}
		case ClientSatAtTable:
			if len(eventParts) < 4 {
				panic(event)
			}
			table, err := strconv.Atoi(eventParts[3])
			if err != nil {
				panic(event)
			}
			outputQueue.Push(event)
			answer := manager.SetClientTable(eventParts[2], table, time)
			if answer.Code != 0 {
				outputQueue.Push(answer.String())
			}
		case ClientWaits:
			outputQueue.Push(event)
			answer := manager.WaitForTable(eventParts[2], time)
			if answer.Code != 0 {
				outputQueue.Push(answer.String())
			}
		case ClientLeft:
			outputQueue.Push(event)
			answer := manager.RemoveClient(eventParts[2], time)
			if answer.Code != 0 {
				outputQueue.Push(answer.String())
			}
		}
	}
	sortedLeaveEvents := manager.RemoveAllClientsSorted()
	for _, leaveEvent := range sortedLeaveEvents {
		outputQueue.Push(leaveEvent.String())
	}
	return outputQueue
}
