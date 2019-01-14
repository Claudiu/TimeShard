package timeshard

import (
	"testing"
)

func TestForwardIterator_Retain(t *testing.T) {
	c := NewBatch()
	c.Insert(200, []byte("Sample text"))
	c.Insert(500, []byte("Sample text"))

	iter := c.Snapshot().Iterator(false)

	for iter.HasNext() {
		if iter.GetMeta(MetaRetain) != 200 {
			t.FailNow()
		}

		break
	}
}

func TestReverseIterator_Retain(t *testing.T) {
	c := NewBatch()
	c.Insert(200, []byte("Sample text"))
	c.Insert(500, []byte("Sample text"))

	iter := c.Snapshot().Iterator(true)

	for iter.HasNext() {
		if iter.GetMeta(MetaRetain) != 500 {
			t.FailNow()
		}

		break
	}
}

func TestSnapshot_Squash(t *testing.T) {
	c := NewBatch()

	c.Insert(0, []byte("a"))
	c.Insert(1, []byte("l"))
	c.Insert(2, []byte("b"))
	c.Insert(3, []byte("a"))
	c.Insert(4, []byte("t"))
	c.Insert(5, []byte("r"))
	c.Insert(6, []byte("o"))
	c.Insert(7, []byte("s"))

	iter := c.Squash(0).Snapshot().Iterator(true)
	for iter.HasNext() {
		if string(iter.Value()) != "albatros" {
			t.Fail()
		}
	}
}

func TestBatch_Delete(t *testing.T) {
	c := NewBatch()

	c.Insert(0, []byte("alb"))
	c.Insert(3, []byte("a"))
	c.Insert(4, []byte("t"))
	c.Insert(5, []byte("r"))
	c.Insert(6, []byte("o"))
	c.Insert(7, []byte("s"))

	c.Delete(0, 3)

	iter := c.Squash(0).Snapshot().Iterator(true)
	for iter.HasNext() {
		if string(iter.Value()) != "atros" {
			t.Fail()
		}
	}
}

func TestBatch_DeleteOutOfBounds(t *testing.T) {
	c := NewBatch()

	c.Insert(0, []byte("a"))
	c.Insert(1, []byte("l"))
	c.Insert(2, []byte("b"))
	c.Insert(3, []byte("a"))
	c.Insert(4, []byte("t"))
	c.Insert(5, []byte("r"))
	c.Insert(6, []byte("o"))
	c.Insert(7, []byte("s"))
	c.Delete(0, 3000)
	c.Insert(0, []byte("imi plac merele"))
	c.Delete(0, 9)

	iter := c.Squash(0).Snapshot().Iterator(true)
	for iter.HasNext() {
		if string(iter.Value()) != "merele" {
			t.Fail()
		}
	}
}

func BenchmarkBatch_Insert(b *testing.B) {
	b.StopTimer()
	c := NewBatch()

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		c.Insert(0, []byte("a"))
	}
}

func BenchmarkBatch_Squash(b *testing.B) {
	b.StopTimer()
	c := NewBatch()

	for n := uint64(0); n < 1000; n++ {
		c.Insert(n, []byte("lorem ipsum dolor"))
	}

	snap := c.Snapshot()

	b.StartTimer()

	newBatch := snap.Squash(0)

	for iter := newBatch.Snapshot().Iterator(false); iter.HasNext(); {
		iter.Value()
	}
}
