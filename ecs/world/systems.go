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

func NewSystems() *Systems {
	return &Systems{
		items: make([]System, 0),
	}
}

func (s *Systems) sortSystems() {
	sort.Slice(s.items, func(i, j int) bool {
		return s.items[i].Priority() < s.items[j].Priority()
	})
}

func (s *Systems) AddSystem(system System) {
	// Увеличиваем capacity на 1 вручную, если слайс переполнен
	if len(s.items) == cap(s.items) {
		newItems := make([]System, len(s.items), len(s.items)+1)
		copy(newItems, s.items)
		s.items = newItems
	}
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
			}
		}
	}()
}
