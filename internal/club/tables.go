package club

import (
	"club-tracker/internal/shift"
	"fmt"
)

type table struct {
	isTaken  bool
	BusyTime shift.Time
	Income   int
}

func countIncome(startTime, endTime shift.Time, costPerHour int) int {
	return endTime.Sub(startTime).GetCeilHours() * costPerHour
}

// freeTable frees table for waiting clients
func (m *Manager) freeTable(client *clientInfo, currentTime shift.Time) Response {
	if client.table != -1 {
		curTable := client.table
		m.leaveTable(client, currentTime)

		for m.waitingClients.Len() > 0 {
			newClientName := m.waitingClients.Pop()
			if newClient, ok := m.clients[newClientName]; ok {
				newClient.table = curTable
				newClient.startTime = currentTime
				m.tables[curTable].isTaken = true
				m.freeTables--
				return Response{Code: ClientGotToTable, Table: newClient.table + 1, Time: currentTime, UserName: newClientName}
			}
		}
	}
	return responseZero()
}

// leaveTable frees table from client, and counts income
func (m *Manager) leaveTable(client *clientInfo, currentTime shift.Time) {
	if client.table == -1 {
		return
	}
	m.tables[client.table].Income += countIncome(client.startTime, currentTime, m.costPerHour)
	m.tables[client.table].BusyTime = m.tables[client.table].BusyTime.Add(currentTime.Sub(client.startTime))
	m.tables[client.table].isTaken = false
	m.freeTables++
	client.table = -1
	client.startTime = shift.Time{}
}

func (m *Manager) GetAllTables() []table {
	return m.tables
}

func (t table) String() string {
	return fmt.Sprintf("%d %s", t.Income, t.BusyTime.String())
}
