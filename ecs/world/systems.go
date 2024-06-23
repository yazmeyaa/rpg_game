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
	dt    time.Duration
	items []System
}

func NewSystems(systemsCount int) *Systems {
	return &Systems{
		items: make([]System, 0, systemsCount),
	}
}

func (s *Systems) SetDT(dt time.Duration) {
	s.dt = dt
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

func (s *Systems) Update() {
	for _, system := range s.items {
		system.Compute(s.dt)
	}
}

func (s *Systems) StartUpdating(ctx context.Context, updateInterval time.Duration) {
	ticker := time.NewTicker(updateInterval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.Update()
				continue
			}
		}
	}()
}
