package club

import (
	"github.com/Fedosdekudrille/club-tracker/internal/shift"
	"github.com/Fedosdekudrille/club-tracker/pkg/queue"
	"slices"
)

const (
	Zero                = iota
	ClientForcedToLeave = iota + 10
	ClientGotToTable
	Error
)

type clientInfo struct {
	table     int
	startTime shift.Time
}

// Manager is responsible for managing clients, tables and events
type Manager struct {
	clients        map[string]*clientInfo
	waitingClients queue.Queue[string]

	tables     []table
	freeTables int

	startTime, endTime shift.Time
	costPerHour        int
}

func NewManager(startTime, endTime shift.Time, costPerHour int, tableNum int) *Manager {
	return &Manager{
		clients:        make(map[string]*clientInfo),
		waitingClients: queue.NewQueue[string](),
		tables:         make([]table, tableNum),
		startTime:      startTime,
		endTime:        endTime,
		costPerHour:    costPerHour,
		freeTables:     tableNum,
	}
}

func (m *Manager) AddClient(name string, currentTime shift.Time) Response {
	if shift.Compare(currentTime, m.startTime) < 0 || shift.Compare(currentTime, m.endTime) > 0 {
		return responseErr("NotOpenYet", currentTime)
	}
	if _, ok := m.clients[name]; ok {
		return responseErr("YouShallNotPass", currentTime)
	}
	m.clients[name] = &clientInfo{
		table: -1,
	}
	return responseZero()
}

func (m *Manager) SetClientTable(name string, table int, currentTime shift.Time) Response {
	if client, ok := m.clients[name]; ok {
		table--
		if m.tables[table].isTaken {
			return responseErr("PlaceIsBusy", currentTime)
		}
		m.leaveTable(client, currentTime)

		client.table = table
		m.tables[table].isTaken = true
		client.startTime = currentTime
		m.freeTables--
		return responseZero()
	}
	return responseErr("ClientUnknown", currentTime)
}

func (m *Manager) WaitForTable(name string, currentTime shift.Time) Response {
	if client, ok := m.clients[name]; ok {
		if m.freeTables > 0 {
			return responseErr("ICanWaitNoLonger!", currentTime)
		}
		if m.waitingClients.Len() > len(m.tables) {
			m.RemoveClient(name, currentTime)
			return newResponse(currentTime, ClientForcedToLeave, name, 0, "")
		}
		res := m.freeTable(client, currentTime)
		m.waitingClients.Push(name)
		return res
	}
	return responseErr("ClientUnknown", currentTime)
}

func (m *Manager) RemoveClient(name string, currentTime shift.Time) Response {
	if client, ok := m.clients[name]; ok {
		res := m.freeTable(client, currentTime)
		delete(m.clients, name)
		return res
	}
	return responseErr("ClientUnknown", currentTime)
}

func (m *Manager) RemoveAllClientsSorted() (sortedClients []Response) {
	sortedClients = make([]Response, 0, len(m.clients))
	for name, client := range m.clients {
		sortedClients = append(sortedClients, newResponse(m.endTime, ClientForcedToLeave, name, 0, ""))
		m.freeTable(client, m.endTime)
	}
	slices.SortFunc(sortedClients, func(a, b Response) int {
		if a.UserName < b.UserName {
			return -1
		} else if a.UserName > b.UserName {
			return 1
		}
		return 0
	})
	m.clients = make(map[string]*clientInfo)
	m.waitingClients = queue.NewQueue[string]()
	return sortedClients
}

func (m *Manager) GetStartTime() shift.Time {
	return m.startTime
}

func (m *Manager) GetEndTime() shift.Time {
	return m.endTime
}
