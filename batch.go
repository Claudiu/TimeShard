package timeshard

import (
	"time"
)

type Batch struct {
	Shard
}

func NewBatch() *Batch {
	return &Batch{Shard{
		make([]byte, 0),
		make([]uint64, 0),
		},
	}
}

func (c *Batch) Merge(other *Batch) *Batch {
	snap1 := c.Snapshot()
	snap2 := other.Snapshot()

	integrity := snap1.assertIntegrity() && snap2.assertIntegrity()
	if !integrity {
		return nil
	}

	first, last := &snap1, &snap2
	if snap1.LastActivity() > snap2.LastActivity() {
		first, last = &snap2, &snap1
	}

	// Data does not need to be moved...
	//for i := 0; i < len((*last).meta); i += MetaSize {
	//	(*last).meta[i] += uint64(len((*first).data))
	//}

	dataMerged := append((*first).data, (*last).data...)
	metaMerged := append((*first).meta, (*last).meta...)

	doc := NewBatch()
	doc.data = append(doc.data, dataMerged...)
	doc.meta = append(doc.meta, metaMerged...)

	return doc
}

func (c *Batch) MarshalJSON() ([]byte, error) {
	return c.Snapshot().MarshalJSON()
}

func (c *Batch) Snapshot() *Snapshot {
	snap := NewSnapshot()
	snap.Shard = Shard{
		make([]byte, len(c.data)),
		make([]uint64, len(c.meta)),
	}

	copy(snap.data, c.data)
	copy(snap.meta, c.meta)

	return snap
}

func (c *Batch) Squash(s uint64) *Batch {
	return c.Snapshot().Squash(s)
}

func (c *Batch) add(rawBytes []byte, action, retain uint64) {
	s, l := c.pushData(&rawBytes)
	c.pushMeta(s, l, retain, OpInsert)
}

func (c *Batch) Insert(at uint64, rawBytes []byte) {
	c.add(rawBytes, OpInsert, at)
}

func (c *Batch) Delete(at uint64, count uint64) {
	now := time.Now().UnixNano()

	s := uint64(len(c.data))
	c.meta = append(c.meta, []uint64{s, count, OpDelete, at, uint64(now)}...)
}
