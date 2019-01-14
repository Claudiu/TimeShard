package timeshard

import (
	"bytes"
	"unicode/utf8"
)

type Snapshot struct {
	Shard
}

func NewSnapshot() *Snapshot {
	return &Snapshot{Shard{
		make([]byte, 0),
		make([]uint64, 0),
		},
	}
}

func (snap *Snapshot) Iterator(reverse bool) Iterator {
	var iter Iterator

	if reverse {
		iter = &ReverseIterator{
			over: snap,
		}

		iter.Init()
		return iter
	}

	iter = &ForwardIterator{
		over: snap,
	}

	iter.Init()
	return iter
}

func (snap *Snapshot) LastActivity() uint64 {
	data, _ := snap.Get(0, MetaTimestamp)
	return data
}

func (snap *Snapshot) Squash(count uint64) *Batch {
	//newSnap := NewEmpty()

	mirror := ""

	iter := snap.Iterator(false)

	current := uint64(0)
	for iter.HasNext() && (current <= count || count == 0) {
		maxLen := uint64(utf8.RuneCountInString(mirror))

		switch iter.GetMeta(MetaOperation) {
		case OpInsert:
			pos := iter.GetMeta(MetaRetain)
			if pos > maxLen {
				pos = maxLen
			}

			n := string(bytes.Runes(iter.Value()))
			mirror = mirror[:pos] + n + mirror[pos:]

			current++
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
	}

	b := NewBatch()
	b.Insert(0, []byte(mirror))

	return b
}
