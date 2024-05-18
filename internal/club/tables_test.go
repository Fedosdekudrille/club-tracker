package club

import (
	"github.com/Fedosdekudrille/club-tracker/internal/shift"
	"testing"
)

func TestManager_freeTable(t *testing.T) {
	manager := NewManager(shift.NewTime(9, 0), shift.NewTime(17, 0), 5, 10)
	manager.clients["test"] = &clientInfo{table: 0, startTime: shift.NewTime(9, 0)}
	manager.clients["test2"] = &clientInfo{table: -1, startTime: shift.Time{}}
	manager.waitingClients.Push("test2")
	t.Run("Positive", func(t *testing.T) {
		manager.freeTable(manager.clients["test"], shift.NewTime(10, 0))
		if manager.clients["test"].table != -1 {
			t.Error("Table should be -1")
		}
		if manager.clients["test2"].table != 0 {
			t.Error("Should sit at table 0")
		}
	})
	t.Run("Negative", func(t *testing.T) {
		res := manager.freeTable(manager.clients["test"], shift.NewTime(10, 0))
		if res.Code != Zero {
			t.Error("Code should be Zero")
		}
	})
}
