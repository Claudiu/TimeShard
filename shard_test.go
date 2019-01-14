package timeshard

import "testing"



func TestShard_Len(t *testing.T) {
	batch := NewBatch()
	batch.Insert(0, []byte{})
	batch.Delete(0, 1)

	if batch.Len() != 2 {
		t.Fail()
	}
}

func TestShard_Get(t *testing.T) {
	batch := NewBatch()
	batch.Insert(0, []byte{})
	batch.Delete(0, 1)

	if val, ok := batch.Get(0, MetaOperation);
		val != OpInsert && ok != true {
		t.Fail()
	}
}

func TestShard_IsEmpty(t *testing.T) {
	batch := NewBatch()

	if batch.IsEmpty() != true {
		t.Fail()
	}
}