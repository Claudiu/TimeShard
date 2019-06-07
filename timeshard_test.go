package timeshard

import (
	"testing"
)

func TestForwardIterator_Retain(t *testing.T) {
	c := NewBlock()
	c.Insert(200, []byte("Sample text"), nil)
	c.Insert(500, []byte("Sample text"), nil)

	iter := c.Iterator(false)

	for iter.HasNext() {
		if iter.GetMeta(MetaRetain) != 200 {
			t.FailNow()
		}

		break
	}
}

func TestReverseIterator_Retain(t *testing.T) {
	c := NewBlock()
	c.Insert(200, []byte("Sample text"), nil)
	c.Insert(500, []byte("Sample text"), nil)

	iter := c.Iterator(true)

	for iter.HasNext() {
		if iter.GetMeta(MetaRetain) != 500 {
			t.FailNow()
		}

		break
	}
}

func TestSnapshot_Squash(t *testing.T) {
	c := NewBlock()

	c.Insert(0, []byte("a"), nil)
	c.Insert(1, []byte("l"), nil)
	c.Insert(2, []byte("b"), nil)
	c.Insert(3, []byte("a"), nil)
	c.Insert(4, []byte("t"), nil)
	c.Insert(5, []byte("r"), nil)
	c.Insert(6, []byte("o"), nil)
	c.Insert(7, []byte("s"), nil)

	iter := c.Squash(0).Iterator(true)
	for iter.HasNext() {
		if string(iter.Value()) != "albatros" {
			t.Fail()
		}
	}
}

func TestSnapshot_SquashIssue1(t *testing.T) {
	c := NewBlock()

	c.Insert(0, []byte("alb"), nil)
	c.Insert(1, []byte("a"), nil)
	c.Insert(2, []byte("t"), nil)
	c.Insert(3, []byte("r"), nil)
	c.Insert(4, []byte("o"), nil)
	c.Insert(5, []byte("s"), nil)

	iter := c.Squash(1).Iterator(true)
	for iter.HasNext() {
		if string(iter.Value()) != "alb" {
			t.Fail()
		}
	}
}

func TestSnapshot_SquashEmoji(t *testing.T) {
	c := NewBlock()

	c.Insert(0, []byte("😀"), nil)
	c.Insert(1, []byte("a"), nil)
	c.Insert(2, []byte("t"), nil)
	c.Insert(3, []byte("r"), nil)
	c.Insert(4, []byte("o"), nil)
	c.Insert(5, []byte("s"), nil)

	iter := c.Squash(1).Iterator(true)
	for iter.HasNext() {
		if string(iter.Value()) != "😀" {
			t.Fail()
		}
	}
}

func TestBatch_Delete(t *testing.T) {
	c := NewBlock()

	c.Insert(0, []byte("alb"), nil)
	c.Insert(3, []byte("a"), nil)
	c.Insert(4, []byte("t"), nil)
	c.Insert(5, []byte("r"), nil)
	c.Insert(6, []byte("o"), nil)
	c.Insert(7, []byte("s"), nil)

	c.Delete(0, 3)

	iter := c.Squash(0).Iterator(true)
	for iter.HasNext() {
		if string(iter.Value()) != "atros" {
			t.Fail()
		}
	}
}

func TestBatch_DeleteOutOfBounds(t *testing.T) {
	c := NewBlock()

	c.Insert(0, []byte("a"), nil)
	c.Insert(1, []byte("l"), nil)
	c.Insert(2, []byte("b"), nil)
	c.Insert(3, []byte("a"), nil)
	c.Insert(4, []byte("t"), nil)
	c.Insert(5, []byte("r"), nil)
	c.Insert(6, []byte("o"), nil)
	c.Insert(7, []byte("s"), nil)
	c.Delete(0, 3000)
	c.Insert(0, []byte("imi plac merele"), nil)
	c.Delete(0, 9)

	iter := c.Squash(0).Iterator(true)
	for iter.HasNext() {
		if string(iter.Value()) != "merele" {
			t.Fail()
		}
	}
}

func BenchmarkBatch_Insert(b *testing.B) {
	b.StopTimer()
	c := NewBlock()

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		c.Insert(0, []byte("a"), nil)
	}
}

func BenchmarkBatch_Squash(b *testing.B) {
	b.StopTimer()
	c := NewBlock()

	for n := uint64(0); n < 1000; n++ {
		c.Insert(n, []byte("lorem ipsum dolor"), nil)
	}

	b.StartTimer()

	NewSnapshot := c.Squash(0)

	for iter := NewSnapshot.Iterator(false); iter.HasNext(); {
		iter.Value()
	}
}

func BenchmarkBatch_SquashReverse(b *testing.B) {
	b.StopTimer()
	c := NewBlock()

	for n := uint64(0); n < 1000; n++ {
		c.Insert(n, []byte("lorem ipsum dolor"), nil)
	}

	snap := c

	b.StartTimer()

	NewSnapshot := snap.Squash(0)

	for iter := NewSnapshot.Iterator(true); iter.HasNext(); {
		iter.Value()
	}
}
