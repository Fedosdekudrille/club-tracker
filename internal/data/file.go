package data

import (
	"bufio"
	"github.com/Fedosdekudrille/club-tracker/internal/club"
	"github.com/Fedosdekudrille/club-tracker/internal/shift"
	"github.com/Fedosdekudrille/club-tracker/pkg/queue"
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

// MustParseClubManager parses all management data from scanner
func MustParseClubManager(scanner *bufio.Scanner) *club.Manager {
	scanner.Scan()
	tableNumTxt := scanner.Text()
	tableNum, err := strconv.Atoi(tableNumTxt)
	if err != nil || tableNum < 0 {
		panic(tableNumTxt)
	}

	zeroTime := shift.NewTime(0, 0)
	maxTime := shift.NewTime(23, 59)
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
	if shift.Compare(shiftTime[0], shiftTime[1]) > 0 || shift.Compare(shiftTime[0], zeroTime) == -1 || shift.Compare(shiftTime[1], maxTime) == 1 {
		panic(strTimes)
	}

	scanner.Scan()
	costPerHourTxt := scanner.Text()
	costPerHour, err := strconv.Atoi(costPerHourTxt)
	if err != nil || costPerHour < 0 {
		panic(costPerHourTxt)
	}

	return club.NewManager(shiftTime[0], shiftTime[1], costPerHour, tableNum)
}

// MustQueueEvents parses all events from scanner and queues them adding internal events
func MustQueueEvents(scanner *bufio.Scanner, manager *club.Manager) queue.Queue[string] {
	outputQueue := queue.NewQueue[string]()
	prevTime := shift.NewTime(0, 0)
	nameRegex := regexp.MustCompile(namePattern)
	shiftEnded := false
	for scanner.Scan() {
		event := scanner.Text()
		code, time, name, table := parseEvent(event, prevTime, nameRegex)
		prevTime = time
		if code == ClientSatAtTable && (table < 1 || table > manager.GetTableNum()) {
			panic(event)
		}
		if !shiftEnded && shift.Compare(time, manager.GetEndTime()) > 0 {
			outputQueue.Push(manager.GetEndTime().String())
			endShift(manager, &outputQueue)
			shiftEnded = true
		}
		outputQueue.Push(event)
		var answer club.Response
		switch code {
		case ClientCame:
			answer = manager.AddClient(name, time)
		case ClientSatAtTable:
			answer = manager.SetClientTable(name, table, time)
		case ClientWaits:
			answer = manager.WaitForTable(name, time)
		case ClientLeft:
			answer = manager.RemoveClient(name, time)
		}
		if answer.Code != 0 {
			outputQueue.Push(answer.String())
		}
	}
	if !shiftEnded {
		endShift(manager, &outputQueue)
	}
	return outputQueue
}
func parseEvent(event string, prevTime shift.Time, nameRegexp *regexp.Regexp) (code int, time shift.Time, name string, table int) {
	eventParts := strings.Split(event, " ")
	if len(eventParts) < 3 {
		panic(event)
	}

	time, err := shift.Parse(eventParts[0])
	if err != nil || shift.Compare(time, prevTime) < 0 {
		panic(event)
	}
	prevTime = time

	code, err = strconv.Atoi(eventParts[1])
	if err != nil {
		panic(event)
	}

	name = eventParts[2]
	if !nameRegexp.MatchString(name) {
		panic(event)
	}

	if code == ClientSatAtTable {
		if len(eventParts) < 4 {
			panic(event)
		}
		table, err = strconv.Atoi(eventParts[3])
		if err != nil {
			panic(event)
		}
	} else {
		table = -1
	}
	return code, time, name, table

}

func endShift(manager *club.Manager, outputQueue *queue.Queue[string]) {
	sortedLeaveEvents := manager.RemoveAllClientsSorted()
	for _, leaveEvent := range sortedLeaveEvents {
		outputQueue.Push(leaveEvent.String())
	}
}
