package timeshard

import (
	"bytes"
	"errors"
	"sync"
	"time"
	"unicode/utf8"
)

type Snapshot struct {
	Shard
	sync.Mutex
}

func NewSnapshot() *Snapshot {
	return &Snapshot{
		Shard: Shard{
			make([]byte, 0),
			make([]uint64, 0),
		},
		Mutex: sync.Mutex{},
	}
}

func (snapshot *Snapshot) Iterator(reverse bool) Iterator {
	var iterator Iterator

	if reverse {
		iterator = &ReverseIterator{
			over: snapshot.Clone(),
		}

		iterator.Init()
		return iterator
	}

	iterator = &ForwardIterator{
		over: snapshot.Clone(),
	}

	iterator.Init()
	return iterator
}

func (snapshot *Snapshot) LastActivity() uint64 {
	data, _ := snapshot.Get(0, MetaTimestamp)
	return data
}

func (snapshot *Snapshot) Squash(count uint64) *Snapshot {
	//newSnap := NewEmpty()

	mirror := ""

	iter := snapshot.Iterator(false)

	current := uint64(0)
	for iter.HasNext() && (current < count || count == 0) {
		maxLen := uint64(utf8.RuneCountInString(mirror))

		switch iter.GetMeta(MetaOperation) {
		case OpInsert:
			pos := iter.GetMeta(MetaRetain)
			if pos > maxLen {
				pos = maxLen
			}

			n := string(bytes.Runes(iter.Value()))
			mirror = mirror[:pos] + n + mirror[pos:]
		case OpDelete:
			pos := iter.GetMeta(MetaRetain)
			if pos > maxLen {
				continue
			}

			affected := iter.GetMeta(MetaDataByteSize)
			if pos+affected > maxLen {
				affected = maxLen - pos
			}

			mirror = mirror[:pos] + mirror[pos+affected:]
		}

		current++
	}

	b := NewSnapshot()
	b.Insert(0, []byte(mirror))

	return b
}

// add will add an operation into our Shard
func (snapshot *Snapshot) add(rawBytes []byte, action, retain uint64) {
	snapshot.Lock()
	defer snapshot.Unlock()

	var s, l uint64
	if foundIndex := bytes.Index(snapshot.data, rawBytes); foundIndex != -1 {
		s = uint64(foundIndex)
		l = uint64(len(rawBytes))
	} else {
		s, l = snapshot.pushData(&rawBytes)
	}

	snapshot.pushMeta(s, l, retain, OpInsert)
}

// Insert will add an OpInsert into our Shard
func (snapshot *Snapshot) Insert(at uint64, rawBytes []byte) {
	snapshot.add(rawBytes, OpInsert, at)
}

// Delete will add an OpDelete into our Shard
// TODO: Use add
func (snapshot *Snapshot) Delete(at uint64, count uint64) {
	snapshot.Lock()
	defer snapshot.Unlock()

	now := time.Now().UnixNano()

	currentLength := uint64(len(snapshot.data))
	data := []uint64{currentLength, count, OpDelete, at, uint64(now)}
	snapshot.meta = append(snapshot.meta, data...)
}

// Clone will clone the current snapshot, useful for iterating
func (snapshot *Snapshot) Clone() *Snapshot {
	snapshot.Lock()
	defer snapshot.Unlock()

	snap := NewSnapshot()

	snap.Shard = Shard{
		make([]byte, len(snapshot.data)),
		make([]uint64, len(snapshot.meta)),
	}

	copy(snap.data, snapshot.data)
	copy(snap.meta, snapshot.meta)

	return snap
}

func (snapshot *Snapshot) applyCommits(otherSnapshot *Snapshot, target *Snapshot) (*Snapshot, error) {
	if integrity := snapshot.assertIntegrity(); integrity {
		return nil, errors.New("first Batch failed integrity check")
	}

	if integrity := otherSnapshot.assertIntegrity(); integrity {
		return nil, errors.New("second Batch failed integrity check")
	}

	// Verify whether snapshot is older then otherSnapshot
	first, last := &snapshot, &otherSnapshot
	if snapshot.LastActivity() > otherSnapshot.LastActivity() {
		first, last = &otherSnapshot, &snapshot
	}

	// Data does not need to be moved...
	//for i := 0; i < len((*last).meta); i += MetaSize {
	//	(*last).meta[i] += uint64(len((*first).data))
	//}

	dataMerged := append((*first).data, (*last).data...)
	metaMerged := append((*first).meta, (*last).meta...)

	var targetSnapshot *Snapshot

	if target == nil {
		targetSnapshot = NewSnapshot()
	}

	targetSnapshot.data = append(targetSnapshot.data, dataMerged...)
	targetSnapshot.meta = append(targetSnapshot.meta, metaMerged...)

	return targetSnapshot, nil
}
