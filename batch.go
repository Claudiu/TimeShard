package timeshard

import (
	"fmt"
	"time"
)

// Batch is a series of alterable shards, if you need to read it
// You need to create a Snapshot by calling the method.
type Batch struct {
	Shard
}

// NewBatch will initialize an empty batch (with empty slices).
func NewBatch() *Batch {
	return &Batch{Shard: NewShard()}
}

// Merge will concatenate two batches into a single one chronologically
func (c *Batch) Merge(other *Batch) (*Batch, error) {
	snap1 := c.Snapshot()
	snap2 := other.Snapshot()

	if integrity := snap1.assertIntegrity(); integrity {
		return nil, fmt.Errorf("first Batch failed integrity check")
	}

	if integrity := snap2.assertIntegrity(); integrity {
		return nil, fmt.Errorf("second Batch failed integrity check")
	}

	// Verify whether snap1 is older then snap2
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

	return doc, nil
}

// MarshalJSON is a shorthand operation for creating a Snapshot
// and then calling Marshal on it.
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

// Squash is a shorthand operation for creating a Snapshot
// and then calling Squash on it (reducing it to a single operation)
func (c *Batch) Squash(s uint64) *Batch {
	return c.Snapshot().Squash(s)
}

// add will add an operation into our Shard
func (c *Batch) add(rawBytes []byte, action, retain uint64) {
	s, l := c.pushData(&rawBytes)
	c.pushMeta(s, l, retain, OpInsert)
}

// Insert will add an OpInsert into our Shard
func (c *Batch) Insert(at uint64, rawBytes []byte) {
	c.add(rawBytes, OpInsert, at)
}

// Delete will add an OpDelete into our Shard
// TODO: Use add
func (c *Batch) Delete(at uint64, count uint64) {
	now := time.Now().UnixNano()

	s := uint64(len(c.data))
	c.meta = append(c.meta, []uint64{s, count, OpDelete, at, uint64(now)}...)
}
