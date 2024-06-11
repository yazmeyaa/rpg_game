package utils_test

import (
	"testing"

	"github.com/rpg_game/utils"
)

func TestSparseSetAdd(t *testing.T) {
	ss := utils.NewSparseSet(32)
	ss.Add(12)
	if ss.Contains(12) == false {
		t.Error("Value not added")
	}
}

func TestSparseSetContains(t *testing.T) {
	ss := utils.NewSparseSet(32)
	ss.Add(16)

	notExisting := ss.Contains(10)
	existing := ss.Contains(16)
	defaultValue := ss.Contains(0)

	if notExisting == true {
		t.Error("Element must not be in sparse set")
		return
	}
	if existing == false {
		t.Error("Element must be in sparse set")
	}
	if defaultValue == true {
		t.Error("Element 0 must not be in sparse set")
	}
}

func TestSparseSetRange(t *testing.T) {
	ss := utils.NewSparseSet(32)
	for idx := range 10 {
		ss.Add(idx * 2)
	}

	var expectedValues []int = []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}

	i := 0
	ss.Range(func(id utils.Id) {
		if expectedValues[i] != id {
			t.Errorf("Expected value: <%d>, recieved: <%d>", expectedValues[i], id)
		}
		i++
	})
}

func TestRemove(t *testing.T) {
	ss := utils.NewSparseSet(32)

	ss.Add(12)
	exist := ss.Contains(12)
	if !exist {
		t.Error("Added value is not exist")
	}

	ss.Remove(12)
	exist = ss.Contains(12)
	if exist {
		t.Error("Value exist after element remove")
	}
}
