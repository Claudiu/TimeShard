package timeshard

import "testing"

func TestShard_Len(t *testing.T) {
	snap := NewBlock()
	snap.Insert(0, []byte{}, nil)
	snap.Delete(0, 1)

	if snap.Len() != 2 {
		t.Fail()
	}
}

func TestShard_Get(t *testing.T) {
	snap := NewBlock()
	snap.Insert(0, []byte{}, nil)
	snap.Delete(0, 1)

	if val, ok := snap.Get(0, MetaOperation); val != OpInsert && ok != true {
		t.Fail()
	}
}

func TestShard_IsEmpty(t *testing.T) {
	snap := NewBlock()

	if snap.IsEmpty() != true {
		t.Fail()
	}
}
