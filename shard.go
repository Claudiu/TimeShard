package timeshard

import (
	"container/heap"
	"time"
)

type Shard struct {
	data []byte
	meta []uint64
}

type MetaType uint64

const (
	// Where data is stored in our slice, index of []byte
	MetaDataIndex MetaType = iota

	// How many runes this action will affect (i.e. for DeleteOp, how many runes we will delete)
	// For InsertOp, the size of our data in bytes (used for keeping data in the slice)
	MetaDataByteSize

	// The type of our operation Delete or Insert
	MetaOperation

	// How many runes to skip
	MetaRetain

	// When was the action triggered, used for sorting and keeping data integrity
	MetaTimestamp
)

const (
	// Will Insert N runes at position
	OpInsert = 0

	// Will Delete N runes at position
	OpDelete = 1
)

// The size of our Meta
const MetaSize = 5

func NewShard() Shard {
	newShard := Shard{
		data: make([]byte, 0),
		meta: make([]uint64, 0),
	}

	newShard.Init()

	return newShard
}

func (s *Shard) Init() {
	heap.Init(s)
}

// Len returns the length of a shard (how many operations it produces)
// We get this by dividing the length of our meta slice, to the MetaSize count constant
// I.E. A batch with one insert will return 1 when we call Len()
func (s Shard) Len() int { return len(s.meta) / MetaSize }

func (s Shard) Less(i, j int) bool {
	exp1, _ := s.Get(i, MetaTimestamp)
	exp2, _ := s.Get(j, MetaTimestamp)
	return exp1 < exp2
}

func (s Shard) Swap(i, j int) {
	for x := 0; x < MetaSize; x++ {
		s.meta[i+x], s.meta[j+x] = s.meta[j+x], s.meta[i+x]
	}
}

func (s *Shard) Push(x interface{}) {
	(*s).meta = append((*s).meta, x.([]uint64)...)
}

func (s *Shard) Pop() interface{} {
	old := (*s).meta
	n := len(old)
	x := old[n-MetaSize]
	(*s).meta = old[0 : n-MetaSize]
	return x
}

// Get returns two things, the meta value and a boolean if the value was found.
// If the index overflows our slice, it will return 0 and false
func (s *Shard) Get(index int, k MetaType) (uint64, bool) {
	if s.Len() >= index && s.assertIntegrity() {
		return 0, false
	}

	return s.meta[uint64(index)*MetaSize+uint64(k)], true
}

// IsEmpty asserts whether we have 0 actions or not in our shard.
// I.E. An empty Batch will return 0
func (s *Shard) IsEmpty() bool {
	return s.Len() == 0
}

// This will assert if data passes simple checks (such as if it's divisible by MetaSize)
func (s *Shard) assertIntegrity() bool {
	return len(s.meta)%MetaSize == 0
}

func (s *Shard) pushData(b *[]byte) (dataStart, dataLength uint64) {
	dataStart = uint64(len(s.data))
	dataLength = uint64(len(*b))

	s.data = append(s.data, *b...)
	return
}

func (s *Shard) pushMeta(start, l, retain, action uint64) {
	now := time.Now().UnixNano()
	heap.Push(s, []uint64{start, l, action, retain, uint64(now)})
}
