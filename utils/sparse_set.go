package utils

type SparseSet struct {
	size   int
	dense  []Id
	sparse []Id
}

type Id = int

func NewSparseSet(maxSize uint32) *SparseSet {
	return &SparseSet{
		size:   0,
		sparse: make([]Id, maxSize),
		dense:  make([]Id, 0, maxSize),
	}
}

func (ss *SparseSet) Contains(x Id) bool {
	return x < len(ss.sparse) && ss.sparse[x] < ss.size && ss.dense[ss.sparse[x]] == x
}

func (ss *SparseSet) Add(val Id) {
	if ss.Contains(val) {
		return
	}
	ss.sparse[val] = ss.size
	ss.dense = append(ss.dense, val)
	ss.size++
}

func (ss *SparseSet) Remove(x Id) {
	if !ss.Contains(x) {
		return
	}
	idx := ss.sparse[x]
	last := ss.dense[ss.size-1]
	ss.dense[idx] = last
	ss.sparse[last] = idx
	ss.size--
	ss.dense = ss.dense[:ss.size]
}

func (ss *SparseSet) Clear() {
	ss.size = 0
}

func (ss *SparseSet) GetAll() []Id {
	result := make([]Id, ss.size)
	copy(result, ss.dense)
	return result
}

func (ss *SparseSet) Range(cb func(id Id)) {
	for _, id := range ss.dense[:ss.size] {
		cb(id)
	}
}
