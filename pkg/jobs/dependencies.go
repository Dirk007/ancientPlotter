package jobs

import (
	"github.com/Dirk007/ancientPlotter/pkg/broadcast"
	"github.com/Dirk007/ancientPlotter/pkg/feeder"
)

type ContextDependencies struct {
	Alive *broadcast.BroadcastChannel[string]
	Logs  *broadcast.BroadcastChannel[string]
	Stats *broadcast.BroadcastChannel[feeder.Stats]
	Jobs  *JobCollection
}

func NewContextDependencies() *ContextDependencies {
	alive := broadcast.NewBroadcastChannel[string]()
	logs := broadcast.NewBroadcastChannel[string]()
	stats := broadcast.NewBroadcastChannel[feeder.Stats]()
	jobs := NewJobCollection()
	return &ContextDependencies{
		Alive: alive,
		Logs:  logs,
		Stats: stats,
		Jobs:  jobs,
	}
}
