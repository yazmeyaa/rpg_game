package world

import (
	"context"
	"sort"
	"time"
)

type System interface {
	Compute(dt time.Duration)
	Priority() uint8
}

type Systems struct {
	items []System
}

func NewSystems(systemsCount int) *Systems {
	return &Systems{
		items: make([]System, 0, systemsCount),
	}
}

func (s *Systems) sortSystems() {
	sort.Slice(s.items, func(i, j int) bool {
		return s.items[i].Priority() < s.items[j].Priority()
	})
}

func (s *Systems) AddSystem(system System) {
	s.items = append(s.items, system)

	s.sortSystems()
}

func (s *Systems) Update(dt time.Duration) {
	for _, system := range s.items {
		system.Compute(dt)
	}
}

func (s *Systems) StartUpdating(ctx context.Context, updateInterval time.Duration) {
	ticker := time.NewTicker(updateInterval)
	last := time.Now()

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				dt := time.Since(last)
				s.Update(dt)
				last = time.Now()
				continue
			}
		}
	}()
}
