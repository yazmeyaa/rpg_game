package world

import (
	"fmt"
	"sort"
)

type System interface {
	Compute()
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

func (s *Systems) Update() {
	for _, system := range s.items {
		fmt.Printf("SYSTEM: >>>>>>>>>>>>>>>>>>>>>>>> \n\n%+v\n>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n\n", system)
		system.Compute()
	}
}
